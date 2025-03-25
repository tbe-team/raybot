import type { Component } from 'vue'
import {
  ChartBarBig,
  Command,
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
    name: 'Commands',
    path: '/commands',
    icon: Command,
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
