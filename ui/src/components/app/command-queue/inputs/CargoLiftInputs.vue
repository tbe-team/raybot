<script setup lang="ts">
import { FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
</script>

<template>
  <FormField v-slot="{ componentField }" name="inputs.position">
    <FormItem>
      <FormLabel>Lift position (cm)</FormLabel>
      <Input v-bind="componentField" type="number" placeholder="Enter lift position" />
      <FormMessage />
    </FormItem>
  </FormField>
  <FormField v-slot="{ componentField }" name="inputs.motorSpeed">
    <FormItem>
      <FormLabel>Motor Speed (0-100%)</FormLabel>
      <Input v-bind="componentField" type="number" placeholder="Enter motor speed" />
      <FormMessage />
    </FormItem>
  </FormField>
</template>
