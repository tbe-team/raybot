<script setup lang="ts">
import CopyButton from '@/components/shared/CopyButton.vue'
import { HoverCard, HoverCardContent, HoverCardTrigger } from '@/components/ui/hover-card'
import { useVersionQuery } from '@/composables/use-version'
import { Loader2 } from 'lucide-vue-next'

const { data: version, isPending } = useVersionQuery()
const isOpen = ref(false)

const versionText = computed(() => {
  if (!version.value)
    return ''

  return `Version: ${version.value.version}
          Build Date: ${version.value.buildDate}
          Go Version: ${version.value.goVersion}`
})
</script>

<template>
  <HoverCard :open="isOpen" @update:open="isOpen = $event">
    <HoverCardTrigger as-child>
      <span class="flex items-center gap-1 text-sm text-muted-foreground">
        <button
          class="px-1 -mx-1 transition-opacity duration-200 rounded hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
          :class="{ 'opacity-0': isPending, 'opacity-100': !isPending }"
          @click="isOpen = !isOpen"
        >
          {{ version?.version.substring(0, 7) }}
        </button>
        <template v-if="isPending">
          <Loader2 class="w-3 h-3 animate-spin" />
          <span>Loading...</span>
        </template>
      </span>
    </HoverCardTrigger>
    <HoverCardContent side="bottom" align="start" class="w-[400px]">
      <div class="p-2">
        <div class="flex items-center justify-between mb-2">
          <div class="font-semibold">
            Build Information
          </div>
          <CopyButton v-if="version && !isPending" :text="versionText" />
        </div>

        <div
          class="transition-opacity duration-200"
          :class="{ 'opacity-0': isPending, 'opacity-100': !isPending }"
        >
          <div v-if="version" class="mb-2 last:mb-0">
            <div class="font-semibold">
              {{ version.version }}
            </div>
            <div class="text-sm">
              <div>Build Date: {{ version.buildDate }}</div>
              <div>Go Version: {{ version.goVersion }}</div>
            </div>
          </div>
          <div v-else class="text-muted-foreground">
            No build information available
          </div>
        </div>

        <div v-if="isPending" class="flex items-center gap-2 mt-2">
          <Loader2 class="w-4 h-4 animate-spin" />
          <span>Loading build information...</span>
        </div>
      </div>
    </HoverCardContent>
  </HoverCard>
</template>
