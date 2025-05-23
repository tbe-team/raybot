<script setup lang="ts">
import { ArrowLeft, ArrowRight } from 'lucide-vue-next'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
</script>

<template>
  <FormField v-slot="{ componentField }" name="inputs.location">
    <FormItem>
      <FormLabel>Location</FormLabel>
      <Input v-bind="componentField" type="text" placeholder="Enter location" />
      <FormMessage />
    </FormItem>
  </FormField>
  <FormField v-slot="{ componentField }" name="inputs.direction">
    <FormItem>
      <FormLabel>Move direction</FormLabel>
      <Select v-bind="componentField">
        <FormControl>
          <SelectTrigger>
            <SelectValue placeholder="Select direction" />
          </SelectTrigger>
        </FormControl>
        <SelectContent>
          <SelectItem value="FORWARD">
            <div class="flex items-center gap-2">
              <ArrowRight class="w-4 h-4" />
              <span>Forward</span>
            </div>
          </SelectItem>
          <SelectItem value="BACKWARD">
            <div class="flex items-center gap-2">
              <ArrowLeft class="w-4 h-4" />
              <span>Backward</span>
            </div>
          </SelectItem>
        </SelectContent>
      </Select>
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
