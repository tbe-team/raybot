<script setup lang="ts">
import CommandDetailSheet from '@/components/app/command-queue/CommandDetailSheet.vue'
import CreateCommandForm from '@/components/app/command-queue/CreateCommandForm.vue'
import CurrentProcessingCommandCard from '@/components/app/command-queue/CurrentProcessingCommandCard.vue'
import WaitingCommandList from '@/components/app/command-queue/WaitingCommandList.vue'
import PageContainer from '@/components/shared/PageContainer.vue'
import { Button } from '@/components/ui/button'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'
import { Loader2 } from 'lucide-vue-next'

const isDetailOpen = ref(false)

const selectedCommandId = ref<number | null>(null)

function handleCommandSelected(commandId: number) {
  selectedCommandId.value = commandId
  isDetailOpen.value = true
}
</script>

<template>
  <PageContainer>
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-xl font-semibold tracking-tight">
          Robot command queue
        </h1>
        <p class="text-sm text-muted-foreground">
          Manage and monitor robot commands in real-time
        </p>
      </div>
      <TooltipProvider>
        <Tooltip>
          <TooltipTrigger as-child>
            <Button variant="outline" size="icon">
              <Loader2 class="w-4 h-4 animate-spin" />
            </Button>
          </TooltipTrigger>
          <TooltipContent>
            <p>Refreshing in real-time</p>
          </TooltipContent>
        </Tooltip>
      </TooltipProvider>
    </div>

    <div class="grid grid-cols-1 gap-6 md:grid-cols-3">
      <div class="space-y-6 md:col-span-2">
        <CurrentProcessingCommandCard @view-details="handleCommandSelected" />
        <WaitingCommandList @view-details="handleCommandSelected" />
      </div>
      <div>
        <CreateCommandForm />
      </div>
    </div>

    <CommandDetailSheet
      v-if="selectedCommandId"
      v-model:is-open="isDetailOpen"
      :command-id="selectedCommandId"
    />
  </PageContainer>
</template>
