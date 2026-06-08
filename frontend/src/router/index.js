import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
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

export default router