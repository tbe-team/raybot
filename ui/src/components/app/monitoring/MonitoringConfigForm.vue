<script setup lang="ts">
import type { MonitoringConfig } from '@/types/config'
import { useQueryClient } from '@tanstack/vue-query'
import { toTypedSchema } from '@vee-validate/zod'
import { Loader } from 'lucide-vue-next'
import { useForm } from 'vee-validate'
import { z } from 'zod'
import { Button } from '@/components/ui/button'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Switch } from '@/components/ui/switch'
import { useMonitoringConfigMutation } from '@/composables/use-config'

interface Props {
  initialValues: MonitoringConfig
}
const props = defineProps<Props>()

const monitoringConfigSchema = z.object({
  voltageLow: z.object({
    enable: z.boolean(),
    threshold: z.number().min(0, 'Threshold must be positive'),
  }),
  voltageHigh: z.object({
    enable: z.boolean(),
    threshold: z.number().min(0, 'Threshold must be positive'),
  }),
  cellVoltageHigh: z.object({
    enable: z.boolean(),
    threshold: z.number().min(0, 'Threshold must be positive'),
  }),
  cellVoltageLow: z.object({
    enable: z.boolean(),
    threshold: z.number().min(0, 'Threshold must be positive'),
  }),
  cellVoltageDiff: z.object({
    enable: z.boolean(),
    threshold: z.number().min(0, 'Threshold must be positive'),
  }),
  currentHigh: z.object({
    enable: z.boolean(),
    threshold: z.number().min(0, 'Threshold must be positive'),
  }),
  tempHigh: z.object({
    enable: z.boolean(),
    threshold: z.number().min(0, 'Threshold must be positive'),
  }),
  percentLow: z.object({
    enable: z.boolean(),
    threshold: z.number().min(0).max(100, 'Percentage must be between 0 and 100'),
  }),
  healthLow: z.object({
    enable: z.boolean(),
    threshold: z.number().min(0).max(100, 'Health must be between 0 and 100'),
  }),
})

const queryClient = useQueryClient()
const { mutate, isPending } = useMonitoringConfigMutation()

const form = useForm({
  validationSchema: toTypedSchema(monitoringConfigSchema),
  initialValues: props.initialValues,
})

const onSubmit = form.handleSubmit((values) => {
  mutate(values, {
    onSuccess: () => {
      queryClient.setQueryData(['config', 'monitoring'], values)
      notification.success('Monitoring configuration updated successfully!')
    },
    onError: () => {
      notification.error('Failed to update monitoring configuration')
    },
  })
})
</script>

