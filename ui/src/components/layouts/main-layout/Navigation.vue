<script setup lang="ts">
import type { LucideIcon } from 'lucide-vue-next'
import { SidebarGroup, SidebarMenu, SidebarMenuButton, SidebarMenuItem } from '@/components/ui/sidebar'

interface Item {
  title: string
  path: string
  icon: LucideIcon
  isActive?: boolean
}

interface Props {
  items: Item[]
}
const props = defineProps<Props>()

export type { Item }
</script>

<template>
  <SidebarGroup>
    <SidebarMenu>
      <SidebarMenuItem v-for="item in props.items" :key="item.title">
        <SidebarMenuButton as-child :is-active="item.isActive">
          <RouterLink
            :to="item.path"
            active-class="font-semibold bg-muted"
          >
            <component :is="item.icon" v-if="item.icon" />
            <span>{{ item.title }}</span>
          </RouterLink>
        </SidebarMenuButton>
      </SidebarMenuItem>
    </SidebarMenu>
  </SidebarGroup>
</template>
