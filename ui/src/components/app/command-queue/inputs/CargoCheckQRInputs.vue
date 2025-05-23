<script setup lang="ts">
import { FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
</script>

<template>
  <FormField v-slot="{ componentField }" name="inputs.qrCode">
    <FormItem>
      <FormLabel>QR Code</FormLabel>
      <Input v-bind="componentField" type="text" placeholder="Enter QR code" />
      <FormMessage />
    </FormItem>
  </FormField>
</template>
