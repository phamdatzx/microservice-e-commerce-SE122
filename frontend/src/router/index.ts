import { createRouter, createWebHistory } from 'vue-router'
import AuthenticationView from '@/views/authentication/AuthenticationView.vue'
import ActivationView from '@/views/authentication/ActivationView.vue'
import CategoryView from '@/views/seller/CategoryView.vue'
import { default as SellerHomeView } from '@/views/seller/HomeView.vue'
import ProductView from '@/views/seller/ProductView.vue'
import HomeView from '@/views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
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
    {
      path: '/seller',
      name: 'seller',
      component: SellerHomeView,
      children: [
        { path: 'category', name: 'category', component: CategoryView },
        { path: 'product', name: 'product', component: ProductView },
      ],
    },
  ],
})

export default router
