<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { useConfirmationStore } from '@/stores/confirmation-store'
import { DialogOverlay } from 'reka-ui'

const confirmation = useConfirmationStore()

function handleAction() {
  confirmation.onAction()
  confirmation.closeConfirmation()
}
</script>

<template>
  <Dialog :open="confirmation.open" @update:open="confirmation.closeConfirmation">
    <DialogOverlay class="z-50" />
    <DialogContent>
      <DialogHeader>
        <DialogTitle>{{ confirmation.title }}</DialogTitle>
        <DialogDescription>
          {{ confirmation.description }}
        </DialogDescription>
      </DialogHeader>
      <DialogFooter>
        <Button
          variant="outline"
          @click="confirmation.closeConfirmation"
        >
          {{ confirmation.cancelLabel }}
        </Button>
        <Button
          @click="handleAction"
        >
          {{ confirmation.actionLabel }}
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
