<script setup lang="ts">
import type { HTTPConfig } from '@/types/config'
import { Button } from '@/components/ui/button'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Switch } from '@/components/ui/switch'
import { HTTP_CONFIG_QUERY_KEY, useHTTPConfigMutation } from '@/composables/use-config'
import { useQueryClient } from '@tanstack/vue-query'
import { toTypedSchema } from '@vee-validate/zod'
import { Loader } from 'lucide-vue-next'
import { useForm } from 'vee-validate'
import { z } from 'zod'

interface Props {
  initialValues: HTTPConfig
}
const props = defineProps<Props>()

const httpConfigSchema = z.object({
  port: z.number().int().min(1024, 'Port must be at least 1024').max(65535, 'Port must be at most 65535'),
  swagger: z.boolean().default(false),
})

const queryClient = useQueryClient()
const { mutate, isPending } = useHTTPConfigMutation()

const form = useForm({
  validationSchema: toTypedSchema(httpConfigSchema),
  initialValues: props.initialValues,
})

const onSubmit = form.handleSubmit((values) => {
  mutate(values, {
    onSuccess: () => {
      queryClient.setQueryData([HTTP_CONFIG_QUERY_KEY], values)
      notification.success('HTTP configuration updated successfully!')
    },
    onError: () => {
      notification.error('Failed to update HTTP configuration')
    },
  })
})
</script>

<template>
  <form class="flex flex-col w-full max-w-lg space-y-6" @submit="onSubmit">
    <h3 class="pb-2 text-lg font-medium border-b">
      HTTP Configuration
    </h3>

    <FormField v-slot="{ field }" name="port">
      <FormItem>
        <FormLabel>Port</FormLabel>
        <FormControl>
          <Input
            v-model="field.value"
            type="number"
            placeholder="Enter port number"
            :disabled="isPending"
            class="[appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none"
          />
        </FormControl>
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField v-slot="{ value, handleChange }" name="swagger">
      <FormItem class="flex flex-row items-center justify-between p-4 border rounded-lg">
        <div class="space-y-0.5">
          <FormLabel>Enable Swagger</FormLabel>
        </div>
        <FormControl>
          <Switch
            :model-value="value"
            :disabled="isPending"
            aria-readonly
            @update:model-value="handleChange"
          />
        </FormControl>
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
