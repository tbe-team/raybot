import { createApp } from 'vue'
import App from './App.vue'
import { registerPlugins } from './plugins'

import './index.css'

const app = createApp(App)

registerPlugins(app)

app.mount('#app')
