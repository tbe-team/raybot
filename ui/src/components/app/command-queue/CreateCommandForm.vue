<script setup lang="ts">
import type { CommandType } from '@/types/command'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { useCreateCommandMutation } from '@/composables/use-command'
import { RaybotError } from '@/types/error'
import { toTypedSchema } from '@vee-validate/zod'
import { ArrowDown, ArrowUp, Loader2, MapPin, Package, QrCode, StopCircle } from 'lucide-vue-next'
import { useForm } from 'vee-validate'
import { createCommandSchema } from './schemas'

const { values, handleSubmit, setFieldValue, resetForm } = useForm({
  validationSchema: toTypedSchema(createCommandSchema),
  initialValues: {
    type: 'STOP',
    inputs: {},
  },
})

const commandType = computed(() => values.type)
function setCommandType(type: CommandType) {
  setFieldValue('type', type)
}

const { mutate: createCommand, isPending } = useCreateCommandMutation()

const onSubmit = handleSubmit((values) => {
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
})

function clearForm() {
  resetForm()
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
              <Select
                :disabled="isPending"
                :model-value="commandType"
                @update:model-value="(val) => setCommandType(val as CommandType)"
              >
                <FormControl>
                  <SelectTrigger>
                    <SelectValue placeholder="Select command type" />
                  </SelectTrigger>
                </FormControl>
                <SelectContent>
                  <SelectItem value="STOP">
                    <div class="flex items-center gap-2">
                      <StopCircle class="w-4 h-4" />
                      <span>Stop</span>
                    </div>
                  </SelectItem>
                  <SelectItem value="MOVE_FORWARD">
                    <div class="flex items-center gap-2">
                      <ArrowUp class="w-4 h-4" />
                      <span>Move Forward</span>
                    </div>
                  </SelectItem>
                  <SelectItem value="MOVE_BACKWARD">
                    <div class="flex items-center gap-2">
                      <ArrowDown class="w-4 h-4" />
                      <span>Move Backward</span>
                    </div>
                  </SelectItem>
                  <SelectItem value="MOVE_TO">
                    <div class="flex items-center gap-2">
                      <MapPin class="w-4 h-4" />
                      <span>Move To</span>
                    </div>
                  </SelectItem>
                  <SelectItem value="CARGO_OPEN">
                    <div class="flex items-center gap-2">
                      <Package class="w-4 h-4" />
                      <span>Cargo Open</span>
                    </div>
                  </SelectItem>
                  <SelectItem value="CARGO_CLOSE">
                    <div class="flex items-center gap-2">
                      <Package class="w-4 h-4" />
                      <span>Cargo Close</span>
                    </div>
                  </SelectItem>
                  <SelectItem value="CARGO_LIFT">
                    <div class="flex items-center gap-2">
                      <Package class="w-4 h-4" />
                      <span>Cargo Lift</span>
                    </div>
                  </SelectItem>
                  <SelectItem value="CARGO_LOWER">
                    <div class="flex items-center gap-2">
                      <Package class="w-4 h-4" />
                      <span>Cargo Lower</span>
                    </div>
                  </SelectItem>
                  <SelectItem value="CARGO_CHECK_QR" disabled>
                    <div class="flex items-center gap-2">
                      <QrCode class="w-4 h-4" />
                      <span>Cargo Check QR</span>
                    </div>
                  </SelectItem>
                </SelectContent>
              </Select>
              <FormMessage />
            </FormItem>
          </FormField>

          <!-- Dynamic inputs based on command type -->
          <template v-if="commandType === 'MOVE_TO'">
            <FormField v-slot="{ componentField }" name="inputs.location">
              <FormItem>
                <FormLabel>Location</FormLabel>
                <Input v-bind="componentField" placeholder="Enter location" />
                <FormMessage />
              </FormItem>
            </FormField>
          </template>
          <template v-else-if="commandType === 'CARGO_CHECK_QR'">
            <FormField v-slot="{ componentField }" name="inputs.qrCode">
              <FormItem>
                <FormLabel>QR Code</FormLabel>
                <Input v-bind="componentField" placeholder="Enter QR code" />
                <FormMessage />
              </FormItem>
            </FormField>
          </template>
        </div>
      </CardContent>
      <CardFooter class="flex flex-col gap-2">
        <Button type="submit" class="w-full" :disabled="isPending">
          <Loader2 v-if="isPending" class="w-4 h-4 mr-2 animate-spin" />
          Add command
        </Button>
        <Button type="button" variant="outline" class="w-full" :disabled="isPending" @click="clearForm">
          Clear Form
        </Button>
      </CardFooter>
    </form>
  </Card>
</template>
