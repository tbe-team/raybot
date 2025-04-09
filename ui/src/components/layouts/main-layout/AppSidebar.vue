<script setup lang="ts">
import type { SidebarProps } from '@/components/ui/sidebar'
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
    <SidebarHeader>
      <SidebarMenu>
        <SidebarMenuItem>
          <SidebarMenuButton
            size="lg"
            class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground hover:bg-transparent"
          >
            <div class="flex items-center justify-center bg-blue-600 rounded-lg aspect-square size-8 text-sidebar-primary-foreground">
              <Bot class="size-4" />
            </div>
            <span class="font-semibold truncate">Raybot UI</span>
          </SidebarMenuButton>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarHeader>
    <SidebarContent>
      <Navigation :items="menus" />
    </SidebarContent>
  </Sidebar>
</template>
