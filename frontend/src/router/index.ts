import { createRouter, createWebHistory } from 'vue-router'
import Dashboard from '@/pages/Dashboard.vue'
import Mantenimientos from '@/pages/MaintenanceExplorer.vue'
import Maquinas from '@/pages/MaquinasChecklist.vue'
import Tecnicos from '@/pages/Tecnicos.vue'
import Login from '@/pages/Login.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/', component: Dashboard, meta: { requiresAdmin: true } },
    { path: '/mantenimientos', component: Mantenimientos, meta: { requiresAdmin: true } },
    { path: '/maquinas', component: Maquinas, meta: { requiresAdmin: true } },
    { path: '/tecnicos', component: Tecnicos, meta: { requiresAdmin: true } },
    { path: '/login', component: Login },
  ],
})

router.beforeEach((to) => {
  const hasToken = Boolean(localStorage.getItem('access-admin-token'))
  if (to.meta.requiresAdmin && !hasToken) return '/login'
  if (to.path === '/login' && hasToken) return '/'
  return true
})

export default router
