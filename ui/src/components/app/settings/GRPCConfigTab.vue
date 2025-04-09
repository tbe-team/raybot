<script setup lang="ts">
import { useGRPCConfigQuery } from '@/composables/use-config'
import { Loader } from 'lucide-vue-next'
import GRPCConfigForm from './GRPCConfigForm.vue'

const { data, isPending } = useGRPCConfigQuery()
</script>

<template>
  <div v-if="isPending" class="flex items-center justify-center h-full">
    <Loader class="w-4 h-4 animate-spin" />
  </div>
  <GRPCConfigForm v-if="data" :initial-values="data" />
</template>
