<script setup lang="ts">
import ConfirmationDialog from '@/components/shared/ConfirmationDialog.vue'
import { Separator } from '@/components/ui/separator'
import { SidebarInset, SidebarProvider, SidebarTrigger } from '@/components/ui/sidebar'
import { useLocalStorage } from '@vueuse/core'
import { Notification, Notivue } from 'notivue'
import AppSidebar from './AppSidebar.vue'
import HeaderActions from './HeaderActions.vue'

const open = useLocalStorage('sidebar', true)
</script>

<template>
  <Notivue v-slot="item">
    <Notification :item="item" />
  </Notivue>

  <ConfirmationDialog />

  <SidebarProvider v-model:open="open">
    <AppSidebar />
    <SidebarInset>
      <header class="flex items-center justify-between h-12 border-b shrink-0">
        <div class="flex items-center gap-2 px-4">
          <SidebarTrigger class="-ml-1" />
          <Separator orientation="vertical" class="h-4 mr-2" />
        </div>
        <HeaderActions class="mr-4" />
      </header>
      <main class="flex flex-col flex-1">
        <RouterView />
      </main>
    </SidebarInset>
  </SidebarProvider>
</template>
