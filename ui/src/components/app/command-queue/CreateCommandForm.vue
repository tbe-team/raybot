<script setup lang="ts">
import type { CommandType } from '@/types/command'
import { useQueryClient } from '@tanstack/vue-query'
import { toTypedSchema } from '@vee-validate/zod'
import { Loader2 } from 'lucide-vue-next'
import { useForm } from 'vee-validate'
import { useCommandConfig } from '@/components/app/command-queue/use-command-config'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { COMMAND_QUEUE_QUERY_KEY, CURRENT_PROCESSING_COMMAND_QUERY_KEY, useCreateCommandMutation } from '@/composables/use-command'
import { RaybotError } from '@/types/error'
import CommandTypeSelect from './CommandTypeSelect.vue'
import DynamicInputs from './inputs/DynamicInputs.vue'
import { createCommandSchema } from './schemas'

const queryClient = useQueryClient()

const { values, handleSubmit, setFieldValue } = useForm({
  validationSchema: toTypedSchema(createCommandSchema),
  initialValues: {
    type: 'STOP_MOVEMENT',
    inputs: {},
  },
})

const commandType = computed(() => values.type!)

const { mutate: createCommand, isPending } = useCreateCommandMutation()

const { commandConfig, updateCommandConfigFromInputs } = useCommandConfig()

const onSubmit = handleSubmit((values) => {
  createCommand(values, {
    onSuccess: () => {
      notification.success('Command created successfully')
      queryClient.invalidateQueries({ queryKey: [COMMAND_QUEUE_QUERY_KEY] })
      queryClient.invalidateQueries({ queryKey: [CURRENT_PROCESSING_COMMAND_QUERY_KEY] })
    },
    onError: (error) => {
      if (error instanceof RaybotError) {
        notification.error({
          title: error.errorCode,
          message: error.message,
        })
      }
      else {
        notification.error('Failed to create command')
      }
    },
  })

  updateCommandConfigFromInputs(commandType.value, values.inputs)
})

function setCommandType(type: CommandType) {
  setFieldValue('type', type)

  const configMap: Partial<Record<CommandType, unknown>> = {
    MOVE_TO: commandConfig.value.moveTo,
    MOVE_FORWARD: commandConfig.value.moveForward,
    MOVE_BACKWARD: commandConfig.value.moveBackward,
    CARGO_OPEN: commandConfig.value.cargoOpen,
    CARGO_CLOSE: commandConfig.value.cargoClose,
    CARGO_LIFT: commandConfig.value.cargoLift,
    CARGO_LOWER: commandConfig.value.cargoLower,
  }
  const inputs = configMap[type] ?? {}
  setFieldValue('inputs', { ...inputs })
}
</script>

<template>
  <Card class="sticky top-6">
    <CardHeader>
      <CardTitle>Create command</CardTitle>
    </CardHeader>

    <form @submit.prevent="onSubmit">
      <CardContent>
        <div class="space-y-4">
          <FormField name="type">
            <FormItem>
              <FormLabel>Command type</FormLabel>
              <CommandTypeSelect
                :disabled="isPending"
                :model-value="commandType"
                @update:model-value="(val) => setCommandType(val as CommandType)"
              />
              <FormMessage />
            </FormItem>
          </FormField>
          <DynamicInputs :command-type="commandType" />
        </div>
      </CardContent>

      <CardFooter class="flex flex-col gap-2">
        <Button type="submit" class="w-full" :disabled="isPending">
          <Loader2 v-if="isPending" class="w-4 h-4 mr-2 animate-spin" />
          Create
        </Button>
      </CardFooter>
    </form>
  </Card>
</template>
