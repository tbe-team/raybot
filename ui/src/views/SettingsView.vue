<script setup lang="ts">
import CloudConfigTab from '@/components/app/settings/CloudConfigTab.vue'
import HardwareConfigTab from '@/components/app/settings/HardwareConfigTab.vue'
import HTTPConfigTab from '@/components/app/settings/HTTPConfigTab.vue'
import WifiConfigTab from '@/components/app/settings/WifiConfigTab.vue'
import PageContainer from '@/components/shared/PageContainer.vue'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'

const route = useRoute()
const router = useRouter()
const tab = route.query.tab as string | undefined ?? 'hardware'

function handleTabChange(value: string | number) {
  router.replace({ query: { tab: value } })
}
</script>

<template>
  <PageContainer>
    <div>
      <h1 class="text-xl font-semibold tracking-tight">
        Settings
      </h1>
      <p class="text-sm text-muted-foreground">
        Configure your robot's settings
        <RouterLink to="/reboot" class="text-blue-500">
          (Reboot to apply changes)
        </RouterLink>
      </p>
    </div>
    <Tabs :default-value="tab" @update:model-value="handleTabChange">
      <TabsList>
        <TabsTrigger value="hardware">
          Hardware
        </TabsTrigger>
        <TabsTrigger value="cloud">
          Cloud
        </TabsTrigger>
        <TabsTrigger value="http">
          HTTP
        </TabsTrigger>
        <TabsTrigger value="wifi">
          WiFi
        </TabsTrigger>
      </TabsList>

      <TabsContent value="hardware">
        <HardwareConfigTab />
      </TabsContent>
      <TabsContent value="cloud">
        <CloudConfigTab />
      </TabsContent>
      <TabsContent value="http">
        <HTTPConfigTab />
      </TabsContent>
      <TabsContent value="wifi">
        <WifiConfigTab />
      </TabsContent>
    </Tabs>
  </PageContainer>
</template>
