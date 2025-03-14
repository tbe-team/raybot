<script setup lang="ts">
import PageContainer from '@/components/shared/PageContainer.vue'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { useQueryRobotState } from '@/composables/use-robot-state'
import { AlertCircle, Loader } from 'lucide-vue-next'

const { data: robotState, isPending, isError, error } = useQueryRobotState({ doNotShowLoading: true }, 1000)

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
            No State Found
          </h2>
          <p class="text-sm text-muted-foreground">
            The state appears to be empty
          </p>
        </div>
      </Card>
    </div>
    <div v-else class="flex flex-col w-full">
      <!-- <StatePage :robot-state="robotState" /> -->
      <CardHeader>
        <div class="space-y-2">
          <CardTitle>
            Robot state
          </CardTitle>
          <CardDescription>
            The current state of the robot is updated once per second.
          </CardDescription>
        </div>
      </CardHeader>
      <CardContent class="grid grid-cols-4 gap-4">
        <Card class="h-full col-span-2 rounded-sm shadow-lg">
          <CardHeader>
            <CardTitle>Battery State</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="mx-auto space-x-2">
              <div class="grid grid-cols-2 gap-2">
                <p><span class="font-medium">Current: </span>{{ robotState.battery.current }} mA</p>
                <p><span class="font-medium">Battery level: </span><span :class="getBatteryColor(robotState.battery.percent)">{{ robotState.battery.percent }}%</span></p>
                <p><span class="font-medium">Voltage: </span>{{ robotState.battery.voltage }} V</p>
                <p><span class="font-medium">Health: </span>{{ robotState.battery.health }}%</p>
                <p><span class="font-medium">Cell voltages: </span><span>{{ robotState.battery.cellVoltages.join(', ') }} V</span></p>
                <p><span class="font-medium">Fault status: </span><span class="text-green-600">{{ robotState.battery.fault }}</span></p>
                <p><span class="font-medium">Temperature: </span><span :class="getTemperatureColor(robotState.battery.temp)">{{ robotState.battery.temp }}°C</span></p>
                <p><span class="font-medium">Last updated: </span>{{ robotState.battery.updatedAt }}</p>
              </div>
            </div>
          </CardContent>
        </Card>
        <Card class="h-full rounded-sm shadow-lg">
          <CardHeader>
            <CardTitle>Charge State</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="mx-auto space-x-2">
              <div class="flex flex-col gap-2 ">
                <p>
                  <span class="font-medium">Current limit: </span>{{ robotState.charge.currentLimit }} mA
                </p>
                <p>
                  <span class="font-medium">Charging: </span>
                  <span :class="robotState.charge.enabled ? 'text-green-500' : 'text-red-500'">{{ robotState.charge.enabled ? 'Yes' : 'No' }}</span>
                </p>
                <p>
                  <span class="font-medium">Last updated: </span> {{ robotState.charge.updatedAt }}
                </p>
              </div>
            </div>
          </CardContent>
        </Card>
        <Card class="h-full rounded-sm shadow-lg ">
          <CardHeader>
            <CardTitle>Discharge State</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="mx-auto space-x-2">
              <div class="flex flex-col gap-2 ">
                <p>
                  <span class="font-medium">Current limit: </span>{{ robotState.discharge.currentLimit }} mA
                </p>
                <p>
                  <span class="font-medium">Discharging: </span>
                  <span :class="robotState.discharge.enabled ? 'text-green-500' : 'text-red-500'">{{ robotState.discharge.enabled ? 'Yes' : 'No' }}</span>
                </p>
                <p>
                  <span class="font-medium">Last updated: </span>{{ robotState.discharge.updatedAt }}
                </p>
              </div>
            </div>
          </CardContent>
        </Card>
        <Card class="h-full rounded-sm shadow-lg">
          <CardHeader>
            <CardTitle>Lift Motor State</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="mx-auto space-x-2">
              <div class="grid grid-rows-1 gap-2 ">
                <div class="grid grid-cols-2 gap-2">
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
                <p><span class="font-medium">Last updated:</span> {{ robotState.liftMotor.updatedAt }}</p>
              </div>
            </div>
          </CardContent>
        </Card>
        <Card class="h-full rounded-sm shadow-lg ">
          <CardHeader>
            <CardTitle>Drive Motor State</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="mx-auto space-x-2">
              <div class="grid grid-rows-1 gap-2 ">
                <p><span class="font-medium">Direction: </span>{{ robotState.driveMotor.direction }}</p>
                <p><span class="font-medium">Speed: </span>{{ robotState.driveMotor.speed }} %</p>
                <div class="grid grid-cols-2 gap-2">
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
                <p><span class="font-medium">Last updated: </span>{{ robotState.driveMotor.updatedAt }}</p>
              </div>
            </div>
          </CardContent>
        </Card>
        <Card class="h-full rounded-sm shadow-lg">
          <CardHeader>
            <CardTitle>Distance Sensor State</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="mx-auto space-x-2">
              <div class="grid grid-rows-2 gap-2 ">
                <p><span class="font-medium">Front: </span>{{ robotState.distanceSensor.frontDistance }} cm</p>
                <p><span class="font-medium">Back: </span>{{ robotState.distanceSensor.backDistance }} cm</p>
                <p><span class="font-medium">Down: </span>{{ robotState.distanceSensor.downDistance }} cm</p>
                <p><span class="font-medium">Last updated: </span>{{ robotState.distanceSensor.updatedAt }}</p>
              </div>
            </div>
          </CardContent>
        </Card>
        <Card class="h-full rounded-sm shadow-lg">
          <CardHeader>
            <CardTitle>Location State</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="mx-auto space-x-2">
              <div class="grid grid-rows-1 gap-2 ">
                <p><span class="font-medium">RFID tag: </span>{{ robotState.location.currentLocation }}</p>
                <p><span class="font-medium">Last updated:</span> {{ robotState.location.updatedAt }}</p>
              </div>
            </div>
          </CardContent>
        </Card>
      </CardContent>
    </div>
  </PageContainer>
</template>
