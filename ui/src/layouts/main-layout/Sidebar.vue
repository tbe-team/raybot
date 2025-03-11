<script setup lang="ts">
import {
  Sidebar,
  SidebarContent,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarRail,
} from '@/components/ui/sidebar'

import Logo from './Logo.vue'
import { routes } from './menu'
</script>

<template>
  <Sidebar collapsible="icon" class="bg-popover">
    <SidebarHeader class="flex-row items-center mt-3">
      <Logo />
    </SidebarHeader>
    <SidebarContent class="mx-2 mt-6">
      <SidebarMenu>
        <SidebarMenuItem v-for="route in routes" :key="route.name">
          <SidebarMenuButton as-child>
            <RouterLink
              :to="route.path"
              exact-active-class="text-primary bg-muted"
              class="text-xs font-medium transition-colors rounded-lg text-muted-foreground hover:bg-accent hover:text-accent-foreground"
            >
              <component :is="route.icon" class="w-4 h-4" />
              {{ route.name }}
            </RouterLink>
          </SidebarMenuButton>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarContent>
    <SidebarRail />
  </Sidebar>
</template>
