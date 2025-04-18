<script setup lang="ts">
import { useColorMode } from '@vueuse/core'
import ConfirmationDialog from './components/shared/ConfirmationDialog.vue'

useColorMode()
</script>

<template>
  <RouterView />
  <ConfirmationDialog />
</template>
