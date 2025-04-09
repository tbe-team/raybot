<script setup lang="ts">
import { useCloudConfigQuery } from '@/composables/use-config'
import { Loader } from 'lucide-vue-next'
import CloudConfigForm from './CloudConfigForm.vue'

const { data, isPending } = useCloudConfigQuery()
</script>

<template>
  <div v-if="isPending" class="flex items-center justify-center h-full">
    <Loader class="w-4 h-4 animate-spin" />
  </div>
  <CloudConfigForm v-if="data" :initial-values="data" />
</template>
