<script setup lang="ts">
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Separator } from '@/components/ui/separator'
</script>

<template>
  <FormField v-slot="{ componentField }" name="inputs.position">
    <FormItem>
      <FormLabel>Lower position (cm)</FormLabel>
      <Input v-bind="componentField" type="number" placeholder="Enter lower position" />
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

  <Separator />

  <div class="space-y-4">
    <h4 class="text-sm font-medium">
      Bottom Obstacle Tracking
    </h4>

    <div class="grid grid-cols-2 gap-4">
      <FormField v-slot="{ componentField }" name="inputs.bottomObstacleTracking.enterDistance">
        <FormItem>
          <FormLabel>Enter Distance (cm)</FormLabel>
          <FormControl>
            <Input v-bind="componentField" type="number" />
          </FormControl>
          <FormMessage />
        </FormItem>
      </FormField>

      <FormField v-slot="{ componentField }" name="inputs.bottomObstacleTracking.exitDistance">
        <FormItem>
          <FormLabel>Exit Distance (cm)</FormLabel>
          <FormControl>
            <Input v-bind="componentField" type="number" />
          </FormControl>
          <FormMessage />
        </FormItem>
      </FormField>
    </div>
  </div>
</template>
