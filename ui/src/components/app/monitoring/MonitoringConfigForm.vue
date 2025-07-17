<script setup lang="ts">
import type { BatteryMonitoringConfig } from '@/types/config'
import { useQueryClient } from '@tanstack/vue-query'
import { toTypedSchema } from '@vee-validate/zod'
import { Loader } from 'lucide-vue-next'
import { useForm } from 'vee-validate'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Switch } from '@/components/ui/switch'
import { BATTERY_MONITORING_CONFIG_QUERY_KEY, useBatteryMonitoringConfigMutation } from '@/composables/use-config'
import { monitoringConfigSchema } from './schemas'

interface Props {
  initialValues: BatteryMonitoringConfig
}
const props = defineProps<Props>()

const queryClient = useQueryClient()
const { mutate, isPending } = useBatteryMonitoringConfigMutation()

const form = useForm({
  validationSchema: toTypedSchema(monitoringConfigSchema),
  initialValues: props.initialValues,
})

const onSubmit = form.handleSubmit((values) => {
  mutate(values, {
    onSuccess: () => {
      notification.success('Monitoring configuration updated successfully!')
      queryClient.invalidateQueries({ queryKey: [BATTERY_MONITORING_CONFIG_QUERY_KEY] })
    },
    onError: (err) => {
      notification.error(err.message)
    },
  })
})
</script>

