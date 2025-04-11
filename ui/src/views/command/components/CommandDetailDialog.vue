<script setup lang="ts">
import type { CommandStatus } from '@/types/command'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { formatDate, formatDurationFromISO } from '@/lib/date'
import { VisuallyHidden } from 'reka-ui'
import { AlertCircle, CheckCircle, Clock, Loader, XCircle } from 'lucide-vue-next'
import { Badge } from '@/components/ui/badge'
import { useGetCommandQuery } from '@/composables/use-command'



const STATUS_LABELS: Record<CommandStatus, { label: string, class: string, icon: Component }> = {
QUEUED: { label: 'Queued', class: 'bg-yellow-500/10 text-yellow-500 hover:bg-yellow-500/20', icon: Clock },
PROCESSING: { label: 'Processing', class: 'bg-blue-500/10 text-blue-500 hover:bg-blue-500/20', icon: Loader },
SUCCEEDED: { label: 'Succeeded', class: 'bg-green-500/10 text-green-500 hover:bg-green-500/20', icon: CheckCircle },
FAILED: { label: 'Failed', class: 'bg-red-500/10 text-red-500 hover:bg-red-500/20', icon: AlertCircle },
CANCELED: { label: 'Canceled', class: 'bg-gray-500/10 text-gray-500 hover:bg-gray-500/20', icon: XCircle },
}

interface Props {
  commandId : number
}

const props = defineProps<Props>()
const commandId = computed(()=> props.commandId)
const open =defineModel<boolean>('open', {required: true})
const {data, refetch: refetchCommand} = useGetCommandQuery(commandId, {axiosOpts: {doNotShowLoading: true}, refetchInterval: 1000})

function getCommandTypeName(type: string) {
  return type.replace(/_/g, ' ').toLowerCase().replace(/\b\w/g, l => l.toUpperCase())
}

function getStatusBadge(status: CommandStatus){
  const { icon, label, class: className } = STATUS_LABELS[status]
  return h(Badge, {class: className}, [h(icon, {class: `w-4 h-4 mr-2 ${status === 'PROCESSING' ? 'animate-spin' : ''}`}), label])
}

function getSourceName(source: string) {
  return source === 'CLOUD' ? 'Cloud' : 'Application'
}

watch(open, (val)=> {
  if (val) {
    refetchCommand()
  }
})
</script>

<template>
  <Dialog v-model:open="open">
    <DialogContent class="sm:max-w-lg" >
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
      <div v-if="data" class="mt-4 space-y-4">
        <!-- ID -->
        <div class="flex">
          <div class="w-1/4 font-medium">
            ID:
          </div>
          <div>{{ data?.id }}</div>
        </div>

        <!-- Type -->
        <div class="flex">
          <div class="w-1/4 font-medium">
            Type:
          </div>
          <div>{{ getCommandTypeName(data?.type || '') }}</div>
        </div>

        <!-- Status -->
        <div class="flex items-center">
          <div class="w-1/4 font-medium">
            Status:
          </div>
          <div>
            <component :is="getStatusBadge(data?.status || 'QUEUED')"/>
          </div>
        </div>

        <!-- Source -->
        <div class="flex">
          <div class="w-1/4 font-medium">
            Source:
          </div>
          <div>{{ getSourceName(data?.source || '') }}</div>
        </div>

        <!-- Created -->
        <div class="flex">
          <div class="w-1/4 font-medium">
            Created:
          </div>
          <div>{{ formatDate(data?.createdAt || '') }}</div>
        </div>

        <!-- Completed -->
        <div class="flex" v-if="data?.completedAt">
          <div class="w-1/4 font-medium">
            Completed:
          </div>
          <div>{{ formatDate(data?.completedAt || '') }}</div>
        </div>

        <!-- Duration -->
        <div class="flex" v-if="data?.completedAt">
          <div class="w-1/4 font-medium">
            Duration:
          </div>
          <div>{{ formatDurationFromISO(data?.createdAt || '', data?.completedAt || '') }}</div>  
        </div>

        <!-- Error -->
        <div class="space-y-2" v-if="data?.error">
          <div class="font-medium">
            Error:
          </div>
          <div class="p-4 text-red-500 rounded-lg bg-muted-foreground/30 dark:bg-black">
            <pre class="text-sm whitespace-pre-wrap ">{{ data?.error }}</pre>
          </div>
        </div>

        <!-- Inputs -->
        <div v-if="data?.inputs" class="space-y-2">
          <div class="font-medium">
            Inputs:
          </div>
          <div class="p-4 rounded-lg bg-muted-foreground/30 dark:bg-black">
            <pre class="text-sm whitespace-pre-wrap ">{{ JSON.stringify(data?.inputs, null, 2) }}</pre>
          </div>
        </div>
      </div>
    </DialogContent>
  </Dialog>
</template>
