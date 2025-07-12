<script setup lang="ts">
import type { Leds } from '@/types/robot-state'
import { CircleCheck, XCircle } from 'lucide-vue-next'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Separator } from '@/components/ui/separator'
import { formatDate } from '@/lib/date'

interface Props {
  led: Leds
}

const props = defineProps<Props>()
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle>LED</CardTitle>
    </CardHeader>
    <CardContent>
      <div class="grid grid-cols-2 gap-2 text-sm">
        <span class="font-medium text-muted-foreground">Alert LED Connection</span>
        <span>
          <Badge class="text-white" :class="props.led.alertLed.connection.connected ? 'bg-success hover:!bg-success' : 'bg-destructive hover:bg-destructive'">
            <CircleCheck v-if="props.led.alertLed.connection.connected" class="mr-2 size-3" />
            <XCircle v-else class="mr-2 size-3" />
            {{ props.led.alertLed.connection.connected ? 'Connected' : 'Disconnected' }}
          </Badge>
        </span>
        <!-- Not Connected -->
        <template v-if="!props.led.alertLed.connection.connected">
          <span class="font-medium text-muted-foreground">Last Connected</span>
          <span>
            {{ props.led.alertLed.connection.lastConnectedAt ? formatDate(props.led.alertLed.connection.lastConnectedAt) : 'Never' }}
          </span>
          <span class="font-medium text-muted-foreground">
            Error:
          </span>
          <span class="text-destructive">
            {{ props.led.alertLed.connection.error }}
          </span>
        </template>

        <!-- Connected -->
        <template v-else>
          <span class="font-medium text-muted-foreground">Mode</span>
          <span>
            <Badge :variant="props.led.alertLed.state.mode === 'OFF' ? 'secondary' : 'default'">
              {{ props.led.alertLed.state.mode }}
            </Badge>
          </span>
        </template>
        <Separator class="col-span-2 my-1" />
        <span class="font-medium text-muted-foreground">System LED Connection</span>
        <span>
          <Badge class="text-white" :class="props.led.systemLed.connection.connected ? 'bg-success hover:!bg-success' : 'bg-destructive hover:bg-destructive'">
            <CircleCheck v-if="props.led.systemLed.connection.connected" class="mr-2 size-3" />
            <XCircle v-else class="mr-2 size-3" />
            {{ props.led.systemLed.connection.connected ? 'Connected' : 'Disconnected' }}
          </Badge>
        </span>
        <!-- Not Connected -->
        <template v-if="!props.led.systemLed.connection.connected">
          <span class="font-medium text-muted-foreground">Last Connected</span>
          <span>
            {{ props.led.systemLed.connection.lastConnectedAt ? formatDate(props.led.systemLed.connection.lastConnectedAt) : 'Never' }}
          </span>
          <span class="font-medium text-muted-foreground">
            Error:
          </span>
          <span class="text-destructive">
            {{ props.led.systemLed.connection.error }}
          </span>
        </template>
        <!-- Connected -->
        <template v-else>
          <span class="font-medium text-muted-foreground">Mode</span>
          <span>
            <Badge :variant="props.led.systemLed.state.mode === 'OFF' ? 'secondary' : 'default'">
              {{ props.led.systemLed.state.mode }}
            </Badge>
          </span>
        </template>
      </div>
    </CardContent>
  </Card>
</template>
