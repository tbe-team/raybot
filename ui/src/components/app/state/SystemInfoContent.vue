<script setup lang="ts">
import { CircleAlert } from 'lucide-vue-next'
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
    <CardHeader>
      <CardTitle>System Information</CardTitle>
    </CardHeader>
    <CardContent>
      <!-- Loading state -->
      <div v-if="isPending" class="flex justify-center items-center py-6">
        <div class="w-6 h-6 rounded-full border-2 animate-spin border-primary border-t-transparent" />
      </div>
      <!-- Error state -->
      <div v-else-if="isError" class="py-6">
        <div class="text-center text-red-500">
          <CircleAlert class="mx-auto mb-2 w-8 h-8" />
          <p class="text-sm font-medium">
            An error occurred
          </p>
          <p class="mt-1 text-xs text-muted-foreground">
            {{ error?.message }}
          </p>
        </div>
      </div>
      <!-- Data state -->
      <div v-else-if="data" class="grid grid-cols-2 gap-2 text-sm">
        <span class="font-medium text-muted-foreground">Local IP</span>
        <span>
          {{ data.localIp }}
        </span>
        <span class="font-medium text-muted-foreground">Uptime</span>
        <span>
          {{ formatUptimeShort(data.uptime) }}
        </span>
        <span class="font-medium text-muted-foreground">CPU Usage</span>
        <span :class="[getUsageColor(data.cpuUsage)]">
          {{ data.cpuUsage.toFixed(1) }}%
        </span>
        <span class="font-medium text-muted-foreground">Memory Usage</span>
        <span :class="[getUsageColor(data.memoryUsage)]">
          {{ data.memoryUsage.toFixed(1) }}%
        </span>
        <span class="font-medium text-muted-foreground">Total Memory</span>
        <span>
          {{ formatMemory(data.totalMemory) }}

        </span>
      </div>
    </CardContent>
  </Card>
</template>
