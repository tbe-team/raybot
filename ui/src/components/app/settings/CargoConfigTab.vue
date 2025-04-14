<script setup lang="ts">
import { useCargoConfigQuery } from '@/composables/use-config'
import { Loader } from 'lucide-vue-next'
import CargoConfigForm from './CargoConfigForm.vue'

const { data, isPending } = useCargoConfigQuery()
</script>

<template>
  <div v-if="isPending" class="flex items-center justify-center h-full">
    <Loader class="w-4 h-4 animate-spin" />
  </div>
  <CargoConfigForm v-if="data" :initial-values="data" />
</template>
