<script setup lang="ts">
import type { LocationState } from '@/types/state'
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
          <div class="grid grid-rows-1 gap-2 text-sm">
            <div class="grid grid-cols-2 gap-2">
              <p><span class="font-medium">Current: </span>{{ props.location.currentLocation }}</p>
              <p><span class="font-medium">Target: </span>{{ props.location.targetLocation }}</p>
            </div>
            <p><span class="font-medium">Last Updated:</span> {{ (props.location.updatedAt).toISOString() }}</p>
          </div>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
