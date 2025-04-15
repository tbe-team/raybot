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

const steps = computed(() => [
  {
    step: 1,
    title: 'Created',
    description: formatDate(props.command.createdAt),
    icon: Circle,
    completed: !!props.command.createdAt,
  },
  {
    step: 2,
    title: 'Started',
    description: props.command.startedAt ? formatDate(props.command.startedAt) : '',
    icon: Loader2,
    completed: !!props.command.startedAt,
  },
  {
    step: 3,
    title: 'Completed',
    description: props.command.completedAt ? formatDate(props.command.completedAt) : '',
    icon: Check,
    completed: !!props.command.completedAt,
  },
  {
    step: 4,
    title: 'Last Updated',
    description: formatDate(props.command.updatedAt),
    icon: Circle,
    completed: !!props.command.updatedAt,
  },
])

const getVisibleSteps = computed(() => {
  if (!props.command)
    return []

  const visibleSteps = [steps.value[0]]

  if (props.command.startedAt) {
    visibleSteps.push(steps.value[1])
  }

  if (props.command.completedAt) {
    visibleSteps.push(steps.value[2])
  }

  if (props.command.updatedAt) {
    visibleSteps.push(steps.value[3])
  }

  return visibleSteps
})
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
