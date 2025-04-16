<script setup lang="ts">
import { PageContainer } from '@/components/shared'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardFooter } from '@/components/ui/card'
import { useSystemRebootMutation } from '@/composables/use-system'
import { AlertTriangle, Loader, RefreshCw } from 'lucide-vue-next'

const { mutate, isPending } = useSystemRebootMutation()
const rebootInitiated = ref(false)

function handleReboot() {
  if (rebootInitiated.value)
    return

  rebootInitiated.value = true

  mutate(undefined, {
    onSuccess: () => {
      notification.success({
        message: 'System will reboot now. Please refresh the page after a moment.',
        title: 'Rebooting',
      })
    },
    onError: (error) => {
      rebootInitiated.value = false
      notification.error({
        message: error.message,
        title: 'Error',
      })
    },
  })
}
</script>

<template>
  <PageContainer>
    <div class="flex flex-col w-full">
      <div class="mb-6">
        <h1 class="text-xl font-semibold tracking-tight">
          System Reboot
        </h1>
        <p class="text-sm text-muted-foreground">
          Reboot the system to apply configuration changes
        </p>
      </div>

      <Card>
        <CardContent class="pt-6">
          <div class="flex flex-col gap-4">
            <div class="flex items-start gap-3 p-4 border border-yellow-200 rounded-lg bg-yellow-50 dark:bg-yellow-950/50 dark:border-yellow-800">
              <AlertTriangle class="w-6 h-6 text-yellow-600 dark:text-yellow-400 mt-0.5" />
              <div>
                <h3 class="mb-1 text-sm font-medium text-yellow-800 dark:text-yellow-300">
                  Warning
                </h3>
                <p class="text-sm text-yellow-600 dark:text-yellow-400">
                  Rebooting the system will temporarily interrupt all services. Make sure all important operations are completed before proceeding.
                </p>
              </div>
            </div>

            <div class="flex items-start gap-3 p-4 border border-blue-200 rounded-lg bg-blue-50 dark:bg-blue-950/50 dark:border-blue-800">
              <RefreshCw class="w-5 h-5 text-blue-600 dark:text-blue-400 mt-0.5" />
              <div>
                <h3 class="mb-1 text-sm font-medium text-blue-800 dark:text-blue-300">
                  What happens during restart
                </h3>
                <ul class="ml-4 space-y-1 text-sm text-blue-600 list-disc dark:text-blue-400">
                  <li>All system services will be stopped</li>
                  <li>Configuration changes will be applied</li>
                  <li>Services will be restarted with new settings</li>
                  <li>The process typically takes a few seconds</li>
                </ul>
              </div>
            </div>
          </div>
        </CardContent>
        <CardFooter class="flex justify-end pt-2 pb-6">
          <Button
            variant="destructive"
            size="lg"
            :disabled="isPending || rebootInitiated"
            @click="handleReboot"
          >
            <Loader v-if="isPending" class="w-4 h-4 mr-2 animate-spin" />
            <RefreshCw v-else class="w-4 h-4 mr-2" />
            Reboot
          </Button>
        </CardFooter>
      </Card>
    </div>
  </PageContainer>
</template>
