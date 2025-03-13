<script setup lang="ts">
import PageContainer from '@/components/shared/PageContainer.vue'
import { locationState, robotState } from '@/views/state/components/fake-data'
import StatePage from '@/views/state/components/StatePage.vue'
import { Loader } from 'lucide-vue-next'

const isPending = ref(false)
const isError = ref(false)
const error = ref({ message: '' })
</script>

<template>
  <PageContainer>
    <div v-if="isPending" class="flex flex-col items-center justify-center gap-4 pt-20">
      <div class="flex items-center gap-4">
        <Loader class="w-8 h-8 animate-spin text-muted-foreground" />
      </div>
      <p class="text-lg text-muted-foreground">
        Loading state...
      </p>
    </div>

    <div v-else-if="isError" class="flex flex-col items-center justify-center gap-4 pt-20">
      <Card class="flex flex-col items-center gap-4 p-6 text-destructive">
        <AlertCircle class="w-8 h-8" />
        <div class="space-y-2 text-center">
          <h2 class="text-lg font-semibold">
            Failed to load state
          </h2>
          <p class="text-sm text-muted-foreground">
            {{ error?.message || 'An unexpected error occurred' }}
          </p>
        </div>
      </Card>
    </div>

    <div v-else-if="!robotState || !locationState" class="flex flex-col items-center justify-center gap-4 pt-20">
      <Card class="flex flex-col items-center gap-4 p-6">
        <AlertCircle class="w-8 h-8 text-muted-foreground" />
        <div class="space-y-2 text-center">
          <h2 class="text-lg font-semibold">
            No State Found
          </h2>
          <p class="text-sm text-muted-foreground">
            The state appears to be empty
          </p>
        </div>
      </Card>
    </div>

    <div v-else class="flex flex-col w-full gap-4">
      <StatePage :robot-state="robotState" :location-state="locationState" />
    </div>
  </PageContainer>
</template>
