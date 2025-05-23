<script setup lang="ts">
import { FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
</script>

<template>
  <FormField v-slot="{ componentField }" name="inputs.durationMs">
    <FormItem>
      <FormLabel>Duration (ms)</FormLabel>
      <Input v-bind="componentField" type="number" placeholder="Enter duration" />
      <FormMessage />
    </FormItem>
  </FormField>
</template>
