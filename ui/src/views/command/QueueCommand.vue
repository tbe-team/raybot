<script setup lang="ts">
import { PageContainer } from '@/components/shared'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Pagination, PaginationFirst, PaginationLast, PaginationList, PaginationNext, PaginationPrev } from '@/components/ui/pagination'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { useProcessingComandQuery, useQueuedComandQuery } from '@/composables/use-command'
import { formatTimeAgo } from '@/lib/date'
import { AlertCircle, AppWindow, Cloud, Loader } from 'lucide-vue-next'
import { h } from 'vue'
import CommandDetailDialog from './components/CommandDetailDialog.vue'
import StatusBadge from './components/commands-table/StatusBadge.vue'

const route = useRoute()
const router = useRouter()

const page = ref(Number(route.query.page) || 1)
const pageSize = ref(Number(route.query.pageSize) || 10)
const { data, isPending, isError, error, refetch } = useQueuedComandQuery(page, pageSize, { axiosOpts: { doNotShowLoading: true } })
const { data: processingCommand, isPending: isProcessingCommandPending, isError: isProcessingCommandError } = useProcessingComandQuery({ axiosOpts: { doNotShowLoading: true }, refetchInterval: 1000 })

function handlePageChange(p: number) {
  page.value = p
  router.replace({ query: { ...route.query, page: p.toString() } })
}

function getCommandTypeName(type: string) {
  return type.replace(/_/g, ' ').toLowerCase().replace(/\b\w/g, l => l.toUpperCase())
}

function handlePageSizeChange() {
  page.value = 1
  router.replace({ query: { ...route.query, pageSize: pageSize.value.toString(), page: '1' } })
}

watch(processingCommand, (newValue, oldValue) => {
  if (newValue !== oldValue) {
    refetch()
  }
})

function aliasCommandSource(source: string | undefined) {
  switch (source) {
    case 'APP':
      return h('span', { class: 'text-muted-foreground flex items-center gap-2' }, [h(AppWindow, { class: 'w-4 h-4' }), 'Application'])
    case 'CLOUD':
      return h('span', { class: 'text-muted-foreground flex items-center gap-2' }, [h(Cloud, { class: 'w-4 h-4' }), 'Cloud'])
    default:
      return h('span', { class: 'text-muted-foreground' }, source || 'N/A')
  }
}
</script>

<template>
  <PageContainer>
    <div>
      <h1 class="text-xl font-semibold">
        Queue command
      </h1>
      <p class="text-sm text-muted-foreground">
        Queue a command to the robot
      </p>
    </div>
    <div class="grid grid-cols-3 gap-4">
      <div class="col-span-2">
        <Card>
          <CardHeader>
            <CardTitle>
              Command Queue
            </CardTitle>
            <CardDescription>
              Current and queued commands
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div class="flex flex-col gap-4">
              <div class="flex flex-col gap-2">
                <p class="font-medium">
                  Current Command
                </p>
                <div class="w-full">
                  <Card v-if="isError" class="flex flex-col items-center gap-4 p-6 text-destructive">
                    <AlertCircle class="w-8 h-8" />
                    <div class="space-y-2 text-center">
                      <h2 class="text-lg font-semibold">
                        Failed to load state
                      </h2>
                      <p class="text-sm text-muted-foreground">
                        {{ error?.message || 'An unexpected error occurred' }}
                      </p>
                    </div>
                  </Card>
                  <div v-else-if="processingCommand?.items.length === 0" class="flex flex-col w-full gap-2 p-4 border rounded-md shadow-md bg-muted">
                    <div class="flex flex-col items-center gap-4 p-2">
                      <AlertCircle class="w-8 h-8 text-muted-foreground" />
                      <div class="space-y-2 text-center">
                        No command processing
                      </div>
                    </div>
                  </div>
                  <CommandDetailDialog v-else-if="processingCommand?.items[0]" :command="processingCommand?.items[0]">
                    <div class="flex flex-col w-full gap-2 p-6 border rounded-md shadow cursor-pointer bg-muted">
                      <div class="flex items-center justify-between w-full">
                        <p class="font-medium ">
                          {{ getCommandTypeName(processingCommand?.items[0].type) }}
                        </p>
                        <Badge class="text-blue-500 bg-blue-500/10 hover:bg-blue-500/20">
                          <Loader class="w-4 h-4 mr-1 duration-2000 animate-spin" />
                          Processing
                        </Badge>
                      </div>
                      <div class="flex gap-2">
                        <p class="flex items-center gap-1 text-sm text-muted-foreground">
                          <component :is="aliasCommandSource(processingCommand?.items[0].source)" />
                          •
                          <span class="text-xs text-muted-foreground">{{ formatTimeAgo(processingCommand?.items[0].createdAt || '') }}</span>
                        </p>
                      </div>
                    </div>
                  </CommandDetailDialog>
                </div>
              </div>
              <div class="flex flex-col gap-2">
                <p class="font-medium">
                  Waitting Command ({{ data?.totalItems }})
                </p>
                <div v-if="isPending" class="flex items-center justify-center w-full">
                  <Loader class="w-4 h-4 mr-1 duration-2000 animate-spin" />
                  Loading...
                </div>
                <div v-else-if="data?.items.length === 0" class="flex flex-col items-center justify-center w-full gap-2 p-6 border rounded-md shadow-md bg-muted">
                  <AlertCircle class="w-8 h-8 text-muted-foreground" />
                  No commands in queue
                </div>
                <CommandDetailDialog v-for="command in data?.items" :key="command.id" :command="command">
                  <div class="flex flex-col w-full gap-2 p-4 border rounded-md shadow-md cursor-pointer bg-muted">
                    <div class="flex items-center justify-between w-full">
                      <p class="font-medium">
                        {{ getCommandTypeName(command.type) }}
                      </p>
                      <StatusBadge :status="command.status" />
                    </div>
                    <div class="flex gap-2">
                      <p class="flex items-center gap-1 text-sm text-muted-foreground">
                        <component :is="aliasCommandSource(command.source)" />
                        •
                        <span class="text-xs text-muted-foreground">{{ formatTimeAgo(command.createdAt) }}</span>
                      </p>
                    </div>
                  </div>
                </CommandDetailDialog>
                <!-- Pagination -->
                <div class="flex items-center justify-end gap-4">
                  <div class="flex items-center gap-1">
                    <span class="text-sm sr-only sm:not-sr-only">
                      Items per page:
                    </span>
                    <Select v-model="pageSize" @update:model-value="handlePageSizeChange">
                      <SelectTrigger class="w-20">
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectGroup>
                          <SelectItem
                            v-for="option in [10, 20, 30]"
                            :key="option"
                            :value="option"
                          >
                            {{ option }}
                          </SelectItem>
                        </SelectGroup>
                      </SelectContent>
                    </Select>
                  </div>
                  <Pagination :items-per-page="pageSize" :total="data?.totalItems" :sibling-count="1" show-edges :default-page="page" @update:page="handlePageChange">
                    <PaginationList class="flex items-center gap-1">
                      <PaginationFirst />
                      <PaginationPrev />
                      <PaginationNext />
                      <PaginationLast />
                    </PaginationList>
                  </Pagination>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
      <div class="col-span-1">
        <Card>
          <CardHeader>
            <CardTitle>
              Create Command
            </CardTitle>
            <CardDescription>
              Add a new command to the queue
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Input placeholder="Command" />
          </CardContent>
          <CardFooter>
            <Button class="w-full">
              Add Command
            </Button>
          </CardFooter>
        </Card>
      </div>
    </div>
  </PageContainer>
</template>
