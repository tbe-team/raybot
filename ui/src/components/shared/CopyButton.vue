<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { Check, Copy } from 'lucide-vue-next'
import { ref } from 'vue'

const props = defineProps<{
  text: string
  className?: string
}>()

const emit = defineEmits<{
  (e: 'copy'): void
}>()

const hasCopied = ref(false)

async function copyToClipboard() {
  await navigator.clipboard.writeText(props.text)
  hasCopied.value = true
  emit('copy')
  setTimeout(() => {
    hasCopied.value = false
  }, 2000)
}
</script>

<template>
  <Button
    variant="ghost"
    size="sm"
    class="h-8 px-2" :class="[props.className]"
    @click="copyToClipboard"
  >
    <template v-if="hasCopied">
      <Check class="w-4 h-4 mr-1" />
      Copied
    </template>
    <template v-else>
      <Copy class="w-4 h-4 mr-1" />
      Copy
    </template>
  </Button>
</template>
