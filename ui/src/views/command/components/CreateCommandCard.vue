<script setup lang="ts">
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import type { CommandType, CommandInputMap } from '@/types/command';
import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { ref, computed } from 'vue';
import commandsAPI from '@/api/commands';
import { useCreateCommandMutation } from '@/composables/use-command';

const selectedType = ref<CommandType | undefined>()
const inputs = ref<Record<string, string>>({})
const  commandMutation = useCreateCommandMutation()

const inputFields = computed(() => {
  switch (selectedType.value) {
    case 'MOVE_TO':
      return [{ key: 'location', label: 'Location' }]
    case 'CARGO_CHECK_QR':
      return [{ key: 'qrCode', label: 'QR Code' }]
    default:
      return []
  }
})

function getCommandTypeName(type: string) {
  return type.replace(/_/g, ' ').toLowerCase().replace(/\b\w/g, l => l.toUpperCase())
}

async function handleSubmit() {
  try {
    const commandInputs = inputFields.value.reduce((acc, field) => {
      acc[field.key] = inputs.value[field.key]
      return acc
    }, {} as Record<string, string>)

    commandMutation.mutate({
      type: selectedType.value as CommandType,
      inputs: commandInputs
    },{
      onSuccess: () => {
        selectedType.value = undefined
        inputs.value = {}
      }
    })
  } catch (error) {
    console.error('Failed to create command:', error)
  }
}
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle>
        Create Command
      </CardTitle>
      <CardDescription>
        Add a new command to the queue
      </CardDescription>
    </CardHeader>
    <CardContent class="space-y-4">
      <div class="space-y-2">
        <label class="text-sm font-medium">Command Type</label>
        <Select v-model="selectedType">
          <SelectTrigger>
            <SelectValue placeholder="Select Command Type" />
          </SelectTrigger>
          <SelectContent>
            <SelectGroup>
              <SelectItem v-for="type in ['STOP', 'MOVE_FORWARD', 'MOVE_BACKWARD', 'MOVE_TO', 'CARGO_OPEN', 'CARGO_CLOSE', 'CARGO_LIFT', 'CARGO_LOWER', 'CARGO_CHECK_QR']" :key="type" :value="type">
                {{ getCommandTypeName(type) }}
              </SelectItem>
            </SelectGroup>
          </SelectContent>
        </Select>
      </div>

      <div v-for="field in inputFields" :key="field.key" class="space-y-2">
        <label class="text-sm font-medium">{{ field.label }}</label>
        <Input v-model="inputs[field.key]" :placeholder="`Enter ${field.label.toLowerCase()}`" />
      </div>
    </CardContent>
    <CardFooter>
      <Button class="w-full" @click="handleSubmit">
        Add Command
      </Button>
    </CardFooter>
  </Card>
</template>