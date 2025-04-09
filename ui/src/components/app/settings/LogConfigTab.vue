<script setup lang="ts">
import { useLogConfigQuery } from '@/composables/use-config'
import { Loader } from 'lucide-vue-next'
import LogConfigForm from './LogConfigForm.vue'

const { data, isPending } = useLogConfigQuery()
</script>

<template>
  <div v-if="isPending" class="flex items-center justify-center h-full">
    <Loader class="w-4 h-4 animate-spin" />
  </div>
  <LogConfigForm v-if="data" :initial-values="data" />
</template>
