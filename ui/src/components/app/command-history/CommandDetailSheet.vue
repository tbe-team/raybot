<script setup lang="ts">
import type { Command } from '@/types/command'
import { Sheet, SheetContent, SheetDescription, SheetHeader, SheetTitle } from '@/components/ui/sheet'
import { formatDate, formatUptime } from '@/lib/date'
import SourceBadge from '../command-queue/SourceBadge.vue'
import StatusBadge from '../command-queue/StatusBadge.vue'
import { getCommandIcon, getCommandName } from '../command-queue/utils'

const props = defineProps<{
  command: Command
}>()

const isOpen = defineModel<boolean>('isOpen', { required: true })
</script>

<template>
  <Sheet v-model:open="isOpen">
    <SheetContent class="max-h-screen overflow-y-auto sm:max-w-xl">
      <SheetHeader>
        <SheetTitle>
          Command detail
        </SheetTitle>
        <SheetDescription>
          Real-time information about the selected command
        </SheetDescription>
      </SheetHeader>

      <div class="mt-6 space-y-6">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2 ">
            <component :is="getCommandIcon(props.command.type)" />
            <span class="text-sm font-medium">
              {{ getCommandName(props.command.type) }}
            </span>
          </div>
          <StatusBadge :status="props.command.status" />
        </div>

        <div class="grid grid-cols-2 gap-4">
          <div>
            <p class="text-sm text-muted-foreground">
              ID
            </p>
            <p class="font-medium">
              {{ props.command.id }}
            </p>
          </div>
          <div class="flex flex-col items-start gap-2">
            <p class="text-sm text-muted-foreground">
              Source
            </p>
            <SourceBadge :source="props.command.source" />
          </div>
          <div v-if="props.command.completedAt && props.command.startedAt">
            <p class="text-sm text-muted-foreground">
              Duration
            </p>
            <p class="font-medium">
              {{ formatUptime((new Date(props.command.completedAt).getTime() - new Date(props.command.startedAt).getTime()) / 1000) }}
            </p>
          </div>
        </div>

        <div class="space-y-2">
          <p class="text-sm font-medium">
            Command Inputs
          </p>
          <div class="p-4 bg-gray-100 rounded-xl dark:bg-gray-800">
            <span class="font-mono text-sm text-gray-800 break-words whitespace-pre-wrap dark:text-gray-200">
              {{ JSON.stringify(props.command.inputs, null, 4) }}
            </span>
          </div>
        </div>
        <div v-if="props.command.status !== 'PROCESSING' && props.command.status !== 'QUEUED'" class="space-y-2">
          <p class="text-sm font-medium">
            Command Outputs
          </p>
          <div class="p-4 bg-gray-100 rounded-xl dark:bg-gray-800">
            <span class="font-mono text-sm text-gray-800 break-words whitespace-pre-wrap dark:text-gray-200">
              {{ JSON.stringify(props.command.outputs, null, 4) }}
            </span>
          </div>
        </div>

        <div v-if="props.command.error" class="space-y-2">
          <p class="text-sm font-medium text-red-500">
            Error
          </p>
          <div class="p-3 text-red-500 rounded-md bg-red-500/10">
            {{ props.command.error }}
          </div>
        </div>

        <div class="space-y-2">
          <p class="text-sm font-medium">
            Timeline
          </p>
          <div class="space-y-3">
            <div class="flex justify-between text-sm">
              <span class="text-muted-foreground">Created</span>
              <span>{{ formatDate(props.command.createdAt) }}</span>
            </div>
            <div v-if="props.command.startedAt" class="flex justify-between text-sm">
              <span class="text-muted-foreground">Started</span>
              <span>{{ formatDate(props.command.startedAt) }}</span>
            </div>
            <div v-if="props.command.completedAt" class="flex justify-between text-sm">
              <span class="text-muted-foreground">Completed</span>
              <span>{{ formatDate(props.command.completedAt) }}</span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-muted-foreground">Last Updated</span>
              <span>{{ formatDate(props.command.updatedAt) }}</span>
            </div>
          </div>
        </div>
      </div>
    </SheetContent>
  </Sheet>
</template>
