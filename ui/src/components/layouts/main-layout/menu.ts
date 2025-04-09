import type { Item } from './Navigation.vue'
import {
  ChartBarBig,
  Command,
  Power,
  Settings2,
} from 'lucide-vue-next'

export const menus: Item[] = [
  {
    title: 'State',
    path: '/state',
    icon: ChartBarBig,
  },
  // {
  //   title: 'Commands',
  //   path: '/commands',
  //   icon: Command,
  // },
  {
    title: 'Settings',
    path: '/settings',
    icon: Settings2,
  },
  {
    title: 'Restart',
    path: '/restart',
    icon: Power,
  },
]

export { menus as routes }
