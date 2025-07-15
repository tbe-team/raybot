<script setup lang="ts">
import { AlertTriangle, CheckCircle } from 'lucide-vue-next'
import { useSystemStatusQuery } from '@/composables/use-system'

const { data: systemStatus } = useSystemStatusQuery({
  refetchInterval: 3000,
  axiosOpts: {
    doNotShowLoading: true,
  },
})

const showStatus = computed(() => {
  return systemStatus.value?.status
})
</script>

<template>
  <div class="absolute top-0 left-10 z-50 m-2 ml-2 sm:ml-4 lg:ml-4">
    <!-- Error Status -->
    <div
      v-if="showStatus === 'ERROR'"
      class="flex items-center gap-2 px-2 py-1.5 sm:px-3 sm:py-2 bg-destructive text-destructive-foreground rounded-md shadow-lg mb-2"
    >
      <AlertTriangle class="w-4 h-4 flex-shrink-0" />
      <span class="text-xs font-medium whitespace-nowrap">System Error</span>
    </div>

    <!-- Normal Status -->
    <div
      v-if="showStatus === 'NORMAL'"
      class="flex items-center gap-2 px-2 py-1.5 sm:px-3 sm:py-2 bg-green-500 text-white rounded-md shadow-lg"
    >
      <CheckCircle class="w-4 h-4 flex-shrink-0" />
      <span class="text-xs font-medium whitespace-nowrap">System Normal</span>
    </div>
  </div>
</template>
