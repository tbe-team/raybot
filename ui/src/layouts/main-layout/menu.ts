import type { Component } from 'vue'
import {
  ChartBarBig,
  Power,
  Settings2,
} from 'lucide-vue-next'

interface Item {
  name: string
  path: string
  icon: Component
}

const items: Item[] = [
  {
    name: 'State',
    path: '/state',
    icon: ChartBarBig,
  },
  {
    name: 'System config',
    path: '/system/configuration',
    icon: Settings2,
  },
  {
    name: 'Restart',
    path: '/system/restart',
    icon: Power,
  },
]

export { items as routes }
