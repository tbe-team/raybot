<script setup lang="ts">
import type { Outdated } from './outdated'
import type { Cargo } from '@/types/cargo'
import { TriangleAlert } from 'lucide-vue-next'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { clearOutdatedTimeout, setOutdated } from './outdated'

interface Props {
  cargo: Cargo
}

const props = defineProps<Props>()
const expiredAfter = 10
const cargoOutdated = reactive<Outdated>({
  isOutdated: false,
  timeoutId: null,
})

watch(() => props.cargo.updatedAt, (newVal) => {
  setOutdated(newVal, expiredAfter, cargoOutdated)
}, { immediate: true })

onUnmounted(() => {
  clearOutdatedTimeout(cargoOutdated)
})
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle class="flex gap-3 items-center">
        Cargo
        <div v-if="cargoOutdated.isOutdated" class="flex gap-1 items-center text-warning">
          <TriangleAlert class="size-4" />
          <span class="text-xs font-normal">{{ `Last updated ${cargoOutdated.timeAgo}` || 'Outdated' }}</span>
        </div>
      </CardTitle>
    </CardHeader>
    <CardContent>
      <div class="grid grid-cols-2 gap-2 text-sm">
        <span class="font-medium text-muted-foreground">QR Code</span>
        <span>
          {{ props.cargo.qrCode }}
        </span>
        <span class="font-medium text-muted-foreground">Is Open</span>
        <span>
          <Badge class="text-white" :class="props.cargo.isOpen ? '!bg-success' : '!bg-destructive'">
            {{ props.cargo.isOpen ? 'Yes' : 'No' }}
          </Badge>
        </span>
        <span class="font-medium text-muted-foreground">Has Item</span>
        <span>
          <Badge class="text-white" :class="props.cargo.hasItem ? '!bg-success' : '!bg-destructive'">
            {{ props.cargo.hasItem ? 'Yes' : 'No' }}
          </Badge>
        </span>
        <span class="font-medium text-muted-foreground">Bottom Distance</span>
        <span>
          {{ props.cargo.bottomDistance }} cm
        </span>
      </div>
    </CardContent>
  </Card>
</template>
