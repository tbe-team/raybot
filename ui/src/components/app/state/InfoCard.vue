<script setup lang="ts">
import type { Info } from '@/types/info'
import { Clock, Cpu, MemoryStick, Monitor, Wifi } from 'lucide-vue-next'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Progress } from '@/components/ui/progress'
// Props definition
defineProps<{
  info: Info
}>()

function formatUptime(seconds: number): string {
  const hours = Math.floor(seconds / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  const remainingSeconds = seconds % 60

  return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${remainingSeconds.toString().padStart(2, '0')}`
}

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
  <Card class="w-full max-w-md">
    <CardHeader class="pb-3">
      <CardTitle class="flex items-center gap-2 text-lg">
        <Monitor class="w-5 h-5" />
        System Information
      </CardTitle>
    </CardHeader>
    <CardContent class="space-y-4">
      <!-- IP Address -->
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-2">
          <Wifi class="w-4 h-4 text-muted-foreground" />
          <span class="text-sm font-medium">Local IP</span>
        </div>
        <Badge variant="outline" class="!font-mono">
          {{ info.localIp }}
        </Badge>
      </div>

      <!-- CPU Usage -->
      <div class="space-y-2">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2">
            <Cpu class="w-4 h-4 text-muted-foreground" />
            <span class="text-sm font-medium">CPU Usage</span>
          </div>
          <span class="text-sm font-semibold" :class="[getUsageColor(info.cpuUsage)]">
            {{ info.cpuUsage.toFixed(1) }}%
          </span>
        </div>
        <Progress :model-value="info.cpuUsage" class="h-2" />
      </div>

      <!-- Memory Usage -->
      <div class="space-y-2">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2">
            <MemoryStick class="w-4 h-4 text-muted-foreground" />
            <span class="text-sm font-medium">Memory Usage</span>
          </div>
          <span class="text-sm font-semibold" :class="[getUsageColor(info.memoryUsage)]">
            {{ info.memoryUsage.toFixed(1) }}%
          </span>
        </div>
        <Progress :model-value="info.memoryUsage" class="h-2" />
        <div class="text-xs text-right text-muted-foreground">
          Total: {{ formatMemory(info.totalMemory) }}
        </div>
      </div>

      <!-- Uptime -->
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-2">
          <Clock class="w-4 h-4 text-muted-foreground" />
          <span class="text-sm font-medium">Uptime</span>
        </div>
        <span class="font-mono text-sm font-normal text-muted-foreground">{{ formatUptime(info.uptime) }}</span>
      </div>
    </CardContent>
  </Card>
</template>
