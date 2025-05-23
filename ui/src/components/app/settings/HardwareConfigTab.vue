<script setup lang="ts">
import { Loader } from 'lucide-vue-next'
import { useHardwareConfigQuery } from '@/composables/use-config'
import HardwareConfigForm from './HardwareConfigForm.vue'

const { data, isPending } = useHardwareConfigQuery()
</script>

<template>
  <div v-if="isPending" class="flex items-center justify-center h-full">
    <Loader class="w-4 h-4 animate-spin" />
  </div>
  <HardwareConfigForm v-if="data" :initial-values="data" />
</template>
