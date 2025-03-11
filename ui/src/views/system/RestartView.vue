<script setup lang="ts">
import { PageContainer } from '@/components/shared'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { useMutationSystemRestart } from '@/composables/use-system'
import { AlertTriangle, Loader } from 'lucide-vue-next'
import { push } from 'notivue'
import { ref } from 'vue'

const { mutate, isPending } = useMutationSystemRestart()
const restartInitiated = ref(false)

function handleRestart() {
  if (restartInitiated.value)
    return

  restartInitiated.value = true

  mutate(undefined, {
    onSuccess: () => {
      push.success({
        message: 'System will restart in 3 seconds. Please refresh the page after a moment.',
        title: 'Restarting',
      })
    },
    onError: (error) => {
      restartInitiated.value = false
      push.error({
        message: error.message,
        title: 'Error',
      })
    },
  })
}
</script>

<template>
  <PageContainer>
    <Card>
      <CardHeader>
        <CardTitle>System Restart</CardTitle>
        <CardDescription>
          Restart the system to apply configuration changes
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div class="flex items-center gap-2 p-4 border border-yellow-200 rounded-lg bg-yellow-50 dark:bg-yellow-950 dark:border-yellow-800">
          <AlertTriangle class="w-5 h-5 text-yellow-600 dark:text-yellow-400" />
          <p class="text-sm text-yellow-600 dark:text-yellow-400">
            Warning: Restarting the system will temporarily interrupt all services
          </p>
        </div>
      </CardContent>
      <CardFooter>
        <Button variant="destructive" :disabled="isPending || restartInitiated" @click="handleRestart">
          <Loader v-if="isPending" class="w-4 h-4 mr-2 animate-spin" />
          Restart System
        </Button>
      </CardFooter>
    </Card>
  </PageContainer>
</template>
