<script setup lang="ts">
import type { AppConnection } from '@/types/app-connection'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { formatUptime } from '@/lib/date'
import { Cloud } from 'lucide-vue-next'

const props = defineProps<{
  appConnection: AppConnection
}>()

const getPeripheralConnectedCount = computed(() => {
  let count = 0
  if (props.appConnection.espSerialConnection.connected)
    count++
  if (props.appConnection.picSerialConnection.connected)
    count++
  if (props.appConnection.rfidUsbConnection.connected)
    count++
  return count
})

const getTotalPeripherals = computed(() => 3)
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle>Connections</CardTitle>
      <CardDescription>
        {{ getPeripheralConnectedCount }} of {{ getTotalPeripherals }} peripherals connected
      </CardDescription>
    </CardHeader>
    <CardContent>
      <div class="space-y-4">
        <div class="flex items-center gap-4">
          <div
            class="p-2 rounded-full"
            :class="props.appConnection.cloudConnection.connected ? 'bg-success/10' : 'bg-red-500/10'"
          >
            <Cloud
              class="w-4 h-4"
              :class="props.appConnection.cloudConnection.connected ? 'text-success' : 'text-red-500'"
            />
          </div>
          <div class="flex-1">
            <p class="text-sm font-medium">
              Cloud connection
            </p>
            <p class="text-xs text-muted-foreground">
              Uptime: {{ formatUptime(props.appConnection.cloudConnection.uptime) }}
            </p>
          </div>
          <Badge :variant="props.appConnection.cloudConnection.connected ? 'default' : 'destructive'">
            {{ props.appConnection.cloudConnection.connected ? 'Online' : 'Offline' }}
          </Badge>
        </div>
      </div>
    </CardContent>
  </Card>
</template>
