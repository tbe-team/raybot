<script setup lang="ts">
import { Loader } from 'lucide-vue-next'
import { useCommandConfigQuery } from '@/composables/use-config'
import CommandsConfigForm from './CommandConfigForm.vue'

const { data, isPending } = useCommandConfigQuery()
</script>

<template>
  <div v-if="isPending" class="flex items-center justify-center h-full">
    <Loader class="w-4 h-4 animate-spin" />
  </div>
  <CommandsConfigForm v-if="data" :initial-values="data" />
</template>
