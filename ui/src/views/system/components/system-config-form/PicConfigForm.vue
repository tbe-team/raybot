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
import { serialDataBitsOption, serialParityOption, serialStopBitsOption } from './serial-option'
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
                  <SelectItem v-for="option in serialDataBitsOption" :key="option.value" :value="option.value">
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
                  <SelectItem v-for="option in serialStopBitsOption" :key="option.value" :value="option.value">
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
                  <SelectItem v-for="option in serialParityOption" :key="option.value" :value="option.value">
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
