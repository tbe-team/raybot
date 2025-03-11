<script setup lang="ts">
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Switch } from '@/components/ui/switch'
</script>

<template>
  <Card class="rounded-sm">
    <CardHeader>
      <CardTitle>Http Configuration</CardTitle>
    </CardHeader>
    <CardContent class="space-y-4">
      <div class="space-y-2">
        <FormField v-slot="{ componentField }" name="http.port">
          <FormItem>
            <FormLabel>Port</FormLabel>
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
      <div class="flex items-center space-x-2">
        <FormField v-slot="{ value, handleChange }" name="http.enableSwagger">
          <FormItem>
            <div class="flex items-center space-x-2">
              <FormLabel>Enable Swagger</FormLabel>
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
