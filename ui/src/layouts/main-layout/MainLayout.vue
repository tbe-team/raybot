<script setup lang="ts">
import { SidebarInset, SidebarProvider, SidebarTrigger } from '@/components/ui/sidebar'
import { useLocalStorage } from '@vueuse/core'
import { Notification, Notivue } from 'notivue'
import HeaderActions from './HeaderActions.vue'
import AppSidebar from './Sidebar.vue'

const open = useLocalStorage('sidebar', true)
</script>

<template>
  <div class="flex min-h-screen">
    <SidebarProvider v-model:open="open">
      <AppSidebar />
      <SidebarInset>
        <header class="flex items-center h-16 gap-2 px-4 border-b shrink-0">
          <SidebarTrigger />
          <div class="px-4 ml-auto">
            <HeaderActions />
          </div>
        </header>

        <main class="flex flex-col flex-1">
          <RouterView />
        </main>
        <Notivue v-slot="item">
          <Notification :item="item" />
        </Notivue>
      </SidebarInset>
    </SidebarProvider>
  </div>
</template>
