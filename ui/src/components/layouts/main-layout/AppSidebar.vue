<script setup lang="ts">
import type { SidebarProps } from '@/components/ui/sidebar'
import VersionDisplay from '@/components/shared/VersionDisplay.vue'
import { Sidebar, SidebarContent, SidebarHeader, SidebarMenu, SidebarMenuButton, SidebarMenuItem } from '@/components/ui/sidebar'
import { Bot } from 'lucide-vue-next'
import { menus } from './menu'
import Navigation from './Navigation.vue'

const props = withDefaults(defineProps<SidebarProps>(), {
  collapsible: 'icon',
})
</script>

<template>
  <Sidebar v-bind="props" class="bg-popover">
    <SidebarHeader class="py-4">
      <SidebarMenu>
        <SidebarMenuItem>
          <SidebarMenuButton size="lg" as-child>
            <div class="flex items-center gap-3 px-6 py-4">
              <div class="flex items-center justify-center bg-blue-600 rounded-lg aspect-square size-8 text-sidebar-primary-foreground">
                <Bot class="size-4" />
              </div>
              <div class="flex flex-col gap-0.5">
                <span class="text-sm font-semibold leading-none">Raybot UI</span>
                <VersionDisplay />
              </div>
            </div>
          </SidebarMenuButton>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarHeader>
    <SidebarContent>
      <Navigation :items="menus" />
    </SidebarContent>
  </Sidebar>
</template>
