<script setup lang="ts">
import type { CargoConfig } from '@/types/config'
import { Button } from '@/components/ui/button'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { CARGO_CONFIG_QUERY_KEY, useCargoConfigMutation } from '@/composables/use-config'
import { useQueryClient } from '@tanstack/vue-query'
import { toTypedSchema } from '@vee-validate/zod'
import { Loader } from 'lucide-vue-next'
import { useForm } from 'vee-validate'
import { z } from 'zod'

interface Props {
  initialValues: CargoConfig
}
const props = defineProps<Props>()

const cargoConfigSchema = z.object({
  liftPosition: z.number().min(0, 'Lift position must be at least 0'),
  lowerPosition: z.number().min(0, 'Lower position must be at least 0'),
  bottomDistanceHysteresis: z.object({
    lowerThreshold: z.number().int('Lower threshold must be an integer').positive('Lower threshold must be positive').default(15),
    upperThreshold: z.number().int('Upper threshold must be an integer').positive('Upper threshold must be positive').default(20),
  }),
}).refine(data => data.liftPosition < data.lowerPosition, {
  message: 'Lift position must be less than lower position',
  path: ['liftPosition'],
}).superRefine((data, ctx) => {
  if (data.bottomDistanceHysteresis.lowerThreshold >= data.bottomDistanceHysteresis.upperThreshold) {
    ctx.addIssue({
      code: z.ZodIssueCode.custom,
      message: 'Lower threshold must be less than upper threshold',
      path: ['bottomDistanceHysteresis', 'lowerThreshold'],
    })

    ctx.addIssue({
      code: z.ZodIssueCode.custom,
      message: 'Upper threshold must be greater than lower threshold',
      path: ['bottomDistanceHysteresis', 'upperThreshold'],
    })
  }
})

const queryClient = useQueryClient()
const { mutate, isPending } = useCargoConfigMutation()

const form = useForm({
  validationSchema: toTypedSchema(cargoConfigSchema),
  initialValues: props.initialValues,
})

const onSubmit = form.handleSubmit((values) => {
  mutate(values, {
    onSuccess: () => {
      queryClient.setQueryData([CARGO_CONFIG_QUERY_KEY], values)
      notification.success('Cargo configuration updated successfully!')
    },
    onError: () => {
      notification.error('Failed to update cargo configuration')
    },
  })
})
</script>

<template>
  <form class="flex flex-col w-full max-w-lg space-y-6" @submit="onSubmit">
    <h3 class="pb-2 text-lg font-medium border-b">
      Cargo Configuration
    </h3>

    <FormField v-slot="{ componentField }" name="liftPosition">
      <FormItem>
        <FormLabel>Lift Position</FormLabel>
        <FormControl>
          <Input
            v-bind="componentField"
            type="number"
            placeholder="Enter lift position"
            :disabled="isPending"
          />
        </FormControl>
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField v-slot="{ componentField }" name="lowerPosition">
      <FormItem>
        <FormLabel>Lower Position</FormLabel>
        <FormControl>
          <Input
            v-bind="componentField"
            type="number"
            placeholder="Enter lower position"
            :disabled="isPending"
          />
        </FormControl>
        <FormMessage />
      </FormItem>
    </FormField>
    <div class="flex flex-col gap-2">
      <h3 class="pb-2 text-lg font-medium border-b">
        Bottom Distance Hysteresis
      </h3>
    </div>

    <FormField v-slot="{ componentField }" name="bottomDistanceHysteresis.lowerThreshold">
      <FormItem>
        <FormLabel>Lower Threshold</FormLabel>
        <FormControl>
          <Input
            v-bind="componentField"
            type="number"
            placeholder="Enter lower threshold"
            :disabled="isPending"
          />
        </FormControl>
        <FormMessage />
      </FormItem>
    </FormField>
    <FormField v-slot="{ componentField }" name="bottomDistanceHysteresis.upperThreshold">
      <FormItem>
        <FormLabel>Upper Threshold</FormLabel>
        <FormControl>
          <Input
            v-bind="componentField"
            type="number"
            placeholder="Enter upper threshold"
            :disabled="isPending"
          />
        </FormControl>
        <FormMessage />
      </FormItem>
    </FormField>

    <div>
      <Button type="submit" :disabled="isPending">
        <Loader v-if="isPending" class="w-4 h-4 mr-2 animate-spin" />
        Save
      </Button>
    </div>
  </form>
</template>
