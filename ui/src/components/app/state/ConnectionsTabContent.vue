<script setup lang="ts">
import type { AppConnection } from '@/types/app-connection'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { formatDate, formatUptime } from '@/lib/date'
import { Cloud, Cpu, Usb } from 'lucide-vue-next'

const props = defineProps<{
  appConnection: AppConnection
}>()
</script>

<template>
  <div class="space-y-4">
    <Card>
      <CardHeader>
        <CardTitle>Connection Details</CardTitle>
        <CardDescription>Current connection status and details</CardDescription>
      </CardHeader>
      <CardContent>
        <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
          <div class="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle class="text-base">
                  Cloud Connection
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div class="flex items-center gap-4 mb-4">
                  <div
                    class="p-2 rounded-full"
                    :class="props.appConnection.cloudConnection.connected ? 'bg-green-100' : 'bg-destructive/10'"
                  >
                    <Cloud
                      class="w-5 h-5"
                      :class="props.appConnection.cloudConnection.connected ? 'text-green-600' : 'text-destructive'"
                    />
                  </div>
                  <div>
                    <p class="text-sm font-medium">
                      Status
                    </p>
                    <p class="text-lg font-bold">
                      {{ props.appConnection.cloudConnection.connected ? 'Connected' : 'Disconnected' }}
                    </p>
                  </div>
                </div>
                <div class="space-y-2">
                  <div class="flex justify-between text-sm">
                    <span class="text-muted-foreground">Uptime</span>
                    <span>{{ formatUptime(props.appConnection.cloudConnection.uptime) }}</span>
                  </div>
                  <div class="flex justify-between text-sm">
                    <span class="text-muted-foreground">Last Connected</span>
                    <span>{{ props.appConnection.cloudConnection.lastConnectedAt ? formatDate(props.appConnection.cloudConnection.lastConnectedAt) : 'Never' }}</span>
                  </div>
                </div>
                <div v-if="props.appConnection.cloudConnection.error" class="mt-2 text-xs text-destructive">
                  Error: {{ props.appConnection.cloudConnection.error }}
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle class="text-base">
                  ESP Serial Connection
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div class="flex items-center gap-4 mb-4">
                  <div
                    class="p-2 rounded-full"
                    :class="props.appConnection.espSerialConnection.connected ? 'bg-green-100' : 'bg-destructive/10'"
                  >
                    <Cpu
                      class="w-5 h-5"
                      :class="props.appConnection.espSerialConnection.connected ? 'text-green-600' : 'text-destructive'"
                    />
                  </div>
                  <div>
                    <p class="text-sm font-medium">
                      Status
                    </p>
                    <p class="text-lg font-bold">
                      {{ props.appConnection.espSerialConnection.connected ? 'Connected' : 'Disconnected' }}
                    </p>
                  </div>
                </div>
                <div class="space-y-2">
                  <div class="flex justify-between text-sm">
                    <span class="text-muted-foreground">Last Connected</span>
                    <span>{{ props.appConnection.espSerialConnection.lastConnectedAt ? formatDate(props.appConnection.espSerialConnection.lastConnectedAt) : 'Never' }}</span>
                  </div>
                </div>
                <div v-if="props.appConnection.espSerialConnection.error" class="mt-2 text-xs text-destructive">
                  Error: {{ props.appConnection.espSerialConnection.error }}
                </div>
              </CardContent>
            </Card>
          </div>

          <div class="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle class="text-base">
                  PIC Serial Connection
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div class="flex items-center gap-4 mb-4">
                  <div
                    class="p-2 rounded-full"
                    :class="props.appConnection.picSerialConnection.connected ? 'bg-green-100' : 'bg-destructive/10'"
                  >
                    <Cpu
                      class="w-5 h-5"
                      :class="props.appConnection.picSerialConnection.connected ? 'text-green-600' : 'text-destructive'"
                    />
                  </div>
                  <div>
                    <p class="text-sm font-medium">
                      Status
                    </p>
                    <p class="text-lg font-bold">
                      {{ props.appConnection.picSerialConnection.connected ? 'Connected' : 'Disconnected' }}
                    </p>
                  </div>
                </div>
                <div class="space-y-2">
                  <div class="flex justify-between text-sm">
                    <span class="text-muted-foreground">Last Connected</span>
                    <span>{{ props.appConnection.picSerialConnection.lastConnectedAt ? formatDate(props.appConnection.picSerialConnection.lastConnectedAt) : 'Never' }}</span>
                  </div>
                </div>
                <div v-if="props.appConnection.picSerialConnection.error" class="mt-2 text-xs text-destructive">
                  Error: {{ props.appConnection.picSerialConnection.error }}
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle class="text-base">
                  RFID USB Connection
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div class="flex items-center gap-4 mb-4">
                  <div
                    class="p-2 rounded-full"
                    :class="props.appConnection.rfidUsbConnection.connected ? 'bg-green-100' : 'bg-destructive/10'"
                  >
                    <Usb
                      class="w-5 h-5"
                      :class="props.appConnection.rfidUsbConnection.connected ? 'text-green-600' : 'text-destructive'"
                    />
                  </div>
                  <div>
                    <p class="text-sm font-medium">
                      Status
                    </p>
                    <p class="text-lg font-bold">
                      {{ props.appConnection.rfidUsbConnection.connected ? 'Connected' : 'Disconnected' }}
                    </p>
                  </div>
                </div>
                <div class="space-y-2">
                  <div class="flex justify-between text-sm">
                    <span class="text-muted-foreground">Last Connected</span>
                    <span>{{ props.appConnection.rfidUsbConnection.lastConnectedAt ? formatDate(props.appConnection.rfidUsbConnection.lastConnectedAt) : 'Never' }}</span>
                  </div>
                </div>
                <div v-if="props.appConnection.rfidUsbConnection.error" class="mt-2 text-xs text-destructive">
                  Error: {{ props.appConnection.rfidUsbConnection.error }}
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
