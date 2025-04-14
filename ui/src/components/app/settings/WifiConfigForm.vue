<script setup lang="ts">
import type { WifiConfig } from '@/types/config'
import { Button } from '@/components/ui/button'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input, PasswordInput } from '@/components/ui/input'
import { Switch } from '@/components/ui/switch'
import { useWifiConfigMutation, WIFI_CONFIG_QUERY_KEY } from '@/composables/use-config'
import { useQueryClient } from '@tanstack/vue-query'
import { toTypedSchema } from '@vee-validate/zod'
import { Loader } from 'lucide-vue-next'
import { useForm } from 'vee-validate'
import { z } from 'zod'

interface Props {
  initialValues: WifiConfig
}
const props = defineProps<Props>()

const ssidRegex = /^[\w\s\-.]*$/ // alphanumeric, space, -, _, .
const passwordRegex = /^[\x21-\x7E]*$/ // printable characters
const ipv4Regex = /^(?:(?:25[0-5]|2[0-4]\d|[01]?\d{1,2})\.){3}(?:25[0-5]|2[0-4]\d|[01]?\d{1,2})$/

const wifiConfigSchema = z.object({
  ap: z.object({
    enable: z.boolean(),
    ssid: z.string()
      .min(1, 'SSID is required')
      .max(32, 'SSID must be at most 32 characters')
      .regex(ssidRegex, 'SSID can only contain alphanumeric characters, spaces, hyphens, underscores, and dots'),
    password: z.string()
      .min(8, 'Password must be at least 8 characters')
      .max(63, 'Password must be at most 63 characters')
      .regex(passwordRegex, 'Password can only contain printable characters'),
    ip: z.string()
      .min(1, 'IP address is required')
      .regex(ipv4Regex, 'Invalid IPv4 address format'),
  }),
  sta: z.object({
    enable: z.boolean(),
    ssid: z.string()
      .min(1, 'SSID is required')
      .max(32, 'SSID must be at most 32 characters')
      .regex(ssidRegex, 'SSID can only contain alphanumeric characters, spaces, hyphens, underscores, and dots'),
    password: z.string()
      .min(8, 'Password must be at least 8 characters')
      .max(63, 'Password must be at most 63 characters')
      .regex(passwordRegex, 'Password can only contain printable characters'),
  }),
}).superRefine((data, ctx) => {
  if ((data.ap.enable && data.sta.enable) || (!data.ap.enable && !data.sta.enable)) {
    ctx.addIssue({
      code: z.ZodIssueCode.custom,
      message: 'AP and STA cannot be enabled or disabled at the same time',
      path: ['ap.enable'],
    })
    ctx.addIssue({
      code: z.ZodIssueCode.custom,
      message: 'AP and STA cannot be enabled or disabled at the same time',
      path: ['sta.enable'],
    })
  }
})

const queryClient = useQueryClient()
const { mutate, isPending } = useWifiConfigMutation()

const form = useForm({
  validationSchema: toTypedSchema(wifiConfigSchema),
  initialValues: props.initialValues,
})

const onSubmit = form.handleSubmit((values) => {
  mutate(values, {
    onSuccess: () => {
      queryClient.setQueryData([WIFI_CONFIG_QUERY_KEY], values)
      notification.success('WiFi configuration updated successfully!')
    },
    onError: () => {
      notification.error('Failed to update WiFi configuration')
    },
  })
})
</script>

<template>
  <form class="flex flex-col w-full max-w-lg space-y-6" @submit="onSubmit">
    <h3 class="pb-2 text-lg font-medium border-b">
      AP Configuration
    </h3>

    <FormField v-slot="{ value, handleChange }" type="checkbox" name="ap.enable">
      <FormItem class="flex flex-row items-center justify-between p-4 border rounded-lg">
        <div class="space-y-0.5">
          <FormLabel>Enable AP</FormLabel>
        </div>
        <FormControl>
          <Switch
            :model-value="value"
            :disabled="isPending"
            @update:model-value="handleChange"
          />
        </FormControl>
      </FormItem>
      <FormMessage />
    </FormField>

    <FormField v-slot="{ componentField }" name="ap.ssid">
      <FormItem>
        <FormLabel>SSID</FormLabel>
        <FormControl>
          <Input
            v-bind="componentField"
            placeholder="Enter SSID"
            :disabled="isPending"
          />
        </FormControl>
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField v-slot="{ componentField }" name="ap.password">
      <FormItem>
        <FormLabel>Password</FormLabel>
        <FormControl>
          <PasswordInput
            v-bind="componentField"
            placeholder="Enter password"
            :disabled="isPending"
          />
        </FormControl>
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField v-slot="{ componentField }" name="ap.ip">
      <FormItem>
        <FormLabel>IP Address</FormLabel>
        <FormControl>
          <Input
            v-bind="componentField"
            placeholder="Enter IP address"
            :disabled="isPending"
          />
        </FormControl>
        <FormMessage />
      </FormItem>
    </FormField>

    <h3 class="pb-2 text-lg font-medium border-b">
      STA Configuration
    </h3>

    <FormField v-slot="{ value, handleChange }" type="checkbox" name="sta.enable">
      <FormItem class="flex flex-row items-center justify-between p-4 border rounded-lg">
        <div class="space-y-0.5">
          <FormLabel>Enable STA</FormLabel>
        </div>
        <FormControl>
          <Switch
            :model-value="value"
            :disabled="isPending"
            @update:model-value="handleChange"
          />
        </FormControl>
      </FormItem>
      <FormMessage />
    </FormField>

    <FormField v-slot="{ componentField }" name="sta.ssid">
      <FormItem>
        <FormLabel>SSID</FormLabel>
        <FormControl>
          <Input
            v-bind="componentField"
            placeholder="Enter SSID"
            :disabled="isPending"
          />
        </FormControl>
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField v-slot="{ componentField }" name="sta.password">
      <FormItem>
        <FormLabel>Password</FormLabel>
        <FormControl>
          <PasswordInput
            v-bind="componentField"
            placeholder="Enter password"
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
