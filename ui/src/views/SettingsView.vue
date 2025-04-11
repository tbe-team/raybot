<script setup lang="ts">
import CargoConfigTab from '@/components/app/settings/CargoConfigTab.vue'
import CloudConfigTab from '@/components/app/settings/CloudConfigTab.vue'
import GRPCConfigTab from '@/components/app/settings/GRPCConfigTab.vue'
import HardwareConfigTab from '@/components/app/settings/HardwareConfigTab.vue'
import HTTPConfigTab from '@/components/app/settings/HTTPConfigTab.vue'
import LogConfigTab from '@/components/app/settings/LogConfigTab.vue'
import PageContainer from '@/components/shared/PageContainer.vue'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'

const route = useRoute()
const router = useRouter()
const tab = route.query.tab as string | undefined ?? 'log'

function handleTabChange(value: string | number) {
  router.replace({ query: { tab: value } })
}
</script>

<template>
  <PageContainer>
    <Tabs :default-value="tab" @update:model-value="handleTabChange">
      <TabsList>
        <TabsTrigger value="log">
          Log
        </TabsTrigger>
        <TabsTrigger value="hardware">
          Hardware
        </TabsTrigger>
        <TabsTrigger value="cloud">
          Cloud
        </TabsTrigger>
        <TabsTrigger value="grpc">
          GRPC
        </TabsTrigger>
        <TabsTrigger value="http">
          HTTP
        </TabsTrigger>
        <TabsTrigger value="cargo">
          Cargo
        </TabsTrigger>
      </TabsList>

      <TabsContent value="log">
        <LogConfigTab />
      </TabsContent>
      <TabsContent value="hardware">
        <HardwareConfigTab />
      </TabsContent>
      <TabsContent value="cloud">
        <CloudConfigTab />
      </TabsContent>
      <TabsContent value="grpc">
        <GRPCConfigTab />
      </TabsContent>
      <TabsContent value="http">
        <HTTPConfigTab />
      </TabsContent>
      <TabsContent value="cargo">
        <CargoConfigTab />
      </TabsContent>
    </Tabs>
  </PageContainer>
</template>
