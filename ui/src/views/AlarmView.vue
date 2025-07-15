<script setup lang="ts">
import AlarmActiveTab from '@/components/app/alarm/AlarmActiveTab.vue'
import AlarmHistoryTab from '@/components/app/alarm/AlarmHistoryTab.vue'
import PageContainer from '@/components/shared/PageContainer.vue'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'

const route = useRoute()
const router = useRouter()
const page = ref(Number(route.query.page) || 1)
const pageSize = ref(Number(route.query.pageSize) || 10)
const tab = route.query.tab as string | undefined ?? 'active'

function handleTabChange(value: string | number) {
  router.replace({ query: { tab: value } })
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
        <AlarmActiveTab v-model:page="page" v-model:page-size="pageSize" />
      </TabsContent>
      <TabsContent value="history">
        <AlarmHistoryTab v-model:page="page" v-model:page-size="pageSize" />
      </TabsContent>
    </Tabs>
  </PageContainer>
</template>
