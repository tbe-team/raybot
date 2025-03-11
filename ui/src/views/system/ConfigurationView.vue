<script setup lang="ts">
import { useQuerySystemConfig } from '@/composables/use-system'
import { SystemConfigForm } from '@/views/system/components/system-config-form'
import { Loader } from 'lucide-vue-next'

const { isPending, data, isError, error } = useQuerySystemConfig()
</script>

<template>
  <div class="flex items-center justify-center w-full">
    <div v-if="isPending" class="flex items-center justify-center gap-4 pt-20 text-2xl">
      <Loader class="animate-spin" /> Loading
    </div>
    <div v-else-if="isError" class="pt-20 text-2xl">
      Error: {{ error }}
    </div>
    <div v-else-if="!data" class="pt-20 text-2xl">
      System Config Not Found
    </div>
    <div v-else class="flex flex-col w-full gap-4">
      <SystemConfigForm :system-config="data" />
    </div>
  </div>
</template>
