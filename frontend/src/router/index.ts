import { createRouter, createWebHistory } from 'vue-router'
import AuthenticationView from '@/views/authentication/AuthenticationView.vue'
import ActivationView from '@/views/authentication/ActivationView.vue'
import CategoryView from '@/views/seller/CategoryView.vue'
import { default as SellerHomeView } from '@/views/seller/HomeView.vue'
import { default as CustomerHomeView } from '@/views/customer/HomeView/HomeView.vue'
import ProductView from '@/views/seller/ProductView.vue'
import LoginForm from '@/views/authentication/forms/LoginForm.vue'
import RegisterForm from '@/views/authentication/forms/RegisterForm.vue'
import ForgotPasswordForm from '@/views/authentication/forms/ForgotPasswordForm.vue'
import ResetPasswordForm from '@/views/authentication/forms/ResetPasswordForm.vue'
import ProductDetailView from '@/views/customer/ProductDetailView.vue'
import ProfileView from '@/views/customer/ProfileView/ProfileView.vue'
import CartView from '@/views/customer/CartView/CartView.vue'
import SellerView from '@/views/customer/SellerView/SellerView.vue'
import SearchView from '@/views/customer/SearchView/SearchView.vue'
import CheckoutView from '@/views/customer/CheckoutView/CheckoutView.vue'
import ChatView from '@/views/seller/ChatView.vue'
import SellerProfileView from '@/views/seller/ProfileView.vue'
import VoucherManagerView from '@/views/seller/VoucherManagerView.vue'
import OrderManagerView from '@/views/seller/OrderManagerView.vue'
import StatisticView from '@/views/seller/StatisticView.vue'
import AdminHomeView from '@/views/admin/HomeView.vue'
import CategoryManagerView from '@/views/admin/CategoryView.vue'
import UserView from '@/views/admin/UserView.vue'
import ReportView from '@/views/admin/ReportView.vue'
import OrderTrackingView from '@/views/customer/OrderTrackingView/OrderTrackingView.vue'
import MyOrder from '@/views/customer/ProfileView/components/MyOrder.vue'
import CustomerLayout from '@/views/customer/CustomerLayout.vue'

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
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else if (to.path.startsWith('/profile') && from.path.startsWith('/profile')) {
      return
    } else {
      return { top: 0 }
    }
  },
  routes: [
    {
      path: '/',
      component: CustomerLayout,
      children: [
        {
          path: '',
          component: CustomerHomeView,
        },
        {
          path: 'product/:id?',
          component: ProductDetailView,
        },
        {
          path: 'cart',
          component: CartView,
        },
        {
          path: 'profile',
          component: ProfileView,
          children: [
            {
              path: '',
              redirect: '/profile/account-info',
            },
            {
              path: 'account-info',
              name: 'account-info',
              component: () => import('@/components/profile/AccountInfo.vue'),
            },
            {
              path: 'orders',
              name: 'my-orders',
              component: () => import('@/views/customer/ProfileView/components/MyOrder.vue'),
            },
            {
              path: 'address',
              name: 'my-address',
              component: () => import('@/components/profile/MyAddress.vue'),
            },
            {
              path: 'change-password',
              name: 'change-password',
              component: () => import('@/components/profile/ChangePassword.vue'),
            },
            {
              path: 'vouchers',
              name: 'my-vouchers',
              component: () => import('@/views/customer/ProfileView/components/MyVoucher.vue'),
            },
          ],
        },
        {
          path: 'checkout',
          component: CheckoutView,
        },
        {
          path: 'checkout/success',
          component: () => import('@/views/customer/CheckoutView/PaymentSuccessView.vue'),
        },
        {
          path: 'checkout/failure',
          component: () => import('@/views/customer/CheckoutView/PaymentFailureView.vue'),
        },
        {
          path: 'order-tracking/:id',
          component: OrderTrackingView,
        },
        {
          path: 'recently-viewed',
          name: 'recently-viewed',
          component: () => import('@/views/customer/RecentlyViewedView.vue'),
        },
        {
          path: 'seller-page/:sellerId',
          component: SellerView,
        },
        {
          path: 'search',
          name: 'search',
          component: SearchView,
        },
      ],
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
        { path: 'category', name: 'category', component: CategoryView },
        { path: 'product', name: 'product', component: ProductView },
        { path: 'chat', name: 'chat', component: ChatView },
        { path: 'profile', name: 'seller-profile', component: SellerProfileView },
        { path: 'voucher', name: 'voucher', component: VoucherManagerView },
        { path: 'order', name: 'order', component: OrderManagerView },
        { path: 'statistic', name: 'statistic', component: StatisticView },
      ],
    },
    {
      path: '/admin',
      name: 'admin',
      component: AdminHomeView,
      children: [
        { path: 'category', name: 'admin-category', component: CategoryManagerView },
        { path: 'users', name: 'admin-users', component: UserView },
        { path: 'report', name: 'admin-report', component: ReportView },
      ],
    },
  ],
})

export default router
