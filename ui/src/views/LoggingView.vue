<script setup lang="ts">
import LogConsoleConfigTab from '@/components/app/logging/LogConsoleConfigTab.vue'
import LogFileConfigTab from '@/components/app/logging/LogFileConfigTab.vue'
import PageContainer from '@/components/shared/PageContainer.vue'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'

const route = useRoute()
const router = useRouter()
const tab = route.query.tab as string | undefined ?? 'console'

function handleTabChange(value: string | number) {
  router.replace({ query: { tab: value } })
}
</script>

<template>
  <PageContainer>
    <div>
      <h1 class="text-xl font-semibold tracking-tight">
        Logging
      </h1>
      <p class="text-sm text-muted-foreground">
        Log configuration allows users to set up and manage the logging system
        <RouterLink to="/reboot" class="text-blue-500">
          (Reboot to apply changes)
        </RouterLink>
      </p>
    </div>
    <Tabs :default-value="tab" @update:model-value="handleTabChange">
      <TabsList>
        <TabsTrigger value="console">
          Console Log
        </TabsTrigger>
        <TabsTrigger value="file">
          File Log
        </TabsTrigger>
      </TabsList>

      <TabsContent value="console">
        <LogConsoleConfigTab />
      </TabsContent>
      <TabsContent value="file">
        <LogFileConfigTab />
      </TabsContent>
    </Tabs>
  </PageContainer>
</template>
