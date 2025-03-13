<script setup lang="ts">
import type { DischargeState } from '@/types/state'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

interface Props {
  discharge: DischargeState
}

const props = defineProps<Props>()
</script>

<template>
  <div>
    <Card class="h-full rounded-sm shadow-lg ">
      <CardHeader>
        <CardTitle>Discharge State</CardTitle>
      </CardHeader>

      <CardContent>
        <div class="mx-auto space-x-2">
          <div class="flex flex-col gap-2 text-sm">
            <p>
              <span class="font-medium">Current Limit: </span>{{ props.discharge.currentLimit }} A
            </p>
            <p>
              <span class="font-medium">Discharging: </span>
              <span :class="props.discharge.enabled ? 'text-green-500' : 'text-red-500'">{{ props.discharge.enabled ? 'Yes' : 'No' }}</span>
            </p>
            <p>
              <span class="font-medium">Last Updated: </span>{{ (props.discharge.updatedAt).toISOString() }}
            </p>
          </div>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
