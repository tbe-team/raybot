<script setup lang="ts">
import type { SystemConfig } from '@/types/system-config'
import { Button } from '@/components/ui/button'
import {
  CardDescription,
  CardTitle,
} from '@/components/ui/card'
import { useMutationSystemConfig } from '@/composables/use-system'
import { HTTPError, RaybotError } from '@/types/error'
import GrpcConfigForm from '@/views/system/components/system-config-form/GrpcConfigForm.vue'
import HttpConfigForm from '@/views/system/components/system-config-form/HttpConfigForm.vue'
import LogConfigForm from '@/views/system/components/system-config-form/LogConfigForm.vue'
import PicConfigForm from '@/views/system/components/system-config-form/PicConfigForm.vue'
import { systemConfigSchema } from '@/views/system/components/system-config-form/shema'
import { toTypedSchema } from '@vee-validate/zod'
import { Loader } from 'lucide-vue-next'
import { push } from 'notivue'
import { useForm } from 'vee-validate'

interface Props {
  systemConfig: SystemConfig
}
const props = defineProps<Props>()
const { mutate, isPending } = useMutationSystemConfig()

const form = useForm({
  validationSchema: toTypedSchema(systemConfigSchema),
  initialValues: props.systemConfig,
})

watch(() => props.systemConfig, (data) => {
  form.setValues(data)
})

const onSubmit = form.handleSubmit((values) => {
  mutate(values, {
    onSuccess: () => {
      push.success({ message: 'Update successfully', title: 'Success!' })
    },
    onError: (error) => {
      if (error instanceof RaybotError)
        push.error({ message: error.message, title: error.errorCode })
      else if (error instanceof HTTPError)
        push.error({ message: error.message, title: error.status.toString() })
      else
        push.error({ message: error.message, title: error.name })
    },
  })
})
</script>

<template>
  <form @submit="onSubmit">
    <div class="flex flex-col items-start justify-between gap-4 mb-6 sm:flex-row sm:items-center sm:gap-0">
      <div>
        <h1 class="text-xl font-semibold">
          System config
        </h1>
        <p class="text-sm text-muted-foreground">
          Important: Any changes made here require a
          <RouterLink to="/system/restart" class="text-blue-500 underline">
            system restart
          </RouterLink>
          to take effect
        </p>
      </div>
      <Button type="submit" class="w-fit">
        <Loader v-if="isPending" class="w-4 h-4 animate-spin" />
        Save
      </Button>
    </div>
    <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
      <GrpcConfigForm />
      <HttpConfigForm />
      <LogConfigForm />
      <PicConfigForm />
    </div>
  </form>
</template>
