import { createRouter, createWebHistory } from 'vue-router'
import AuthenticationView from '@/views/authentication/AuthenticationView.vue'
import ActivationView from '@/views/authentication/ActivationView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: AuthenticationView,
      props: { formType: 'login' },
    },
    {
      path: '/login',
      name: 'login',
      component: AuthenticationView,
      props: { formType: 'login' },
    },
    {
      path: '/register',
      name: 'register',
      component: AuthenticationView,
      props: { formType: 'register' },
    },
    {
      path: '/forgot-password',
      name: 'forgot-password',
      component: AuthenticationView,
      props: { formType: 'forgot-password' },
    },
    {
      path: '/activate',
      name: 'activate',
      component: ActivationView,
    },
    {
      path: '/reset-password',
      name: 'reset-password',
      component: AuthenticationView,
      props: { formType: 'reset-password' },
    },
  ],
})

export default router
