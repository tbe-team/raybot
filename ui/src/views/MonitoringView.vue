<script setup lang="ts">
import { AlertCircle, Loader } from 'lucide-vue-next'
import MonitoringConfigForm from '@/components/app/monitoring/MonitoringConfigForm.vue'
import PageContainer from '@/components/shared/PageContainer.vue'
import { useBatteryMonitoringConfigQuery } from '@/composables/use-config'

const { data, isPending, isError, error } = useBatteryMonitoringConfigQuery()
</script>

<template>
  <PageContainer>
    <div>
      <h1 class="text-xl font-semibold tracking-tight">
        Monitoring Configuration
      </h1>
      <p class="text-sm text-muted-foreground">
        Configure monitoring alerts and thresholds for your robot
        <RouterLink to="/reboot" class="text-blue-500">
          (Reboot to apply changes)
        </RouterLink>
      </p>
    </div>

    <div v-if="isPending" class="flex justify-center items-center h-64">
      <Loader class="w-6 h-6 animate-spin" />
    </div>

    <div v-else-if="isError" class="flex flex-col gap-4 justify-center items-center pt-20">
      <div class="flex flex-col gap-4 items-center p-6 text-red-500">
        <AlertCircle class="w-8 h-8" />
        <div class="space-y-2 text-center">
          <h2 class="text-lg font-semibold">
            Failed to load monitoring configuration
          </h2>
          <p class="text-sm text-muted-foreground">
            {{ error?.message || 'An unexpected error occurred' }}
          </p>
        </div>
      </div>
    </div>

    <MonitoringConfigForm v-else-if="data" :initial-values="data" />
  </PageContainer>
</template>
