<script setup lang="ts">
import type { ChargeState } from '@/types/state'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

interface Props {
  charge: ChargeState
}

const props = defineProps<Props>()
</script>

<template>
  <div>
    <Card class="h-full rounded-sm shadow-lg">
      <CardHeader>
        <CardTitle>Charge State</CardTitle>
      </CardHeader>

      <CardContent>
        <div class="mx-auto space-x-2">
          <div class="flex flex-col gap-2 text-sm">
            <p>
              <span class="font-medium">Current Limit: </span>{{ props.charge.currentLimit }} A
            </p>
            <p>
              <span class="font-medium">Charging: </span>
              <span :class="props.charge.enabled ? 'text-green-500' : 'text-red-500'">{{ props.charge.enabled ? 'Yes' : 'No' }}</span>
            </p>
            <p>
              <span class="font-medium">Last Updated: </span> {{ (props.charge.updatedAt).toISOString() }}
            </p>
          </div>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
