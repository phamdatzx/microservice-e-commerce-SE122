import { createRouter, createWebHistory } from 'vue-router'
import AuthenticationView from '@/views/authentication/AuthenticationView.vue'
import ActivationView from '@/views/authentication/ActivationView.vue'
import CategoryView from '@/views/seller/CategoryView.vue'
import { default as SellerHomeView } from '@/views/seller/HomeView.vue'
import { default as CustomerHomeView } from '@/views/customer/HomeView.vue'
import ProductView from '@/views/seller/ProductView.vue'
import LoginForm from '@/views/authentication/forms/LoginForm.vue'
import RegisterForm from '@/views/authentication/forms/RegisterForm.vue'
import ForgotPasswordForm from '@/views/authentication/forms/ForgotPasswordForm.vue'
import ResetPasswordForm from '@/views/authentication/forms/ResetPasswordForm.vue'
import ProductDetailView from '@/views/customer/ProductDetailView.vue'
import CartView from '@/views/customer/CartView.vue'

const loginFormProps = {
  title: 'Welcome Back',
  subtitle: 'LOGIN TO CONTINUE',
  mainBtnText: 'LOG IN',
  switchText: 'NEW USER ?',
  switchLinkUrl: '/register',
  switchLinkText: 'SIGN UP',
}
const registerFormProps = {
  title: 'Register',
  subtitle: 'JOIN TO US',
  mainBtnText: 'REGISTER',
  switchText: 'ALREADY USER ?',
  switchLinkUrl: '/login',
  switchLinkText: 'LOG IN',
}
const forgotPasswordFormProps = {
  title: 'Forgot Password',
  subtitle: 'RESET YOUR PASSWORD',
  mainBtnText: 'RESET PASSWORD',
  switchText: 'REMEMBERED YOUR PASSWORD ?',
  switchLinkUrl: '/login',
  switchLinkText: 'LOG IN',
}
const resetFormProps = {
  title: 'Reset Password',
  subtitle: 'SET A NEW PASSWORD',
  mainBtnText: 'SET NEW PASSWORD',
  switchText: 'REMEMBERED YOUR PASSWORD ?',
  switchLinkUrl: '/login',
  switchLinkText: 'LOG IN',
}

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      component: CustomerHomeView,
    },
    {
      path: '/product-detail',
      component: ProductDetailView,
    },
    {
      path: '/cart',
      component: CartView,
    },
    {
      path: '/',
      component: AuthenticationView,
      children: [
        {
          path: '/',
          component: LoginForm,
          meta: loginFormProps,
        },
        {
          path: '/login',
          name: 'login',
          component: LoginForm,
          meta: loginFormProps,
        },
        {
          path: '/register',
          name: 'register',
          component: RegisterForm,
          meta: registerFormProps,
        },
        {
          path: '/forgot-password',
          name: 'forgot-password',
          component: ForgotPasswordForm,
          meta: forgotPasswordFormProps,
        },
        {
          path: '/reset-password',
          name: 'reset-password',
          component: ResetPasswordForm,
          meta: resetFormProps,
        },
      ],
    },
    {
      path: '/activate',
      name: 'activate',
      component: ActivationView,
    },
    {
      path: '/seller',
      name: 'seller',
      component: SellerHomeView,
      children: [
        { path: '/category', name: 'category', component: CategoryView },
        { path: 'product', name: 'product', component: ProductView },
      ],
    },
  ],
})

export default router
