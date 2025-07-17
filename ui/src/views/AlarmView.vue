<script setup lang="ts">
import type { AlarmStatus } from '@/api/alarm'
import ActiveAlarmTabContent from '@/components/app/alarm/ActiveAlarmTabContent.vue'
import AlarmHistoryTabContent from '@/components/app/alarm/AlarmHistoryTabContent.vue'
import PageContainer from '@/components/shared/PageContainer.vue'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { useListAlarmsQuery } from '@/composables/use-alarm'

const route = useRoute()
const router = useRouter()
const page = ref(Number(route.query.page) || 1)
const pageSize = ref(Number(route.query.pageSize) || 10)
const tab = ref(route.query.tab as string | undefined ?? 'active')
const status = computed<AlarmStatus>(() => tab.value === 'history' ? 'DEACTIVE' : 'ACTIVE')

const { data, isFetching, isPending, isError, error, refetch } = useListAlarmsQuery(page, pageSize, status)

function handleTabChange(value: string | number) {
  tab.value = value.toString()
  page.value = 1
  pageSize.value = 10
  router.replace({ query: { tab: value } })
}

function handlePageChange(p: number) {
  page.value = p
  router.replace({ query: { ...route.query, page: p.toString() } })
}

function handlePageSizeChange(ps: number) {
  pageSize.value = ps
  page.value = 1
  router.replace({ query: { ...route.query, pageSize: ps.toString(), page: '1' } })
}
</script>

<template>
  <PageContainer>
    <div>
      <h1 class="text-xl font-semibold tracking-tight">
        Alarm
      </h1>
      <p class="text-sm text-muted-foreground">
        View alarms
      </p>
    </div>
    <Tabs :default-value="tab" @update:model-value="handleTabChange">
      <TabsList>
        <TabsTrigger value="active">
          Active
        </TabsTrigger>
        <TabsTrigger value="history">
          History
        </TabsTrigger>
      </TabsList>

      <TabsContent value="active">
        <ActiveAlarmTabContent
          :page="page"
          :page-size="pageSize"
          :data="data?.items ?? []"
          :total-items="data?.totalItems ?? 0"
          :is-fetching="isFetching"
          :is-pending="isPending"
          :is-error="isError"
          :error="error"
          @update:page="handlePageChange"
          @update:page-size="handlePageSizeChange"
          @update:data="refetch"
        />
      </TabsContent>
      <TabsContent value="history">
        <AlarmHistoryTabContent
          :page="page"
          :page-size="pageSize"
          :data="data?.items ?? []"
          :total-items="data?.totalItems ?? 0"
          :is-fetching="isFetching"
          :is-pending="isPending"
          :is-error="isError"
          :error="error"
          @update:page="handlePageChange"
          @update:page-size="handlePageSizeChange"
          @update:data="refetch"
        />
      </TabsContent>
    </Tabs>
  </PageContainer>
</template>
