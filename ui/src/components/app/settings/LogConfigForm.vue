<script setup lang="ts">
import type { LogConfig } from '@/types/config'
import { Button } from '@/components/ui/button'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Switch } from '@/components/ui/switch'
import { LOG_CONFIG_QUERY_KEY, useLogConfigMutation } from '@/composables/use-config'
import { useQueryClient } from '@tanstack/vue-query'
import { toTypedSchema } from '@vee-validate/zod'
import { Loader } from 'lucide-vue-next'
import { useForm } from 'vee-validate'
import { z } from 'zod'

interface Props {
  initialValues: LogConfig
}
const props = defineProps<Props>()

const logConfigSchema = z.object({
  level: z.union([z.literal('DEBUG'), z.literal('INFO'), z.literal('WARN'), z.literal('ERROR')]),
  format: z.union([z.literal('JSON'), z.literal('TEXT')]),
  addSource: z.boolean(),
})

const queryClient = useQueryClient()
const { mutate, isPending } = useLogConfigMutation()

const form = useForm({
  validationSchema: toTypedSchema(logConfigSchema),
  initialValues: props.initialValues,
})

const onSubmit = form.handleSubmit((values) => {
  mutate(values, {
    onSuccess: () => {
      queryClient.setQueryData([LOG_CONFIG_QUERY_KEY], values)
      notification.success('Log configuration updated successfully!')
    },
    onError: () => {
      notification.error('Failed to update log configuration')
    },
  })
})
</script>

<template>
  <form class="flex flex-col w-full max-w-lg space-y-6" @submit="onSubmit">
    <h3 class="pb-2 text-lg font-medium border-b">
      Log Configuration
    </h3>

    <FormField v-slot="{ componentField }" name="level">
      <FormItem>
        <FormLabel>Log Level</FormLabel>
        <Select v-bind="componentField">
          <FormControl>
            <SelectTrigger :disabled="isPending">
              <SelectValue placeholder="Select log level" />
            </SelectTrigger>
          </FormControl>
          <SelectContent>
            <SelectItem value="DEBUG">
              Debug
            </SelectItem>
            <SelectItem value="INFO">
              Info
            </SelectItem>
            <SelectItem value="WARN">
              Warning
            </SelectItem>
            <SelectItem value="ERROR">
              Error
            </SelectItem>
          </SelectContent>
        </Select>
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField v-slot="{ componentField }" name="format">
      <FormItem>
        <FormLabel>Log Format</FormLabel>
        <Select v-bind="componentField">
          <FormControl>
            <SelectTrigger :disabled="isPending">
              <SelectValue placeholder="Select log format" />
            </SelectTrigger>
          </FormControl>
          <SelectContent>
            <SelectItem value="JSON">
              JSON
            </SelectItem>
            <SelectItem value="TEXT">
              Text
            </SelectItem>
          </SelectContent>
        </Select>
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField v-slot="{ value, handleChange }" type="checkbox" name="addSource">
      <FormItem class="flex flex-row items-center justify-between p-4 border rounded-lg">
        <div class="space-y-0.5">
          <FormLabel>Add Source</FormLabel>
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
