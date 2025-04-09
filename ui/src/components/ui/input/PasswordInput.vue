<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { cn } from '@/lib/utils'
import { EyeIcon, EyeOffIcon } from 'lucide-vue-next'

const props = defineProps<{
  disabled?: boolean
  class?: string
}>()

const modelValue = defineModel<string>()
const showPassword = ref(false)

function togglePasswordVisibility() {
  showPassword.value = !showPassword.value
}
</script>

<template>
  <div class="relative">
    <Input
      v-model="modelValue"
      :type="showPassword ? 'text' : 'password'"
      :class="cn('pr-10 hide-password-toggle', props.class)"
    />
    <Button
      type="button"
      variant="ghost"
      size="sm"
      class="absolute top-0 right-0 h-full px-3 py-2 hover:bg-transparent"
      @click="togglePasswordVisibility"
    >
      <EyeIcon v-if="showPassword && !props.disabled" class="w-4 h-4" aria-hidden="true" />
      <EyeOffIcon v-else class="w-4 h-4" aria-hidden="true" />
      <span class="sr-only">{{ showPassword ? 'Hide password' : 'Show password' }}</span>
    </Button>
  </div>
</template>

<style lang="css" scoped>
.hide-password-toggle::-ms-reveal,
.hide-password-toggle::-ms-clear {
  visibility: hidden;
  pointer-events: none;
  display: none;
}
</style>
