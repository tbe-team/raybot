<script setup lang="ts">
import type { Props as LimitSwitchItemProps } from './LimitSwitchItem.vue'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Separator } from '@/components/ui/separator'
import { useLimitSwitchStateQuery } from '@/composables/use-limit-switch-state'
import LimitSwitchItem from './LimitSwitchItem.vue'

const refetchInterval = defineModel<number>('refetchInterval', { required: true })

const { data: limitSwitches } = useLimitSwitchStateQuery({ axiosOpts: { doNotShowLoading: true }, refetchInterval })

const limitSwitchesArray = computed((): LimitSwitchItemProps[] => {
  const arr: LimitSwitchItemProps[] = []
  if (limitSwitches.value) {
    arr.push({
      name: 'Limit Switch 1',
      switchData: limitSwitches.value.limitSwitch1,
    })
  }
  return arr
})
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle>Limit Switches</CardTitle>
    </CardHeader>
    <CardContent>
      <div class="grid grid-cols-2 gap-2 text-sm">
        <template v-for="(limitSwitch, index) in limitSwitchesArray" :key="index">
          <LimitSwitchItem :name="limitSwitch.name" :switch-data="limitSwitch.switchData" />
          <div v-if="index !== limitSwitchesArray.length - 1" class="col-span-2 px-4">
            <Separator />
          </div>
        </template>
      </div>
    </CardContent>
  </Card>
</template>
