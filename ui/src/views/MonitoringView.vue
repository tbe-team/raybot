<script setup lang="ts">
import { Loader } from 'lucide-vue-next'
import MonitoringConfigForm from '@/components/app/monitoring/MonitoringConfigForm.vue'
import PageContainer from '@/components/shared/PageContainer.vue'
import { useMonitoringConfigQuery } from '@/composables/use-config'

const { data, isPending } = useMonitoringConfigQuery()
</script>

<template>
  <PageContainer>
    <div>
      <h1 class="text-xl font-semibold tracking-tight">
        Monitoring Configuration
      </h1>
      <p class="text-sm text-muted-foreground">
        Configure monitoring alerts and thresholds for your robot
      </p>
    </div>

    <div v-if="isPending" class="flex justify-center items-center h-64">
      <Loader class="w-6 h-6 animate-spin" />
    </div>
    <MonitoringConfigForm v-else-if="data" :initial-values="data" />
  </PageContainer>
</template>
