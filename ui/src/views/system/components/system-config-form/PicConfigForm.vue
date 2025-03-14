<script setup lang="ts">
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'

const picSerialParityOption = [
  { label: 'None', value: 'none' },
  { label: 'Even', value: 'even' },
  { label: 'Odd', value: 'odd' },
]

const dataBitsOption = [
  { label: '5', value: 5 },
  { label: '6', value: 6 },
  { label: '7', value: 7 },
  { label: '8', value: 8 },
]
const stopBitsOption = [
  { label: '1', value: 1 },
  { label: '1.5', value: 1.5 },
  { label: '2', value: 2 },
]
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle>PIC Configuration</CardTitle>
    </CardHeader>
    <CardContent class="space-y-4">
      <div class="space-y-2">
        <FormField v-slot="{ componentField }" name="pic.serial.port">
          <FormItem>
            <FormLabel>Port</FormLabel>
            <FormControl>
              <Input
                v-bind="componentField"
              />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>
      </div>
      <div class="space-y-2">
        <FormField v-slot="{ componentField }" name="pic.serial.baudRate">
          <FormItem>
            <FormLabel>Baud Rate</FormLabel>
            <FormControl>
              <Input
                v-bind="componentField"
                type="number"
              />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>
      </div>
      <div class="space-y-2">
        <FormField v-slot="{ componentField }" name="pic.serial.dataBits">
          <FormItem>
            <FormLabel>Data Bits</FormLabel>
            <Select
              v-bind="componentField"
            >
              <FormControl>
                <SelectTrigger>
                  <SelectValue placeholder="Select data bits" />
                </SelectTrigger>
              </FormControl>
              <SelectContent>
                <SelectGroup>
                  <SelectItem v-for="option in dataBitsOption" :key="option.value" :value="option.value">
                    {{ option.label }}
                  </SelectItem>
                </SelectGroup>
              </SelectContent>
            </Select>
            <FormMessage />
          </FormItem>
        </FormField>
      </div>
      <div class="space-y-2">
        <FormField v-slot="{ componentField }" name="pic.serial.stopBits">
          <FormItem>
            <FormLabel>Stop Bits</FormLabel>
            <Select
              v-bind="componentField"
            >
              <FormControl>
                <SelectTrigger>
                  <SelectValue placeholder="Select stop bits" />
                </SelectTrigger>
              </FormControl>
              <SelectContent>
                <SelectGroup>
                  <SelectItem v-for="option in stopBitsOption" :key="option.value" :value="option.value">
                    {{ option.label }}
                  </SelectItem>
                </SelectGroup>
              </SelectContent>
            </Select>
            <FormMessage />
          </FormItem>
        </FormField>
      </div>
      <div class="space-y-2">
        <FormField v-slot="{ componentField }" name="pic.serial.parity">
          <FormItem>
            <FormLabel>Parity</FormLabel>
            <Select
              v-bind="componentField"
            >
              <FormControl>
                <SelectTrigger>
                  <SelectValue placeholder="Select parity" />
                </SelectTrigger>
              </FormControl>
              <SelectContent>
                <SelectGroup>
                  <SelectItem v-for="option in picSerialParityOption" :key="option.value" :value="option.value">
                    {{ option.label }}
                  </SelectItem>
                </SelectGroup>
              </SelectContent>
            </Select>
            <FormMessage />
          </FormItem>
        </FormField>
      </div>
      <div class="space-y-2">
        <FormField v-slot="{ componentField }" name="pic.serial.readTimeout">
          <FormItem>
            <FormLabel>Read Timeout</FormLabel>
            <FormControl>
              <Input
                v-bind="componentField"
                type="number"
              />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>
      </div>
    </CardContent>
  </Card>
</template>
