import type { RouteRecordRaw } from 'vue-router'
import { useNProgress } from '@/lib/nprogress'
import { createRouter, createWebHistory } from 'vue-router'
import 'nprogress/nprogress.css'

const MainLayout = () => import('@/layouts/main-layout/MainLayout.vue')

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/system',
  },
  {
    path: '/system',
    component: MainLayout,
    children: [
      {
        path: '',
        component: () => import('@/views/system/ConfigurationView.vue'),
        meta: {
          title: 'System Configuration',
        },
      },
      {
        path: 'restart',
        component: () => import('@/views/system/RestartView.vue'),
        meta: {
          title: 'System Restart',
        },
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
  let title = 'Raybot UI'
  if (to.meta.title) {
    title = `${to.meta.title} | ${title}`
  }
  document.title = title
  nprogress.start()
  next()
})

router.afterEach(() => {
  nprogress.done()
})

export default router
