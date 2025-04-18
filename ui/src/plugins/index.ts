import type { App } from 'vue'
import router from '@/router'
import { VueQueryPlugin } from '@tanstack/vue-query'
import { notivue } from './notivue'
import { pinia } from './pinia'
import { queryPluginOpts } from './vue-query'

/**
 * Register plugins
 * @param app - Vue app instance
 * @description This function registers all plugins for the application
 */
export function registerPlugins(app: App) {
  app
    .use(router)
    .use(pinia)
    .use(notivue)
    .use(VueQueryPlugin, queryPluginOpts)
}
