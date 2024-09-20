import { createRouter, createWebHistory } from 'vue-router'
import BoardListView from "../views/BoardListView.vue";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: BoardListView
    },
    {
      path: '/board/:id',
      name: 'board',
      component: () => import('../views/BoardView.vue')
    }
  ]
})

export default router
