<script setup lang="ts">
import type { CommandConfig } from '@/types/command-config'
import { useQueryClient } from '@tanstack/vue-query'
import { toTypedSchema } from '@vee-validate/zod'
import { Loader } from 'lucide-vue-next'
import { useForm } from 'vee-validate'
import { z } from 'zod'
import { Button } from '@/components/ui/button'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { COMMAND_CONFIG_QUERY_KEY, useCommandConfigMutation } from '@/composables/use-config'

interface Props {
  initialValues: CommandConfig
}

const props = defineProps<Props>()

const commandConfigSchema = z.object({
  cargoLift: z.object({
    stableReadCount: z.number().int().positive('Stable read count must be positive').min(1),
  }),
  cargoLower: z.object({
    stableReadCount: z.number().int().positive('Stable read count must be positive').min(1),
    bottomObstacleTracking: z.object({
      enterDistance: z.number().int().min(1),
      exitDistance: z.number().int().min(1),
    }),
  }),
}).superRefine((data, ctx) => {
  if (data.cargoLower.bottomObstacleTracking.enterDistance >= data.cargoLower.bottomObstacleTracking.exitDistance) {
    ctx.addIssue({
      code: z.ZodIssueCode.custom,
      message: 'Enter distance must be less than exit distance',
      path: ['cargoLower.bottomObstacleTracking.enterDistance'],
    })
    ctx.addIssue({
      code: z.ZodIssueCode.custom,
      message: 'Enter distance must be less than exit distance',
      path: ['cargoLower.bottomObstacleTracking.exitDistance'],
    })
  }
})

const queryClient = useQueryClient()
const { mutate, isPending } = useCommandConfigMutation()
const form = useForm({
  validationSchema: toTypedSchema(commandConfigSchema),
  initialValues: props.initialValues,
})

const onSubmit = form.handleSubmit((values) => {
  mutate(values, {
    onSuccess: () => {
      notification.success('Command configuration updated successfully!')
      queryClient.setQueryData([COMMAND_CONFIG_QUERY_KEY], values)
    },
  })
})
</script>

<template>
  <form class="flex flex-col w-full max-w-lg space-y-6" @submit="onSubmit">
    <div class="grid grid-cols-1 gap-8">
      <!-- ESP Controller Section -->
      <div class="space-y-3">
        <h4 class="text-lg font-medium tracking-tight">
          Cargo Lift configuration
        </h4>

        <div class="space-y-6 ps-4">
          <FormField v-slot="{ componentField }" name="cargoLift.stableReadCount">
            <FormItem>
              <FormLabel>Stable Read Count</FormLabel>
              <FormControl>
                <Input v-bind="componentField" type="number" :disabled="isPending" />
              </FormControl>
              <FormMessage />
            </FormItem>
          </FormField>
        </div>
      </div>
    </div>
    <div class="grid grid-cols-1 gap-8">
      <!-- ESP Controller Section -->
      <div class="space-y-3">
        <h4 class="text-lg font-medium tracking-tight">
          Cargo Lower configuration
        </h4>

        <div class="space-y-6 ps-4">
          <FormField v-slot="{ componentField }" name="cargoLower.stableReadCount">
            <FormItem>
              <FormLabel>Stable Read Count</FormLabel>
              <FormControl>
                <Input v-bind="componentField" type="number" :disabled="isPending" />
              </FormControl>
              <FormMessage />
            </FormItem>
          </FormField>
          <div class="space-y-6">
            <h3 class="pb-2 text-lg font-medium border-b">
              Bottom Obstacle Tracking
            </h3>
            <FormField v-slot="{ componentField }" name="cargoLower.bottomObstacleTracking.enterDistance">
              <FormItem>
                <FormLabel>Enter Distance (cm)</FormLabel>
                <FormControl>
                  <Input v-bind="componentField" type="number" :disabled="isPending" />
                </FormControl>
                <FormMessage />
              </FormItem>
            </FormField>
            <FormField v-slot="{ componentField }" name="cargoLower.bottomObstacleTracking.exitDistance">
              <FormItem>
                <FormLabel>Exit Distance (cm)</FormLabel>
                <FormControl>
                  <Input v-bind="componentField" type="number" :disabled="isPending" />
                </FormControl>
                <FormMessage />
              </FormItem>
            </FormField>
          </div>
        </div>
      </div>
    </div>

    <div>
      <Button type="submit" :disabled="isPending">
        <Loader v-if="isPending" class="w-4 h-4 mr-2 animate-spin" />
        Save
      </Button>
    </div>
  </form>
</template>
