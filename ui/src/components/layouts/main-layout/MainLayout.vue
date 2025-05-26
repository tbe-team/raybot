<script setup lang="ts">
import { useLocalStorage } from '@vueuse/core'
import clsx from 'clsx'
import { Notification, Notivue } from 'notivue'
import ConfirmationDialog from '@/components/shared/ConfirmationDialog.vue'
import ScrollToTopButton from '@/components/shared/ScrollToTopButton.vue'
import { Separator } from '@/components/ui/separator'
import { SidebarInset, SidebarProvider, SidebarTrigger } from '@/components/ui/sidebar'
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
      <header
        class="fixed top-0 right-0 z-10 transition-[left,right,width] !duration-200 ease-linear bg-white dark:bg-background shadow"
        :class="clsx(open ? 'lg:left-56 left-0' : 'lg:left-12 left-0')"
      >
        <div class="flex items-center justify-between h-12 border-b shrink-0">
          <div class="flex items-center gap-2 px-4">
            <SidebarTrigger class="-ml-1" />
            <Separator orientation="vertical" class="h-4 mr-2" />
          </div>
          <HeaderActions class="mr-4" />
        </div>
      </header>
      <main class="flex flex-col flex-1 mt-12">
        <RouterView />
      </main>
    </SidebarInset>
    <ScrollToTopButton />
  </SidebarProvider>
</template>
