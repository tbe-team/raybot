<script setup lang="ts">
import { Loader } from 'lucide-vue-next'
import { useWifiConfigQuery } from '@/composables/use-config'
import WifiConfigForm from './WifiConfigForm.vue'

const { data, isPending } = useWifiConfigQuery()
</script>

<template>
  <div v-if="isPending" class="flex items-center justify-center h-full">
    <Loader class="w-4 h-4 animate-spin" />
  </div>
  <WifiConfigForm v-if="data" :initial-values="data" />
</template>
