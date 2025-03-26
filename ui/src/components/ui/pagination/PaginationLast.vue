<script setup lang="ts">
import type { PaginationLastProps } from 'reka-ui'
import type { HTMLAttributes } from 'vue'
import {
  Button,
} from '@/components/ui/button'
import { cn } from '@/lib/utils'
import { ChevronsRight } from 'lucide-vue-next'
import { PaginationLast } from 'reka-ui'
import { computed } from 'vue'

const props = withDefaults(defineProps<PaginationLastProps & { class?: HTMLAttributes['class'] }>(), {
  asChild: true,
})

const delegatedProps = computed(() => {
  const { class: _, ...delegated } = props

  return delegated
})
</script>

<template>
  <PaginationLast v-bind="delegatedProps">
    <Button :class="cn('w-9 h-9 p-0', props.class)" variant="outline">
      <slot>
        <ChevronsRight />
      </slot>
    </Button>
  </PaginationLast>
</template>
