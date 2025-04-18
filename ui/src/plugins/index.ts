import type { App } from 'vue'
import router from '@/router'
import { VueQueryPlugin } from '@tanstack/vue-query'
import { createPinia } from 'pinia'
import { notivue } from './notivue'
import { queryPluginOpts } from './vue-query'

/**
 * Register plugins
 * @param app - Vue app instance
 * @description This function registers all plugins for the application
 */
export function registerPlugins(app: App) {
  app
    .use(router)
    .use(createPinia())
    .use(notivue)
    .use(VueQueryPlugin, queryPluginOpts)
}
