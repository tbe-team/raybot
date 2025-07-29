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
import HeaderWarning from './HeaderWarning.vue'

const open = useLocalStorage('sidebar', true)
</script>

<template>
  <div>
    <Notivue v-slot="item">
      <Notification :item="item" />
    </Notivue>
    <ConfirmationDialog />
    <SidebarProvider v-model:open="open">
      <AppSidebar />
      <SidebarInset>
        <header
          class="fixed top-0 right-0 z-10 transition-[left,right,width] duration-0 md:!duration-200 ease-linear bg-white dark:bg-background shadow left-0"
          :class="clsx(open ? 'md:left-56' : 'md:left-12')"
        >
          <div class="flex justify-between items-center h-12 border-b shrink-0">
            <div class="flex gap-2 items-center px-4">
              <SidebarTrigger class="-ml-1" />
              <Separator orientation="vertical" class="mr-2 h-4" />
            </div>
            <HeaderActions class="mr-4" />
          </div>
          <HeaderWarning />
        </header>
        <main class="flex flex-col flex-1 mt-12">
          <RouterView />
        </main>
      </SidebarInset>
      <ScrollToTopButton />
    </SidebarProvider>
  </div>
</template>
