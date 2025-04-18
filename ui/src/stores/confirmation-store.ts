// confirmation-store.ts
import { defineStore } from 'pinia'

interface ConfirmationState {
  open: boolean
  title: string | null
  description: string | null
  cancelLabel: string | null
  actionLabel: string | null
  onAction: () => void
  onCancel: () => void
}

export const useConfirmationStore = defineStore('confirmation', {
  state: (): ConfirmationState => ({
    open: false,
    title: null,
    description: null,
    cancelLabel: null,
    actionLabel: null,
    onAction: () => {},
    onCancel: () => {},
  }),

  actions: {
    openConfirmation(data: {
      title: string
      description: string
      cancelLabel: string
      actionLabel: string
      onAction: () => void
      onCancel: () => void
    }) {
      this.open = true
      this.title = data.title
      this.description = data.description
      this.cancelLabel = data.cancelLabel
      this.actionLabel = data.actionLabel
      this.onAction = data.onAction
      this.onCancel = data.onCancel
    },

    closeConfirmation() {
      this.open = false
      this.title = null
      this.description = null
      this.cancelLabel = null
      this.actionLabel = null
      this.onAction = () => {}
      this.onCancel = () => {}
    },
  },
})
