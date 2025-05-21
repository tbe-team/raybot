<script setup lang="ts">
import type { CommandType } from '@/types/command'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { useCreateCommandMutation } from '@/composables/use-command'
import { RaybotError } from '@/types/error'

const { mutate: createCommand } = useCreateCommandMutation()
function handleCommand(command: CommandType) {
  const values = {
    type: command,
    inputs: {},
  }
  createCommand(values, {
    onSuccess: () => {
      notification.success('Command created successfully')
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
}
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle>Control</CardTitle>
      <CardDescription>Basic commands to control the robot</CardDescription>
    </CardHeader>
    <CardContent>
      <div class="h-full">
        <div class="grid grid-cols-3 grid-rows-3 gap-4">
          <Button class="flex items-center justify-center col-start-2 row-start-1 " @click="handleCommand('CARGO_LIFT')">
            Lift
          </Button>
          <Button class="flex items-center justify-center col-start-1 row-start-2" @click="handleCommand('MOVE_BACKWARD')">
            Backward
          </Button>
          <Button class="flex items-center justify-center col-start-2 row-start-2" @click="handleCommand('STOP_MOVEMENT')">
            Stop
          </Button>
          <Button class="flex items-center justify-center col-start-3 row-start-2" @click="handleCommand('MOVE_FORWARD')">
            Forward
          </Button>
          <Button class="flex items-center justify-center col-start-1 row-start-3" @click="handleCommand('CARGO_OPEN')">
            Open
          </Button>
          <Button class="flex items-center justify-center col-start-2 row-start-3" @click="handleCommand('CARGO_LOWER')">
            Lower
          </Button>
          <Button class="flex items-center justify-center col-start-3 row-start-3" @click="handleCommand('CARGO_CLOSE')">
            Close
          </Button>
        </div>
      </div>
    </CardContent>
  </Card>
</template>
