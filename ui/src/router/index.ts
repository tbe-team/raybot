import type { RouteRecordRaw } from 'vue-router'
import { useNProgress } from '@/lib/nprogress'
import { createRouter, createWebHistory } from 'vue-router'

const MainLayout = () => import('@/components/layouts/main-layout/MainLayout.vue')

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: { name: 'settings' },
    component: MainLayout,
    children: [
      {
        path: 'state',
        name: 'state',
        component: () => import('@/views/StateView.vue'),
        meta: { title: 'State' },
      },
      {
        path: 'commands',
        name: 'commands',
        component: () => import('@/views/command/CommandsView.vue'),
        meta: { title: 'Commands' },
      },
      {
        path: 'settings',
        name: 'settings',
        component: () => import('@/views/SettingsView.vue'),
        meta: { title: 'Settings' },
      },
      {
        path: 'restart',
        name: 'restart',
        component: () => import('@/views/RestartView.vue'),
        meta: { title: 'System Restart' },
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
