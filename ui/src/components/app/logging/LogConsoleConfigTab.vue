<script setup lang="ts">
import { useLogConfigQuery } from '@/composables/use-config'
import { Loader } from 'lucide-vue-next'
import LogConsoleConfigForm from './LogConsoleConfigForm.vue'

const { data, isPending } = useLogConfigQuery()
</script>

<template>
  <div v-if="isPending" class="flex items-center justify-center h-full">
    <Loader class="w-4 h-4 animate-spin" />
  </div>
  <LogConsoleConfigForm v-if="data" :log-config="data" />
</template>
