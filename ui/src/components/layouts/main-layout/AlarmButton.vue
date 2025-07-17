<script setup lang="ts">
import { Bell } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip'
import { useCountActiveAlarmsQuery } from '@/composables/use-alarm'

const REFRESH_INTERVAL = 5000
const router = useRouter()
const { data, isError } = useCountActiveAlarmsQuery({
  axiosOpts: { doNotShowLoading: true },
  refetchInterval: REFRESH_INTERVAL,
})
</script>

<template>
  <TooltipProvider>
    <Tooltip>
      <TooltipTrigger as-child>
        <Button
          class="relative rounded-lg bg-muted hover:bg-muted-hover" variant="ghost" size="icon" @click="() => {
            router.push('/alarms')
          }"
        >
          <Bell class="size-5" />
          <span v-if="data && data.count > 0" class="block absolute top-0 right-0 rounded-full border bg-destructive dark:border-white size-2" />
        </Button>
      </TooltipTrigger>
      <TooltipContent>
        <p v-if="data && data.count > 0">
          The system has {{ data.count }} alarms. Click to view
        </p>
        <p v-if="data && data.count === 0">
          No alarm found
        </p>
        <p v-else-if="isError">
          Failed to load alarms
        </p>
      </TooltipContent>
    </Tooltip>
  </TooltipProvider>
</template>
