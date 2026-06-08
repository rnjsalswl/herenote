import { createRouter, createWebHistory } from 'vue-router'
import { isLoggedIn } from '@/stores/user.js'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
        path: '/login',
        name: 'login',
        component: () => import('../views/LoginView.vue'),
      },
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/places/:id',
      name: 'place',
      component: () => import('../views/PlaceView.vue'),
    },
    {
      path: '/add-place',
      name: 'add-place',
      component: () => import('../views/AddPlaceView.vue'),
    },
  ],
})

// 로그인 안 했으면 로그인 페이지로
router.beforeEach((to) => {
  if (to.name !== 'login' && !isLoggedIn()) {
    return { name: 'login' }
  }
})

export default router