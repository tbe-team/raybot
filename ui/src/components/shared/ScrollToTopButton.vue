<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { ArrowUp } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Tooltip, TooltipTrigger, TooltipContent } from '@/components/ui/tooltip'

const visible = ref(false)

const handleScroll = () => {
  visible.value = window.scrollY > 200
}

const scrollToTop = () => {
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
        size="icon"
        @click="scrollToTop"
        class="fixed transition-opacity rounded-full shadow-lg z-1000 bottom-6 right-6"
        :class="visible ? 'opacity-100' : 'opacity-0 pointer-events-none'"
      >
        <ArrowUp class="w-5 h-5" />
      </Button>
    </TooltipTrigger>
    <TooltipContent>
      Scroll to top
    </TooltipContent>
  </Tooltip>
</template>
