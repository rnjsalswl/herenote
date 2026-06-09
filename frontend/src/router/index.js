import { createRouter, createWebHistory } from 'vue-router'
import { isLoggedIn } from '@/stores/user.js'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/login', name: 'login', component: () => import('../views/LoginView.vue'), meta: { chrome: false } },
    { path: '/', name: 'home', component: () => import('../views/HomeView.vue'), meta: { chrome: true } },
    { path: '/explore', name: 'explore', component: () => import('../views/ExploreView.vue'), meta: { chrome: true } },
    { path: '/profile', name: 'profile', component: () => import('../views/ProfileView.vue'), meta: { chrome: true } },
    { path: '/places/:id', name: 'place', component: () => import('../views/PlaceView.vue'), meta: { chrome: false } },
    { path: '/compose/:id', name: 'compose', component: () => import('../views/ComposeView.vue'), meta: { chrome: false } },
    { path: '/add-place', name: 'add-place', component: () => import('../views/AddPlaceView.vue'), meta: { chrome: false } },
  ],
})

router.beforeEach((to) => {
  if (to.name !== 'login' && !isLoggedIn()) return { name: 'login' }
})

export default router
