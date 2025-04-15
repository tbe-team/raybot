<script setup lang="ts">
import type { Command } from '@/types/command'
import { Button } from '@/components/ui/button'
import { Stepper, StepperDescription, StepperItem, StepperSeparator, StepperTitle } from '@/components/ui/stepper'
import { formatDate } from '@/lib/date'
import { Check, Circle, Loader2 } from 'lucide-vue-next'

interface Props {
  command: Command
}
const props = defineProps<Props>()
const command = toRef(props, 'command')

const steps = [
  {
    step: 1,
    title: 'Created',
    description: formatDate(command.value.createdAt),
    icon: Circle,
    completed: !!command.value.createdAt,
  },
  {
    step: 2,
    title: 'Started',
    description: command.value.startedAt ? formatDate(command.value.startedAt) : '',
    icon: Loader2,
    completed: !!command.value.startedAt,
  },
  {
    step: 3,
    title: 'Completed',
    description: command.value.completedAt ? formatDate(command.value.completedAt) : '',
    icon: Check,
    completed: !!command.value.completedAt,
  },
  {
    step: 4,
    title: 'Last Updated',
    description: formatDate(command.value.updatedAt),
    icon: Circle,
    completed: !!command.value.updatedAt,
  },
]

const getVisibleSteps = computed(() => {
  if (!command.value)
    return []

  const visibleSteps = [steps[0]]

  if (command.value.startedAt) {
    visibleSteps.push(steps[1])
  }

  if (command.value.completedAt) {
    visibleSteps.push(steps[2])
  }

  if (command.value.updatedAt) {
    visibleSteps.push(steps[3])
  }

  return visibleSteps
})

watch(command, () => {
  if (!command.value)
    return

  steps[0].description = formatDate(command.value.createdAt)
  steps[0].completed = !!command.value.createdAt

  steps[1].description = command.value.startedAt ? formatDate(command.value.startedAt) : ''
  steps[1].completed = !!command.value.startedAt

  steps[2].description = command.value.completedAt ? formatDate(command.value.completedAt) : ''
  steps[2].completed = !!command.value.completedAt

  steps[3].description = formatDate(command.value.updatedAt)
  steps[3].completed = !!command.value.updatedAt
}, { immediate: true })
</script>

<template>
  <Stepper orientation="vertical" class="flex flex-col justify-start w-full max-w-md gap-8 mx-auto">
    <StepperItem
      v-for="step in getVisibleSteps"
      :key="step.step"
      class="relative flex items-start w-full gap-6"
      :step="step.step"
    >
      <StepperSeparator
        v-if="step.step !== steps[steps.length - 1].step"
        class="absolute left-[18px] top-[38px] block h-[90%] w-0.5 shrink-0 rounded-full !bg-primary !opacity-100"
      />
      <Button
        variant="default"
        size="icon"
        class="z-10 !p-0 rounded-full shrink-0"
      >
        <component :is="step.icon" class="!w-4 !h-4" />
      </Button>
      <div class="flex flex-col gap-1">
        <StepperTitle
          class="text-sm font-semibold transition lg:text-base"
        >
          {{ step.title }}
        </StepperTitle>
        <StepperDescription
          class="text-xs transition text-muted-foreground md:not-sr-only lg:text-sm"
        >
          {{ step.description }}
        </StepperDescription>
      </div>
    </StepperItem>
  </Stepper>
</template>
