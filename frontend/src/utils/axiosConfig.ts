import axios from 'axios'
import router from '@/router'
import { eventBus } from '@/utils/eventBus'

// Add a response interceptor
axios.interceptors.response.use(
  (response) => {
    // Any status code that lie within the range of 2xx cause this function to trigger
    return response
  },
  (error) => {
    // Any status codes that falls outside the range of 2xx cause this function to trigger
    if (error.response && error.response.status === 401) {
      // Check if the request was NOT a login request (to allow login form to handle "User not found" etc.)
      // Usually login is a POST to /user/public/login associated with "USER_API_URL" in LoginForm
      const requestUrl = error.config.url
      if (!requestUrl.includes('/login')) {
        localStorage.removeItem('access_token')
        localStorage.removeItem('user_id')
        eventBus.emit('user_logged_out')
        router.push('/')
      }
    }
    return Promise.reject(error)
  },
)
