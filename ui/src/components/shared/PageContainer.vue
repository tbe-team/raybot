<script setup lang="ts">
import { cn } from '@/lib/utils'

interface Props {
  /**
   * Whether to add vertical spacing between children
   * @default true
   */
  spacing?: boolean
  class?: string
}

const props = withDefaults(defineProps<Props>(), {
  spacing: true,
})
</script>

<template>
  <div class="container p-6">
    <div :class="cn(props.spacing && 'space-y-6', props.class)">
      <slot />
    </div>
  </div>
</template>
