<script setup lang="ts">
import type { AppConnection } from '@/types/app-connection'
import { CircleCheck, XCircle } from 'lucide-vue-next'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Separator } from '@/components/ui/separator'
import { formatDate, formatUptimeShort } from '@/lib/date'

interface Props {
  appConnection: AppConnection
}

const props = defineProps<Props>()
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle>Connections</CardTitle>
    </CardHeader>
    <CardContent>
      <div class="grid grid-cols-2 gap-2 text-sm">
        <span class="font-medium text-muted-foreground">ESP Serial Connection</span>
        <span>
          <Badge class="text-white" :class="props.appConnection.espSerialConnection.connected ? '!bg-success' : '!bg-destructive'">
            <CircleCheck v-if="props.appConnection.espSerialConnection.connected" class="mr-2 size-3" />
            <XCircle v-else class="mr-2 size-3" />
            {{ props.appConnection.espSerialConnection.connected ? 'Connected' : 'Disconnected' }}
          </Badge>
        </span>
        <span class="font-medium text-muted-foreground">Last Connected</span>
        <span>{{ props.appConnection.espSerialConnection.lastConnectedAt ? formatDate(props.appConnection.espSerialConnection.lastConnectedAt) : 'Never' }}</span>
        <template v-if="props.appConnection.espSerialConnection.error">
          <span class="font-medium text-muted-foreground">
            Error:
          </span>
          <span class="text-destructive">
            {{ props.appConnection.espSerialConnection.error }}
          </span>
        </template>
        <Separator class="col-span-2 my-1" />
        <span class="font-medium text-muted-foreground">PIC Serial Connection</span>
        <span>
          <Badge class="text-white" :class="props.appConnection.picSerialConnection.connected ? '!bg-success' : '!bg-destructive'">
            <CircleCheck v-if="props.appConnection.picSerialConnection.connected" class="mr-2 size-3" />
            <XCircle v-else class="mr-2 size-3" />
            {{ props.appConnection.picSerialConnection.connected ? 'Connected' : 'Disconnected' }}
          </Badge>
        </span>
        <span class="font-medium text-muted-foreground">Last Connected</span>
        <span>
          {{ props.appConnection.picSerialConnection.lastConnectedAt ? formatDate(props.appConnection.picSerialConnection.lastConnectedAt) : 'Never' }}
        </span>
        <template v-if="props.appConnection.picSerialConnection.error">
          <span class="font-medium text-muted-foreground">
            Error:
          </span>
          <span class="text-destructive">
            {{ props.appConnection.picSerialConnection.error }}
          </span>
        </template>
        <Separator class="col-span-2 my-1" />
        <span class="font-medium text-muted-foreground">RFID USB Connection</span>
        <span>
          <Badge class="text-white" :class="props.appConnection.rfidUsbConnection.connected ? 'bg-success hover:!bg-success' : 'bg-destructive hover:bg-destructive'">
            <CircleCheck v-if="props.appConnection.rfidUsbConnection.connected" class="mr-2 size-3" />
            <XCircle v-else class="mr-2 size-3" />
            {{ props.appConnection.rfidUsbConnection.connected ? 'Connected' : 'Disconnected' }}
          </Badge>
        </span>
        <span class="font-medium text-muted-foreground">Last Connected</span>
        <span>
          {{ props.appConnection.rfidUsbConnection.lastConnectedAt ? formatDate(props.appConnection.rfidUsbConnection.lastConnectedAt) : 'Never' }}
        </span>
        <template v-if="props.appConnection.rfidUsbConnection.error">
          <span class="font-medium text-muted-foreground">
            Error:
          </span>
          <span class="text-destructive">
            {{ props.appConnection.rfidUsbConnection.error }}
          </span>
        </template>
        <Separator class="col-span-2 my-1" />
        <span class="font-medium text-muted-foreground">Cloud Connection</span>
        <span>
          <Badge class="text-white" :class="props.appConnection.cloudConnection.connected ? 'bg-success hover:!bg-success' : 'bg-destructive hover:bg-destructive'">
            <CircleCheck v-if="props.appConnection.cloudConnection.connected" class="mr-2 size-3" />
            <XCircle v-else class="mr-2 size-3" />
            {{ props.appConnection.cloudConnection.connected ? 'Connected' : 'Disconnected' }}
          </Badge>
        </span>
        <span class="font-medium text-muted-foreground">Uptime</span>
        <span>
          {{ formatUptimeShort(props.appConnection.cloudConnection.uptime) }}
        </span>
        <span class="font-medium text-muted-foreground">Last Connected</span>
        <span>
          {{ props.appConnection.cloudConnection.lastConnectedAt ? formatDate(props.appConnection.cloudConnection.lastConnectedAt) : 'Never' }}
        </span>
        <template v-if="props.appConnection.cloudConnection.error">
          <span class="font-medium text-muted-foreground">
            Error:
          </span>
          <span class="text-destructive">
            {{ props.appConnection.cloudConnection.error }}
          </span>
        </template>
      </div>
    </CardContent>
  </Card>
</template>
