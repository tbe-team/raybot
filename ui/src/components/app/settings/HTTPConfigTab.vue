<script setup lang="ts">
import { useHTTPConfigQuery } from '@/composables/use-config'
import { Loader } from 'lucide-vue-next'
import HTTPConfigForm from './HTTPConfigForm.vue'

const { data, isPending } = useHTTPConfigQuery()
</script>

<template>
  <div v-if="isPending" class="flex items-center justify-center h-full">
    <Loader class="w-4 h-4 animate-spin" />
  </div>
  <HTTPConfigForm v-if="data" :initial-values="data" />
</template>
