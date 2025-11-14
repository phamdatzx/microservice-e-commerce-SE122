import { createRouter, createWebHistory } from 'vue-router'
import AuthenticationView from '@/views/authentication/AuthenticationView.vue'
import LoginForm from '../views/authentication/components/LoginForm.vue'
import RegisterForm from '@/views/authentication/components/RegisterForm.vue'
import ForgotPasswordForm from '@/views/authentication/components/ForgotPasswordForm.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: AuthenticationView,
      props: { FormComponent: LoginForm, formType: 'login' },
    },
    {
      path: '/login',
      name: 'login',
      component: AuthenticationView,
      props: { FormComponent: LoginForm, formType: 'login' },
    },
    {
      path: '/register',
      name: 'register',
      component: AuthenticationView,
      props: { FormComponent: RegisterForm, formType: 'register' },
    },
    {
      path: '/forgot-password',
      name: 'forgot-password',
      component: AuthenticationView,
      props: { FormComponent: ForgotPasswordForm, formType: 'forgot-password' },
    },
  ],
})

export default router
