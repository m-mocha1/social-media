import { createApp } from 'vue'
import { connectWebSocket } from '@/wepSock.js'
import App from './App.vue'
import router from './router'

const app = createApp(App)

app.use(router)
connectWebSocket()
app.mount('#app')
