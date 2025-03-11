<script setup lang="ts">
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
</script>

<template>
  <Card class="rounded-sm">
    <CardHeader>
      <CardTitle>GRPC Configuration</CardTitle>
    </CardHeader>
    <CardContent class="space-y-4">
      <div class="space-y-2">
        <FormField v-slot="{ componentField }" name="grpc.port">
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
    </CardContent>
  </Card>
</template>
