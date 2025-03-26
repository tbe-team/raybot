<script setup lang="ts">
import type { PaginationFirstProps } from 'reka-ui'
import type { HTMLAttributes } from 'vue'
import {
  Button,
} from '@/components/ui/button'
import { cn } from '@/lib/utils'
import { ChevronsLeft } from 'lucide-vue-next'
import { PaginationFirst } from 'reka-ui'
import { computed } from 'vue'

const props = withDefaults(defineProps<PaginationFirstProps & { class?: HTMLAttributes['class'] }>(), {
  asChild: true,
})

const delegatedProps = computed(() => {
  const { class: _, ...delegated } = props

  return delegated
})
</script>

<template>
  <PaginationFirst v-bind="delegatedProps">
    <Button :class="cn('w-9 h-9 p-0', props.class)" variant="outline">
      <slot>
        <ChevronsLeft />
      </slot>
    </Button>
  </PaginationFirst>
</template>
