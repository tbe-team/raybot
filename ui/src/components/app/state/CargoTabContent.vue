<script setup lang="ts">
import type { Cargo, CargoDoorMotorState } from '@/types/cargo'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { formatDate } from '@/lib/date'
import { Lock, Package, Unlock } from 'lucide-vue-next'

const props = defineProps<{
  cargo: Cargo
  cargoDoorMotor: CargoDoorMotorState
}>()
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle>Cargo Status</CardTitle>
      <CardDescription>Current cargo information</CardDescription>
    </CardHeader>
    <CardContent>
      <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
        <div class="space-y-4">
          <div class="flex items-center gap-4">
            <div
              class="p-3 rounded-full"
              :class="props.cargo.isOpen ? 'bg-green-100' : 'bg-muted'"
            >
              <Package
                class="w-6 h-6"
                :class="props.cargo.isOpen ? 'text-green-600' : 'text-muted-foreground'"
              />
            </div>
            <div>
              <p class="text-sm text-muted-foreground">
                Cargo Status
              </p>
              <p class="text-xl font-bold">
                {{ props.cargo.isOpen ? 'Open' : 'Closed' }}
              </p>
            </div>
          </div>

          <div class="space-y-2">
            <h3 class="text-sm font-medium">
              Cargo Details
            </h3>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <p class="text-sm text-muted-foreground">
                  QR Code
                </p>
                <p class="text-lg font-medium">
                  {{ props.cargo.qrCode || 'No QR Code' }}
                </p>
              </div>
              <div>
                <p class="text-sm text-muted-foreground">
                  Bottom Distance
                </p>
                <p class="text-lg font-medium">
                  {{ props.cargo.bottomDistance }} cm
                </p>
              </div>
            </div>
          </div>

          <div class="text-xs text-muted-foreground">
            Last updated: {{ formatDate(props.cargo.updatedAt) }}
          </div>
        </div>

        <div class="space-y-4">
          <div class="flex items-center gap-4">
            <div
              class="p-3 rounded-full"
              :class="props.cargoDoorMotor.direction === 'OPEN' ? 'bg-green-100' : 'bg-yellow-100'"
            >
              <Unlock
                v-if="props.cargoDoorMotor.direction === 'OPEN'"
                class="w-6 h-6 text-green-600"
              />
              <Lock
                v-else
                class="w-6 h-6 text-yellow-600"
              />
            </div>
            <div>
              <p class="text-sm text-muted-foreground">
                Cargo Door
              </p>
              <p class="text-xl font-bold">
                {{ props.cargoDoorMotor.direction === 'OPEN' ? 'Open' : 'Closed' }}
              </p>
            </div>
          </div>

          <div class="space-y-2">
            <h3 class="text-sm font-medium">
              Motor Status
            </h3>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <p class="text-sm text-muted-foreground">
                  Speed
                </p>
                <p class="text-lg font-medium">
                  {{ props.cargoDoorMotor.speed }}%
                </p>
              </div>
              <div>
                <p class="text-sm text-muted-foreground">
                  Status
                </p>
                <Badge :variant="props.cargoDoorMotor.isRunning ? 'default' : 'outline'">
                  {{ props.cargoDoorMotor.isRunning ? 'Running' : 'Idle' }}
                </Badge>
              </div>
            </div>
          </div>

          <div class="text-xs text-muted-foreground">
            Last updated: {{ formatDate(props.cargoDoorMotor.updatedAt) }}
          </div>
        </div>
      </div>
    </CardContent>
  </Card>
</template>
