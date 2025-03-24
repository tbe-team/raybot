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
  <Card>
    <CardHeader>
      <CardTitle>GRPC Configuration</CardTitle>
    </CardHeader>
    <CardContent class="space-y-4">
      <section class="space-y-3">
        <span class="font-medium">Server</span>
        <div class="pl-4 space-y-2">
          <FormField v-slot="{ value, handleChange }" name="grpc.server.enable">
            <FormItem>
              <div class="flex items-center space-x-2">
                <FormLabel>Enable</FormLabel>
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
      </section>
      <section class="space-y-3">
        <span class="font-medium">Cloud</span>
        <div class="pl-4 space-y-2">
          <FormField v-slot="{ componentField }" name="grpc.cloud.address">
            <FormItem>
              <FormLabel>Address</FormLabel>
              <FormControl>
                <Input
                  v-bind="componentField"
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          </FormField>
        </div>
      </section>
    </CardContent>
  </Card>
</template>