<template>
  <form class="flex flex-col space-y-6 w-full" @submit="onSubmit">
    <div class="space-y-6">
      <!-- Battery Voltage Section -->
      <Card>
        <CardHeader>
          <CardTitle>Battery Voltage Monitoring</CardTitle>
        </CardHeader>
        <CardContent>
          <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
            <!-- Voltage Low -->
            <div>
              <div class="flex justify-between items-center">
                <h5 class="pb-2 text-base font-medium">
                  Voltage Low Alert
                </h5>
                <FormField v-slot="{ value, handleChange }" name="voltageLow.enable">
                  <FormItem>
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
              </div>

              <FormField v-slot="{ componentField }" name="voltageLow.threshold">
                <FormItem>
                  <FormLabel>Threshold (mV)</FormLabel>
                  <FormControl>
                    <Input
                      v-bind="componentField"
                      type="number"
                      :disabled="isPending || !form.values.voltageLow?.enable"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>
            </div>

            <!-- Voltage High -->
            <div>
              <div class="flex justify-between items-center">
                <h5 class="pb-2 text-base font-medium">
                  Voltage High Alert
                </h5>
                <FormField v-slot="{ value, handleChange }" name="voltageHigh.enable">
                  <FormItem>
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
              </div>

              <FormField v-slot="{ componentField }" name="voltageHigh.threshold">
                <FormItem>
                  <FormLabel>Threshold (mV)</FormLabel>
                  <FormControl>
                    <Input
                      v-bind="componentField"
                      type="number"
                      :disabled="isPending || !form.values.voltageHigh?.enable"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>
            </div>
          </div>
        </CardContent>
      </Card>

      <!-- Cell Voltage Section -->
      <Card>
        <CardHeader>
          <CardTitle>Cell Voltage Monitoring</CardTitle>
        </CardHeader>
        <CardContent>
          <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
            <!-- Cell Voltage High -->
            <div>
              <div class="flex justify-between items-center">
                <h5 class="pb-2 text-base font-medium">
                  Cell Voltage High Alert
                </h5>
                <FormField v-slot="{ value, handleChange }" name="cellVoltageHigh.enable">
                  <FormItem>
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
              </div>

              <FormField v-slot="{ componentField }" name="cellVoltageHigh.threshold">
                <FormItem>
                  <FormLabel>Threshold (mV)</FormLabel>
                  <FormControl>
                    <Input
                      v-bind="componentField"
                      type="number"
                      :disabled="isPending || !form.values.cellVoltageHigh?.enable"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>
            </div>

            <!-- Cell Voltage Low -->
            <div>
              <div class="flex justify-between items-center">
                <h5 class="pb-2 text-base font-medium">
                  Cell Voltage Low Alert
                </h5>
                <FormField v-slot="{ value, handleChange }" name="cellVoltageLow.enable">
                  <FormItem>
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
              </div>

              <FormField v-slot="{ componentField }" name="cellVoltageLow.threshold">
                <FormItem>
                  <FormLabel>Threshold (mV)</FormLabel>
                  <FormControl>
                    <Input
                      v-bind="componentField"
                      type="number"
                      :disabled="isPending || !form.values.cellVoltageLow?.enable"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>
            </div>

            <!-- Cell Voltage Diff -->
            <div>
              <div class="flex justify-between items-center">
                <h5 class="pb-2 text-base font-medium">
                  Cell Voltage Difference Alert
                </h5>
                <FormField v-slot="{ value, handleChange }" name="cellVoltageDiff.enable">
                  <FormItem>
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
              </div>

              <FormField v-slot="{ componentField }" name="cellVoltageDiff.threshold">
                <FormItem>
                  <FormLabel>Threshold (mV)</FormLabel>
                  <FormControl>
                    <Input
                      v-bind="componentField"
                      type="number"
                      :disabled="isPending || !form.values.cellVoltageDiff?.enable"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>
            </div>
          </div>
        </CardContent>
      </Card>

      <!-- Battery Status Section -->
      <Card>
        <CardHeader>
          <CardTitle>Battery Status Monitoring</CardTitle>
        </CardHeader>
        <CardContent>
          <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
            <!-- Current High -->
            <div>
              <div class="flex justify-between items-center">
                <h5 class="pb-2 text-base font-medium">
                  Current High Alert
                </h5>
                <FormField v-slot="{ value, handleChange }" name="currentHigh.enable">
                  <FormItem>
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
              </div>

              <FormField v-slot="{ componentField }" name="currentHigh.threshold">
                <FormItem>
                  <FormLabel>Threshold (mA)</FormLabel>
                  <FormControl>
                    <Input
                      v-bind="componentField"
                      type="number"
                      :disabled="isPending || !form.values.currentHigh?.enable"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>
            </div>

            <!-- Temperature High -->
            <div>
              <div class="flex justify-between items-center">
                <h5 class="pb-2 text-base font-medium">
                  Temperature High Alert
                </h5>
                <FormField v-slot="{ value, handleChange }" name="tempHigh.enable">
                  <FormItem>
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
              </div>

              <FormField v-slot="{ componentField }" name="tempHigh.threshold">
                <FormItem>
                  <FormLabel>Threshold (°C)</FormLabel>
                  <FormControl>
                    <Input
                      v-bind="componentField"
                      type="number"
                      :disabled="isPending || !form.values.tempHigh?.enable"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>
            </div>

            <!-- Battery Percent Low -->
            <div>
              <div class="flex justify-between items-center">
                <h5 class="pb-2 text-base font-medium">
                  Battery Percent Low Alert
                </h5>
                <FormField v-slot="{ value, handleChange }" name="percentLow.enable">
                  <FormItem>
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
              </div>

              <FormField v-slot="{ componentField }" name="percentLow.threshold">
                <FormItem>
                  <FormLabel>Threshold (%)</FormLabel>
                  <FormControl>
                    <Input
                      v-bind="componentField"
                      type="number"
                      :disabled="isPending || !form.values.percentLow?.enable"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>
            </div>

            <!-- Battery Health Low -->
            <div>
              <div class="flex justify-between items-center">
                <h5 class="pb-2 text-base font-medium">
                  Battery Health Low Alert
                </h5>
                <FormField v-slot="{ value, handleChange }" name="healthLow.enable">
                  <FormItem>
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
              </div>

              <FormField v-slot="{ componentField }" name="healthLow.threshold">
                <FormItem>
                  <FormLabel>Threshold (%)</FormLabel>
                  <FormControl>
                    <Input
                      v-bind="componentField"
                      type="number"
                      :disabled="isPending || !form.values.healthLow?.enable"
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              </FormField>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>

    <div>
      <Button type="submit" :disabled="isPending">
        <Loader v-if="isPending" class="mr-2 w-4 h-4 animate-spin" />
        Save
      </Button>
    </div>
  </form>
</template>
