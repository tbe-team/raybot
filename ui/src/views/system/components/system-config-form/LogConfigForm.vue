<script setup lang="ts">
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Switch } from '@/components/ui/switch'

const logLevelOption = [
  { label: 'Debug', value: 'debug' },
  { label: 'Info', value: 'info' },
  { label: 'Warn', value: 'warn' },
  { label: 'Error', value: 'error' },
]
const logFormatOption = [
  { label: 'JSON', value: 'json' },
  { label: 'Text', value: 'text' },
]
</script>

<template>
  <Card class="rounded-sm">
    <CardHeader>
      <CardTitle>Log Configuration</CardTitle>
    </CardHeader>
    <CardContent class="space-y-4">
      <div class="space-y-2">
        <FormField v-slot="{ componentField }" name="log.level">
          <FormItem>
            <FormLabel>Level</FormLabel>
            <Select
              v-bind="componentField"
            >
              <FormControl>
                <SelectTrigger>
                  <SelectValue placeholder="Select log level" />
                </SelectTrigger>
              </FormControl>
              <SelectContent>
                <SelectGroup>
                  <SelectItem v-for="option in logLevelOption" :key="option.value" :value="option.value">
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
        <FormField v-slot="{ componentField }" name="log.format">
          <FormItem>
            <FormLabel>Log Format</FormLabel>
            <Select
              v-bind="componentField"
            >
              <FormControl>
                <SelectTrigger>
                  <SelectValue placeholder="Select log format" />
                </SelectTrigger>
              </FormControl>
              <SelectContent>
                <SelectGroup>
                  <SelectItem v-for="option in logFormatOption" :key="option.value" :value="option.value">
                    {{ option.label }}
                  </SelectItem>
                </SelectGroup>
              </SelectContent>
            </Select>
            <FormMessage />
          </FormItem>
        </FormField>
      </div>
      <div class="flex items-center space-x-2">
        <FormField v-slot="{ value, handleChange }" name="log.addSource">
          <FormItem>
            <div class="flex items-center space-x-2">
              <FormLabel>Add Source</FormLabel>
              <FormControl>
                <Switch
                  class="!mt-0"
                  :model-value="value"
                  @update:model-value="handleChange"
                />
              </FormControl>
            </div>
            <FormMessage />
          </FormItem>
        </FormField>
      </div>
    </CardContent>
  </Card>
</template>
