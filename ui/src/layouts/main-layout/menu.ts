import type { Component } from 'vue'
import {
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
    name: 'System config',
    path: '/system',
    icon: Settings2,
  },
  {
    name: 'Restart',
    path: '/system/restart',
    icon: Power,
  },
]

export { items as routes }
