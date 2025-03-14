<script setup lang="ts">
import type { LocationState } from '@/types/robot-state'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

interface Props {
  location: LocationState
}

const props = defineProps<Props>()
</script>

<script lang="ts">
</script>

<template>
  <div>
    <Card class="h-full rounded-sm shadow-lg">
      <CardHeader>
        <CardTitle>Location State</CardTitle>
      </CardHeader>

      <CardContent>
        <div class="mx-auto space-x-2">
          <div class="grid grid-rows-1 gap-2 ">
            <p><span class="font-medium">RFID tag: </span>{{ props.location.currentLocation }}</p>
            <p><span class="font-medium">Last updated:</span> {{ props.location.updatedAt }}</p>
          </div>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
