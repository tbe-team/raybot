<script setup lang="ts">
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { Pagination, PaginationFirst, PaginationLast, PaginationList, PaginationNext, PaginationPrev } from '@/components/ui/pagination'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { useListQueuedCommandsQuery } from '@/composables/use-command'
import QueuedCommandItem from './QueuedCommandItem.vue'

const emit = defineEmits<{
  (e: 'viewDetails', commandId: number): void
}>()

const page = ref(1)
const pageSize = ref(5)
const { data: commands, refetch } = useListQueuedCommandsQuery(page, pageSize, {
  axiosOpts: {
    doNotShowLoading: true,
  },
})

function handlePageSizeChange() {
  page.value = 1
}
function handlePageChange(p: number) {
  page.value = p
}

const REFRESH_INTERVAL = 2000
setInterval(refetch, REFRESH_INTERVAL)
</script>

<template>
  <Card v-if="commands">
    <CardHeader class="pb-3">
      <CardTitle class="flex items-center justify-between">
        <span>Waiting Commands</span>
        <Badge variant="outline">
          {{ commands.totalItems }} commands
        </Badge>
      </CardTitle>
    </CardHeader>
    <CardContent class="space-y-3">
      <template v-if="commands.totalItems > 0">
        <QueuedCommandItem
          v-for="command in commands.items"
          :key="command.id"
          :command="command"
          @view-details="emit('viewDetails', command.id)"
        />
      </template>
      <div v-else class="py-6 text-center text-muted-foreground">
        No commands in queue
      </div>
    </CardContent>
    <CardFooter class="flex items-center justify-between pt-4 border-t">
      <div class="flex items-center gap-1 text-sm text-muted-foreground">
        <span class="text-sm sr-only sm:not-sr-only">
          Items per page:
        </span>
        <Select v-model="pageSize" @update:model-value="handlePageSizeChange">
          <SelectTrigger class="w-16 h-8 ml-2">
            <SelectValue placeholder="5" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="5">
              5
            </SelectItem>
            <SelectItem value="10">
              10
            </SelectItem>
            <SelectItem value="20">
              20
            </SelectItem>
          </SelectContent>
        </Select>
      </div>
      <Pagination
        :sibling-count="1"
        :items-per-page="pageSize"
        :total="commands.totalItems"
        :default-page="page"
        @update:page="handlePageChange"
      >
        <PaginationList class="flex items-center gap-1">
          <PaginationFirst />
          <PaginationPrev />
          <PaginationNext />
          <PaginationLast />
        </PaginationList>
      </Pagination>
    </CardFooter>
  </Card>
</template>
