<script setup lang="ts">
import { CircleAlert, Server } from 'lucide-vue-next'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { useSystemGetInfoQuery } from '@/composables/use-system'
import { formatUptimeShort } from '@/lib/date'
const REFRESH_INTERVAL_INFORMATION = 10000

const { data, isPending, isError, error } = useSystemGetInfoQuery({
  axiosOpts: { doNotShowLoading: true },
  refetchInterval: REFRESH_INTERVAL_INFORMATION,
})

function formatMemory(bytes: number): string {
  const gb = bytes / (1024)
  return `${gb.toFixed(1)} GB`
}

function getUsageColor(usage: number): string {
  if (usage < 50)
    return 'text-green-600'
  if (usage < 80)
    return 'text-yellow-600'
  return 'text-red-600'
}
</script>

<template>
  <Card>
    <CardHeader class="pb-2">
      <CardTitle class="flex items-center gap-2 text-sm font-medium">
        <Server class="w-5 h-5" />
        System Information
      </CardTitle>
    </CardHeader>
    <!-- Loading state -->
    <CardContent v-if="isPending" class="flex items-center justify-center py-6">
      <div class="w-6 h-6 border-2 rounded-full animate-spin border-primary border-t-transparent" />
    </CardContent>
    <!-- Error state -->
    <CardContent v-else-if="isError" class="py-6">
      <div class="text-center text-red-500">
        <CircleAlert class="w-8 h-8 mx-auto mb-2" />
        <p class="text-sm font-medium">
          An error occurred
        </p>
        <p class="mt-1 text-xs text-muted-foreground">
          {{ error?.message }}
        </p>
      </div>
    </CardContent>
    <!-- Data state -->
    <CardContent v-else-if="data" class="space-y-1">
      <!-- IP Address -->
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-2">
          <span class="text-sm font-medium">Local IP</span>
        </div>
        <span class="font-mono text-sm font-normal text-muted-foreground">{{ data.localIp }}</span>
      </div>

      <!-- Uptime -->
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-2">
          <span class="text-sm font-medium">Uptime</span>
        </div>
        <span class="font-mono text-sm font-normal text-muted-foreground">{{ formatUptimeShort(data.uptime) }}</span>
      </div>

      <!-- CPU Usage -->
      <div class="space-y-2">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2">
            <span class="text-sm font-medium">CPU Usage</span>
          </div>
          <span class="text-sm font-semibold" :class="[getUsageColor(data.cpuUsage)]">
            {{ data.cpuUsage.toFixed(1) }}%
          </span>
        </div>
      </div>

      <!-- Memory Usage -->
      <div class="space-y-2">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2">
            <span class="text-sm font-medium">Memory Usage</span>
          </div>
          <span class="text-sm font-semibold" :class="[getUsageColor(data.memoryUsage)]">
            {{ data.memoryUsage.toFixed(1) }}%
          </span>
        </div>
      </div>

      <!-- Total Memory -->
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-2">
          <span class="text-sm font-medium">Total Memory</span>
        </div>
        <span class="font-mono text-sm font-normal text-muted-foreground">{{ formatMemory(data.totalMemory) }}</span>
      </div>
    </CardContent>
  </Card>
</template>
