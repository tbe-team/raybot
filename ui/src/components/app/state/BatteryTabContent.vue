<script setup lang="ts">
import type { BatteryState, ChargeState, DischargeState } from '@/types/robot-state'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { BatteryFull, Thermometer, TriangleAlert, Zap } from 'lucide-vue-next'

const props = defineProps<{
  battery: BatteryState
  charge: ChargeState
  discharge: DischargeState
}>()
</script>

<template>
  <div class="space-y-4">
    <Card>
      <CardHeader>
        <CardTitle>Battery Details</CardTitle>
        <CardDescription>Current battery status and power management</CardDescription>
      </CardHeader>
      <CardContent>
        <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
          <div class="space-y-4">
            <div>
              <h3 class="mb-2 font-medium">
                Battery Status
              </h3>
              <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
                <div class="flex items-center gap-2">
                  <BatteryFull class="w-5 h-5 text-muted-foreground" />
                  <div>
                    <p class="text-sm font-medium">
                      Charge
                    </p>
                    <p class="text-xl font-bold">
                      {{ battery.percent }}%
                    </p>
                  </div>
                </div>
                <div class="flex items-center gap-2">
                  <Thermometer class="w-5 h-5 text-muted-foreground" />
                  <div>
                    <p class="text-sm font-medium">
                      Temperature
                    </p>
                    <p class="text-xl font-bold">
                      {{ battery.temp }}°C
                    </p>
                  </div>
                </div>
                <div class="flex items-center gap-2">
                  <Zap class="w-5 h-5 text-muted-foreground" />
                  <div>
                    <p class="text-sm font-medium">
                      Voltage
                    </p>
                    <p class="text-xl font-bold">
                      {{ battery.voltage }} V
                    </p>
                  </div>
                </div>
                <div class="flex items-center gap-2">
                  <Zap class="w-5 h-5 text-muted-foreground" />
                  <div>
                    <p class="text-sm font-medium">
                      Current
                    </p>
                    <p class="text-xl font-bold">
                      {{ battery.current }} mA
                    </p>
                  </div>
                </div>
              </div>
            </div>

            <div>
              <h3 class="mb-2 font-medium">
                Cell Voltages
              </h3>
              <div class="grid grid-cols-2 gap-4">
                <Card v-for="(voltage, index) in battery.cellVoltages" :key="index">
                  <CardContent class="p-2">
                    <p class="text-xs text-muted-foreground">
                      Cell {{ index + 1 }}
                    </p>
                    <p class="text-lg font-medium">
                      {{ voltage }} V
                    </p>
                  </CardContent>
                </Card>
              </div>
            </div>
          </div>

          <div class="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle class="text-base">
                  Charging Status
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div class="flex justify-between mb-4">
                  <span class="text-sm font-medium">Status</span>
                  <Badge :variant="charge.enabled ? 'default' : 'secondary'">
                    {{ charge.enabled ? 'Charging' : 'Not Charging' }}
                  </Badge>
                </div>
                <div class="flex justify-between mb-4">
                  <span class="text-sm font-medium">Health</span>
                  <span class="text-sm">{{ battery.health }}%</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-sm font-medium">Limit</span>
                  <span class="text-sm">{{ charge.currentLimit }} mA</span>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle class="text-base">
                  Discharge Status
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div class="flex justify-between mb-4">
                  <span class="text-sm font-medium">Status</span>
                  <Badge :variant="discharge.enabled ? 'default' : 'secondary'">
                    {{ discharge.enabled ? 'Discharging' : 'Not Discharging' }}
                  </Badge>
                </div>
                <div class="flex justify-between">
                  <span class="text-sm font-medium">Limit</span>
                  <span class="text-sm">{{ props.discharge.currentLimit }} mA</span>
                </div>
              </CardContent>
            </Card>

            <Card :class="battery.fault > 0 ? 'border-destructive' : 'border-success'">
              <CardHeader>
                <CardTitle class="text-base">
                  Fault Status
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div class="flex items-center gap-2">
                  <TriangleAlert class="w-5 h-5 text-destructive" />
                  <span class="text-sm font-medium" :class="props.battery.fault > 0 ? 'text-destructive' : 'text-success'">
                    {{ props.battery.fault > 0 ? 'Fault Detected' : 'No Faults' }}
                  </span>
                </div>
                <p v-if="props.battery.fault > 0" class="mt-2 text-sm text-destructive">
                  Fault code: {{ props.battery.fault }}
                </p>
              </CardContent>
            </Card>
          </div>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
