<script setup lang="ts">
import type { ProgressRootProps } from 'reka-ui'
import type { HTMLAttributes } from 'vue'
import { cn } from '@/lib/utils'
import {
  ProgressIndicator,
  ProgressRoot,

} from 'reka-ui'
import { computed } from 'vue'

const props = withDefaults(
  defineProps<ProgressRootProps & { class?: HTMLAttributes['class'], variant?: string, max?: number }>(),
  {
    modelValue: 0,
    max: 100,
  },
)

const delegatedProps = computed(() => {
  const { class: _, ...delegated } = props

  return delegated
})
</script>

<template>
  <ProgressRoot
    v-bind="delegatedProps"
    :class="
      cn(
        'relative h-2 w-full overflow-hidden rounded-full bg-primary/20',
        props.class,
      )
    "
  >
    <ProgressIndicator
      class="flex-1 w-full h-full transition-all"
      :class="props.variant ? `bg-${props.variant}` : 'bg-primary'"
      :style="`transform: translateX(-${max - (props.modelValue ?? 0)}%);`"
    />
  </ProgressRoot>
</template>
