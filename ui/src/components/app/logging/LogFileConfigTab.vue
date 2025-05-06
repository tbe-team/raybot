<script setup lang="ts">
import { useLogConfigQuery } from '@/composables/use-config'
import { Loader } from 'lucide-vue-next'
import LogFileConfigForm from './LogFileConfigForm.vue'

const { data, isPending } = useLogConfigQuery()
</script>

<template>
  <div v-if="isPending" class="flex items-center justify-center h-full">
    <Loader class="w-4 h-4 animate-spin" />
  </div>
  <LogFileConfigForm v-if="data" :log-config="data" />
</template>
