<script setup lang="ts">
import PageContainer from '@/components/shared/PageContainer.vue'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { useQueryRobotState } from '@/composables/use-robot-state'
import { formatDate } from '@/lib/date'
import { AlertCircle, Loader } from 'lucide-vue-next'

const REFRESH_INTERVAL = 1000

const { data: robotState, isPending, isError, error } = useQueryRobotState({
  axiosOpts: { doNotShowLoading: true },
  refetchInterval: REFRESH_INTERVAL,
})

function getBatteryColor(percent: number): string {
  if (percent < 20)
    return 'text-red-5500'
  if (percent < 40)
    return 'text-orange-500'
  if (percent < 60)
    return 'text-yellow-500'
  return 'text-green-500'
}
function getTemperatureColor(temp: number): string {
  if (temp < 15)
    return 'text-blue-500'
  if (temp < 35)
    return 'text-green-500'
  return 'text-red-500'
}
</script>

<template>
  <PageContainer>
    <div v-if="isPending" class="flex flex-col items-center justify-center gap-4 pt-20">
      <div class="flex items-center gap-4">
        <Loader class="w-8 h-8 animate-spin text-muted-foreground" />
      </div>
      <p class="text-lg text-muted-foreground">
        Loading state...
      </p>
    </div>
    <div v-else-if="isError" class="flex flex-col items-center justify-center gap-4 pt-20">
      <Card class="flex flex-col items-center gap-4 p-6 text-destructive">
        <AlertCircle class="w-8 h-8" />
        <div class="space-y-2 text-center">
          <h2 class="text-lg font-semibold">
            Failed to load state
          </h2>
          <p class="text-sm text-muted-foreground">
            {{ error?.message || 'An unexpected error occurred' }}
          </p>
        </div>
      </Card>
    </div>
    <div v-else-if="!robotState" class="flex flex-col items-center justify-center gap-4 pt-20">
      <Card class="flex flex-col items-center gap-4 p-6">
        <AlertCircle class="w-8 h-8 text-muted-foreground" />
        <div class="space-y-2 text-center">
          <h2 class="text-lg font-semibold">
            Robot state not found
          </h2>
          <p class="text-sm text-muted-foreground">
            The robot state appears to be empty
          </p>
        </div>
      </Card>
    </div>
    <div v-else class="flex flex-col w-full">
      <div class="mb-6">
        <h1 class="text-xl font-semibold">
          Robot state
        </h1>
        <p class="text-sm text-muted-foreground">
          The current state of the robot is updated once per second.
        </p>
      </div>
      <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
        <Card class="col-span-1 sm:col-span-2">
          <CardHeader>
            <CardTitle>Battery State</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="mx-auto space-x-2">
              <div class="grid grid-cols-1 gap-2 sm:grid-cols-2">
                <p><span class="font-medium">Current: </span>{{ robotState.battery.current }} mA</p>
                <p><span class="font-medium">Battery level: </span><span :class="getBatteryColor(robotState.battery.percent)">{{ robotState.battery.percent }}%</span></p>
                <p><span class="font-medium">Voltage: </span>{{ robotState.battery.voltage }} V</p>
                <p><span class="font-medium">Health: </span>{{ robotState.battery.health }}%</p>
                <p><span class="font-medium">Cell voltages: </span><span>{{ robotState.battery.cellVoltages.join(', ') }} V</span></p>
                <p><span class="font-medium">Fault status: </span><span class="text-green-600">{{ robotState.battery.fault }}</span></p>
                <p><span class="font-medium">Temperature: </span><span :class="getTemperatureColor(robotState.battery.temp)">{{ robotState.battery.temp }}°C</span></p>
                <p><span class="font-medium">Last updated: </span>{{ formatDate(robotState.battery.updatedAt) }}</p>
              </div>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>Charge State</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="mx-auto space-x-2">
              <div class="flex flex-col gap-2">
                <p>
                  <span class="font-medium">Current limit: </span>{{ robotState.charge.currentLimit }} mA
                </p>
                <p>
                  <span class="font-medium">Charging: </span>
                  <span :class="robotState.charge.enabled ? 'text-green-500' : 'text-red-500'">{{ robotState.charge.enabled ? 'Yes' : 'No' }}</span>
                </p>
                <p>
                  <span class="font-medium">Last updated: </span> {{ formatDate(robotState.charge.updatedAt) }}
                </p>
              </div>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>Discharge State</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="mx-auto space-x-2">
              <div class="flex flex-col gap-2">
                <p>
                  <span class="font-medium">Current limit: </span>{{ robotState.discharge.currentLimit }} mA
                </p>
                <p>
                  <span class="font-medium">Discharging: </span>
                  <span :class="robotState.discharge.enabled ? 'text-green-500' : 'text-red-500'">{{ robotState.discharge.enabled ? 'Yes' : 'No' }}</span>
                </p>
                <p>
                  <span class="font-medium">Last updated: </span>{{ formatDate(robotState.discharge.updatedAt) }}
                </p>
              </div>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>Lift Motor State</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="mx-auto space-x-2">
              <div class="grid grid-rows-1 gap-2">
                <div class="grid grid-cols-1 gap-2 sm:grid-cols-2">
                  <p><span class="font-medium">Current: </span>{{ robotState.liftMotor.currentPosition }} cm</p>
                  <p><span class="font-medium">Target: </span>{{ robotState.liftMotor.targetPosition }} cm</p>
                  <p>
                    <span class="font-medium">Running: </span>
                    <span :class="robotState.liftMotor.isRunning ? 'text-green-500' : 'text-red-500'">{{ robotState.liftMotor.isRunning ? 'Yes' : 'No' }}</span>
                  </p>
                  <p>
                    <span class="font-medium">Enabled: </span>
                    <span :class="robotState.liftMotor.enabled ? 'text-green-500' : 'text-red-500'">{{ robotState.liftMotor.enabled ? 'Yes' : 'No' }}</span>
                  </p>
                </div>
                <p><span class="font-medium">Last updated:</span> {{ formatDate(robotState.liftMotor.updatedAt) }}</p>
              </div>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>Drive Motor State</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="mx-auto space-x-2">
              <div class="grid grid-rows-1 gap-2">
                <p><span class="font-medium">Direction: </span>{{ robotState.driveMotor.direction }}</p>
                <p><span class="font-medium">Speed: </span>{{ robotState.driveMotor.speed }} %</p>
                <div class="grid grid-cols-1 gap-2 sm:grid-cols-2">
                  <p>
                    <span class="font-medium">Running: </span>
                    <span :class="robotState.driveMotor.isRunning ? 'text-green-500' : 'text-red-500'">
                      {{ robotState.driveMotor.isRunning ? 'Yes' : 'No' }}
                    </span>
                  </p>
                  <p>
                    <span class="font-medium">Enabled: </span>
                    <span :class="robotState.driveMotor.enabled ? 'text-green-500' : 'text-red-500'">
                      {{ robotState.driveMotor.enabled ? 'Yes' : 'No' }}
                    </span>
                  </p>
                </div>
                <p><span class="font-medium">Last updated: </span>{{ formatDate(robotState.driveMotor.updatedAt) }}</p>
              </div>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>Distance Sensor State</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="mx-auto space-x-2">
              <div class="grid grid-rows-2 gap-2">
                <p><span class="font-medium">Front: </span>{{ robotState.distanceSensor.frontDistance }} cm</p>
                <p><span class="font-medium">Back: </span>{{ robotState.distanceSensor.backDistance }} cm</p>
                <p><span class="font-medium">Down: </span>{{ robotState.distanceSensor.downDistance }} cm</p>
                <p><span class="font-medium">Last updated: </span>{{ formatDate(robotState.distanceSensor.updatedAt) }}</p>
              </div>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>Location State</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="mx-auto space-x-2">
              <div class="grid grid-rows-1 gap-2">
                <p><span class="font-medium">Current location: </span>{{ robotState.location.currentLocation }}</p>
                <p><span class="font-medium">Last updated:</span> {{ formatDate(robotState.location.updatedAt) }}</p>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Cargo Door Motor State</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="mx-auto space-x-2">
              <div class="grid grid-rows-1 gap-2">
                <p><span class="font-medium">Direction: </span>{{ robotState.cargoDoorMotor.direction }}</p>
                <p><span class="font-medium">Speed: </span>{{ robotState.cargoDoorMotor.speed }} %</p>
                <div class="grid grid-cols-1 gap-2 sm:grid-cols-2">
                  <p>
                    <span class="font-medium">Running: </span>
                    <span :class="robotState.cargoDoorMotor.isRunning ? 'text-green-500' : 'text-red-500'">
                      {{ robotState.cargoDoorMotor.isRunning ? 'Yes' : 'No' }}
                    </span>
                  </p>
                  <p>
                    <span class="font-medium">Enabled: </span>
                    <span :class="robotState.cargoDoorMotor.enabled ? 'text-green-500' : 'text-red-500'">
                      {{ robotState.cargoDoorMotor.enabled ? 'Yes' : 'No' }}
                    </span>
                  </p>
                </div>
                <p><span class="font-medium">Last updated:</span> {{ formatDate(robotState.cargoDoorMotor.updatedAt) }}</p>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Cargo State</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="mx-auto space-x-2">
              <div class="grid grid-rows-1 gap-2">
                <p>
                  <span class="font-medium">Door Open: </span>
                  <span :class="robotState.cargo.isOpen ? 'text-green-500' : 'text-red-500'">
                    {{ robotState.cargo.isOpen ? 'Yes' : 'No' }}
                  </span>
                </p>
                <p><span class="font-medium">QR Code: </span>{{ robotState.cargo.qrCode || 'None' }}</p>
                <p><span class="font-medium">Bottom Distance: </span>{{ robotState.cargo.bottomDistance }} cm</p>
                <p><span class="font-medium">Last updated:</span> {{ formatDate(robotState.cargo.updatedAt) }}</p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  </PageContainer>
</template>
