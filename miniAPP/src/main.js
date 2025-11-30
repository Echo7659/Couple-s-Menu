import { createSSRApp } from 'vue'
import App from './App.vue'
import store from '@/store'
import * as Pinia from 'pinia'
import { createUnistorage } from 'pinia-plugin-unistorage'
// #ifdef MP-WEIXIN
import share from '@/utils/share'
// #endif
// 引入UnoCSS
import 'virtual:uno.css'
export function createApp() {
  const app = createSSRApp(App)
  app.mixin(share)
  store.use(createUnistorage())
  app.use(store)
  // #ifdef MP-WEIXIN
  app.mixin(share)
  // #endif
  return {
    app,
    Pinia
  }
}
