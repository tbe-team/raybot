<script setup lang="ts">
import type { Command } from '@/types/command'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'
import { formatDate } from '@/lib/date'
import { VisuallyHidden } from 'reka-ui'
import StatusBadge from './commands-table/StatusBadge.vue'

const props = defineProps<{
  command: Command
}>()

function getCommandTypeName(type: string) {
  return type.replace(/_/g, ' ').toLowerCase().replace(/\b\w/g, l => l.toUpperCase())
}

function getSourceName(source: string) {
  return source === 'CLOUD' ? 'Cloud' : 'Application'
}
</script>

<template>
  <Dialog>
    <DialogTrigger as-child>
      <slot />
    </DialogTrigger>
    <DialogContent class="sm:max-w-md">
      <DialogHeader>
        <DialogTitle>
          Command Details
        </DialogTitle>
        <DialogDescription>
          <VisuallyHidden>
            View the details of the command
          </VisuallyHidden>
        </DialogDescription>
      </DialogHeader>
      <div class="mt-4 space-y-4">
        <!-- ID -->
        <div class="flex">
          <div class="w-1/4 font-medium">
            ID:
          </div>
          <div>{{ props.command.id }}</div>
        </div>

        <!-- Type -->
        <div class="flex">
          <div class="w-1/4 font-medium">
            Type:
          </div>
          <div>{{ getCommandTypeName(props.command.type) }}</div>
        </div>

        <!-- Status -->
        <div class="flex items-center">
          <div class="w-1/4 font-medium">
            Status:
          </div>
          <div>
            <StatusBadge :status="props.command.status" />
          </div>
        </div>

        <!-- Source -->
        <div class="flex">
          <div class="w-1/4 font-medium">
            Source:
          </div>
          <div>{{ getSourceName(props.command.source) }}</div>
        </div>

        <!-- Created -->
        <div class="flex">
          <div class="w-1/4 font-medium">
            Created:
          </div>
          <div>{{ formatDate(props.command.createdAt) }}</div>
        </div>

        <!-- Inputs -->
        <div v-if="props.command.inputs" class="space-y-2">
          <div class="font-medium">
            Inputs:
          </div>
          <div class="p-4 rounded-lg bg-muted-foreground/30 dark:bg-black">
            <pre class="text-sm whitespace-pre-wrap ">{{ JSON.stringify(props.command.inputs, null, 2) }}</pre>
          </div>
        </div>
      </div>
    </DialogContent>
  </Dialog>
</template>
