import type { RouteRecordRaw } from 'vue-router'
import { useNProgress } from '@/lib/nprogress'
import { createRouter, createWebHistory } from 'vue-router'

const MainLayout = () => import('@/components/layouts/main-layout/MainLayout.vue')

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: { name: 'state' },
    component: MainLayout,
    children: [
      {
        path: 'state',
        name: 'state',
        component: () => import('@/views/StateView.vue'),
        meta: { title: 'State' },
      },
      {
        path: 'command-queue',
        name: 'command-queue',
        component: () => import('@/views/CommandQueueView.vue'),
        meta: { title: 'Command queue' },
      },
      {
        path: 'command-history',
        name: 'command-history',
        component: () => import('@/views/CommandHistoryView.vue'),
        meta: { title: 'Command history' },
      },
      {
        path: 'settings',
        name: 'settings',
        component: () => import('@/views/SettingsView.vue'),
        meta: { title: 'Settings' },
      },
      {
        path: 'reboot',
        name: 'reboot',
        component: () => import('@/views/RebootView.vue'),
        meta: { title: 'Reboot' },
      },
    ],
  },
  {
    path: '/404',
    component: () => import('@/views/NotFoundView.vue'),
  },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

const nprogress = useNProgress()

router.beforeEach((to, _, next) => {
  nprogress.start()
  document.title = to.meta.title ? `${to.meta.title} | Raybot UI` : 'Raybot UI'
  next()
})

router.afterEach(() => {
  nprogress.done()
})

export default router