<template>
  <form class="flex flex-col space-y-6 w-full" @submit="onSubmit">
    <div class="grid grid-cols-1 gap-8">
      <!-- Battery Voltage Section -->
      <div class="space-y-3">
        <h4 class="text-lg font-medium tracking-tight">
          Battery Voltage Monitoring
        </h4>

        <div class="px-4 space-y-6">
          <!-- Voltage Low -->
          <div class="space-y-4">
            <h5 class="pb-2 text-base font-medium border-b">
              Voltage Low Alert
            </h5>
            <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
              <FormField v-slot="{ value, handleChange }" name="voltageLow.enable">
                <FormItem class="flex flex-row justify-between items-center p-4 rounded-lg border">
                  <FormLabel>Enable Alert</FormLabel>
                  <FormControl>
                    <Switch
                      :model-value="value"
                      :disabled="isPending"
                      @update:model-value="handleChange"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>

              <FormField v-slot="{ componentField }" name="voltageLow.threshold">
                <FormItem>
                  <FormLabel>Threshold (V)</FormLabel>
                  <FormControl>
                    <Input
                      v-bind="componentField"
                      type="number"
                      step="0.1"
                      placeholder="14.0"
                      :disabled="isPending"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>
            </div>
          </div>

          <!-- Voltage High -->
          <div class="space-y-4">
            <h5 class="pb-2 text-base font-medium border-b">
              Voltage High Alert
            </h5>
            <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
              <FormField v-slot="{ value, handleChange }" name="voltageHigh.enable">
                <FormItem class="flex flex-row justify-between items-center p-4 rounded-lg border">
                  <FormLabel>Enable Alert</FormLabel>
                  <FormControl>
                    <Switch
                      :model-value="value"
                      :disabled="isPending"
                      @update:model-value="handleChange"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>

              <FormField v-slot="{ componentField }" name="voltageHigh.threshold">
                <FormItem>
                  <FormLabel>Threshold (V)</FormLabel>
                  <FormControl>
                    <Input
                      v-bind="componentField"
                      type="number"
                      step="0.1"
                      placeholder="18.0"
                      :disabled="isPending"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>
            </div>
          </div>
        </div>
      </div>

      <!-- Cell Voltage Section -->
      <div class="space-y-3">
        <h4 class="text-lg font-medium tracking-tight">
          Cell Voltage Monitoring
        </h4>

        <div class="px-4 space-y-6">
          <!-- Cell Voltage High -->
          <div class="space-y-4">
            <h5 class="pb-2 text-base font-medium border-b">
              Cell Voltage High Alert
            </h5>
            <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
              <FormField v-slot="{ value, handleChange }" name="cellVoltageHigh.enable">
                <FormItem class="flex flex-row justify-between items-center p-4 rounded-lg border">
                  <FormLabel>Enable Alert</FormLabel>
                  <FormControl>
                    <Switch
                      :model-value="value"
                      :disabled="isPending"
                      @update:model-value="handleChange"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>

              <FormField v-slot="{ componentField }" name="cellVoltageHigh.threshold">
                <FormItem>
                  <FormLabel>Threshold (V)</FormLabel>
                  <FormControl>
                    <Input
                      v-bind="componentField"
                      type="number"
                      step="0.01"
                      placeholder="4.30"
                      :disabled="isPending"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>
            </div>
          </div>

          <!-- Cell Voltage Low -->
          <div class="space-y-4">
            <h5 class="pb-2 text-base font-medium border-b">
              Cell Voltage Low Alert
            </h5>
            <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
              <FormField v-slot="{ value, handleChange }" name="cellVoltageLow.enable">
                <FormItem class="flex flex-row justify-between items-center p-4 rounded-lg border">
                  <FormLabel>Enable Alert</FormLabel>
                  <FormControl>
                    <Switch
                      :model-value="value"
                      :disabled="isPending"
                      @update:model-value="handleChange"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>

              <FormField v-slot="{ componentField }" name="cellVoltageLow.threshold">
                <FormItem>
                  <FormLabel>Threshold (V)</FormLabel>
                  <FormControl>
                    <Input
                      v-bind="componentField"
                      type="number"
                      step="0.01"
                      placeholder="3.80"
                      :disabled="isPending"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>
            </div>
          </div>

          <!-- Cell Voltage Diff -->
          <div class="space-y-4">
            <h5 class="pb-2 text-base font-medium border-b">
              Cell Voltage Difference Alert
            </h5>
            <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
              <FormField v-slot="{ value, handleChange }" name="cellVoltageDiff.enable">
                <FormItem class="flex flex-row justify-between items-center p-4 rounded-lg border">
                  <FormLabel>Enable Alert</FormLabel>
                  <FormControl>
                    <Switch
                      :model-value="value"
                      :disabled="isPending"
                      @update:model-value="handleChange"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>

              <FormField v-slot="{ componentField }" name="cellVoltageDiff.threshold">
                <FormItem>
                  <FormLabel>Threshold (V)</FormLabel>
                  <FormControl>
                    <Input
                      v-bind="componentField"
                      type="number"
                      step="0.01"
                      placeholder="0.50"
                      :disabled="isPending"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>
            </div>
          </div>
        </div>
      </div>

      <!-- Battery Status Section -->
      <div class="space-y-3">
        <h4 class="text-lg font-medium tracking-tight">
          Battery Status Monitoring
        </h4>

        <div class="px-4 space-y-6">
          <!-- Current High -->
          <div class="space-y-4">
            <h5 class="pb-2 text-base font-medium border-b">
              Current High Alert
            </h5>
            <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
              <FormField v-slot="{ value, handleChange }" name="currentHigh.enable">
                <FormItem class="flex flex-row justify-between items-center p-4 rounded-lg border">
                  <FormLabel>Enable Alert</FormLabel>
                  <FormControl>
                    <Switch
                      :model-value="value"
                      :disabled="isPending"
                      @update:model-value="handleChange"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>

              <FormField v-slot="{ componentField }" name="currentHigh.threshold">
                <FormItem>
                  <FormLabel>Threshold (A)</FormLabel>
                  <FormControl>
                    <Input
                      v-bind="componentField"
                      type="number"
                      step="0.1"
                      placeholder="6.0"
                      :disabled="isPending"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>
            </div>
          </div>

          <!-- Temperature High -->
          <div class="space-y-4">
            <h5 class="pb-2 text-base font-medium border-b">
              Temperature High Alert
            </h5>
            <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
              <FormField v-slot="{ value, handleChange }" name="tempHigh.enable">
                <FormItem class="flex flex-row justify-between items-center p-4 rounded-lg border">
                  <FormLabel>Enable Alert</FormLabel>
                  <FormControl>
                    <Switch
                      :model-value="value"
                      :disabled="isPending"
                      @update:model-value="handleChange"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>

              <FormField v-slot="{ componentField }" name="tempHigh.threshold">
                <FormItem>
                  <FormLabel>Threshold (°C)</FormLabel>
                  <FormControl>
                    <Input
                      v-bind="componentField"
                      type="number"
                      step="1"
                      placeholder="60"
                      :disabled="isPending"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>
            </div>
          </div>

          <!-- Battery Percent Low -->
          <div class="space-y-4">
            <h5 class="pb-2 text-base font-medium border-b">
              Battery Percent Low Alert
            </h5>
            <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
              <FormField v-slot="{ value, handleChange }" name="percentLow.enable">
                <FormItem class="flex flex-row justify-between items-center p-4 rounded-lg border">
                  <FormLabel>Enable Alert</FormLabel>
                  <FormControl>
                    <Switch
                      :model-value="value"
                      :disabled="isPending"
                      @update:model-value="handleChange"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>

              <FormField v-slot="{ componentField }" name="percentLow.threshold">
                <FormItem>
                  <FormLabel>Threshold (%)</FormLabel>
                  <FormControl>
                    <Input
                      v-bind="componentField"
                      type="number"
                      step="1"
                      min="0"
                      max="100"
                      placeholder="20"
                      :disabled="isPending"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>
            </div>
          </div>

          <!-- Battery Health Low -->
          <div class="space-y-4">
            <h5 class="pb-2 text-base font-medium border-b">
              Battery Health Low Alert
            </h5>
            <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
              <FormField v-slot="{ value, handleChange }" name="healthLow.enable">
                <FormItem class="flex flex-row justify-between items-center p-4 rounded-lg border">
                  <FormLabel>Enable Alert</FormLabel>
                  <FormControl>
                    <Switch
                      :model-value="value"
                      :disabled="isPending"
                      @update:model-value="handleChange"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>

              <FormField v-slot="{ componentField }" name="healthLow.threshold">
                <FormItem>
                  <FormLabel>Threshold (%)</FormLabel>
                  <FormControl>
                    <Input
                      v-bind="componentField"
                      type="number"
                      step="1"
                      min="0"
                      max="100"
                      placeholder="60"
                      :disabled="isPending"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div>
      <Button type="submit" :disabled="isPending">
        <Loader v-if="isPending" class="mr-2 w-4 h-4 animate-spin" />
        Save
      </Button>
    </div>
  </form>
</template>
