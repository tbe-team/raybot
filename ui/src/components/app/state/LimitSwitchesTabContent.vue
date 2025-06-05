<script setup lang="ts">
import { AlertCircle, Loader } from 'lucide-vue-next'
import LimitSwitchItem from '@/components/app/state/LimitSwitchItem.vue'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { useLimitSwitchStateQuery } from '@/composables/use-limit-switch-state'

const props = defineProps<{
  refreshInterval: number
}>()

const refetchInterval = computed(() => props.refreshInterval)

const { data: limitSwitchState, isPending, isError, error } = useLimitSwitchStateQuery({
  refetchInterval,
  axiosOpts: {
    doNotShowLoading: true,
  },
})
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle>Limit Switches</CardTitle>
    </CardHeader>
    <CardContent v-if="isPending">
      <div class="flex items-center justify-center gap-2">
        <Loader class="w-8 h-8 animate-spin text-muted-foreground" />
        <p class="text-muted-foreground">
          Loading limit switch state...
        </p>
      </div>
    </CardContent>
    <CardContent v-else-if="isError">
      <div class="flex items-center justify-center gap-2">
        <AlertCircle class="w-8 h-8 text-destructive" />
        <p class="text-destructive">
          {{ error?.message }}
        </p>
      </div>
    </CardContent>
    <CardContent v-else-if="limitSwitchState">
      <div class="grid grid-cols-2 gap-4 md:grid-cols-3 xl:grid-cols-5">
        <LimitSwitchItem :switch-data="limitSwitchState.limitSwitch1" name="Limit Switch 1" />
      </div>
    </CardContent>
  </Card>
</template>
