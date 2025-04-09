<script setup lang="ts">
import type { LocationState } from '@/types/robot-state'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { formatDate } from '@/lib/date'
import { MapPin } from 'lucide-vue-next'

const props = defineProps<{
  location: LocationState
}>()
</script>

<template>
  <Card>
    <CardHeader class="pb-2">
      <CardTitle class="flex items-center gap-2 text-sm font-medium">
        <MapPin class="w-6 h-6" />
        Location
      </CardTitle>
    </CardHeader>
    <CardContent>
      <div class="text-lg font-medium">
        {{ props.location.currentLocation === '' ? 'N/A' : props.location.currentLocation }}
      </div>
      <p class="mt-2 text-xs text-muted-foreground">
        Updated: {{ formatDate(props.location.updatedAt) }}
      </p>
    </CardContent>
  </Card>
</template>
