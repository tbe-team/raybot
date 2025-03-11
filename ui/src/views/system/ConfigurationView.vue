<script setup lang="ts">
import { PageContainer } from '@/components/shared'
import { Card } from '@/components/ui/card'
import { useQuerySystemConfig } from '@/composables/use-system'
import { SystemConfigForm } from '@/views/system/components/system-config-form'
import { AlertCircle, Loader } from 'lucide-vue-next'

const { isPending, data, isError, error } = useQuerySystemConfig()
</script>

<template>
  <PageContainer>
    <div v-if="isPending" class="flex flex-col items-center justify-center gap-4 pt-20">
      <div class="flex items-center gap-4">
        <Loader class="w-8 h-8 animate-spin text-muted-foreground" />
      </div>
      <p class="text-lg text-muted-foreground">
        Loading system configuration...
      </p>
    </div>

    <div v-else-if="isError" class="flex flex-col items-center justify-center gap-4 pt-20">
      <Card class="flex flex-col items-center gap-4 p-6 text-destructive">
        <AlertCircle class="w-8 h-8" />
        <div class="space-y-2 text-center">
          <h2 class="text-lg font-semibold">
            Failed to load configuration
          </h2>
          <p class="text-sm text-muted-foreground">
            {{ error?.message || 'An unexpected error occurred' }}
          </p>
        </div>
      </Card>
    </div>

    <div v-else-if="!data" class="flex flex-col items-center justify-center gap-4 pt-20">
      <Card class="flex flex-col items-center gap-4 p-6">
        <AlertCircle class="w-8 h-8 text-muted-foreground" />
        <div class="space-y-2 text-center">
          <h2 class="text-lg font-semibold">
            No Configuration Found
          </h2>
          <p class="text-sm text-muted-foreground">
            The system configuration appears to be empty
          </p>
        </div>
      </Card>
    </div>

    <div v-else class="flex flex-col w-full gap-4">
      <SystemConfigForm :system-config="data" />
    </div>
  </PageContainer>
</template>
