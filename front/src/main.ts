import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import keycloak from './keycloak'
import { connectWebSocket } from './services/websocketService'
import axios from 'axios'

keycloak.init({ onLoad: 'login-required' }).then((authenticated) => {
    console.log('authenticated', authenticated)
    if (authenticated) {
      axios.interceptors.request.use(
        (config) => {
          const localToken = keycloak.token;
          config.headers.Authorization = `Bearer ${localToken}`;
          return config;
        },

        (error) => Promise.reject(error),
      );
        connectWebSocket();

        const app = createApp(App)

        app.use(createPinia())
        app.use(router)
        app.mount('#app');
    } else {
      window.location.reload();
    }
  }).catch((error) => {
    console.error('Keycloak initialization failed:', error);
    window.location.reload();
  });