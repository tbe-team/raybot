<script setup lang="ts">
import type { LocationState } from '@/types/robot-state'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

interface Props {
  location: LocationState
}

const props = defineProps<Props>()
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle>Location</CardTitle>
    </CardHeader>
    <CardContent>
      <div class="grid grid-cols-2 gap-2 text-sm">
        <span class="font-medium text-muted-foreground">Current Location</span>
        <span>
          {{ props.location.currentLocation === '' ? 'N/A' : props.location.currentLocation }}
        </span>
      </div>
    </CardContent>
  </Card>
</template>
