<script setup lang="ts">
import { ArrowUp } from 'lucide-vue-next'
import { onMounted, onUnmounted, ref } from 'vue'
import { Button } from '@/components/ui/button'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'

const visible = ref(false)

function handleScroll() {
  visible.value = window.scrollY > 200
}

function scrollToTop() {
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(() => {
  window.addEventListener('scroll', handleScroll)
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})
</script>

<template>
  <Tooltip>
    <TooltipTrigger as-child>
      <Button
        size="icon" class="fixed transition-opacity rounded-full shadow-lg z-1000 bottom-6 right-6"
        :class="visible ? 'opacity-100' : 'opacity-0 pointer-events-none'"
        @click="scrollToTop"
      >
        <ArrowUp class="w-5 h-5" />
      </Button>
    </TooltipTrigger>
    <TooltipContent>
      Scroll to top
    </TooltipContent>
  </Tooltip>
</template>
