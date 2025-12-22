import axios from 'axios'
import { message } from 'antd'

const api = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
})

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  (response) => {
    const { data } = response
    if (data.code === 0) {
      return data.data
    } else {
      message.error(data.message || '请求失败')
      return Promise.reject(data)
    }
  },
  (error) => {
    if (error.response) {
      const { status } = error.response
      if (status === 401) {
        message.error('请先登录')
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        window.location.href = '/login'
      } else if (status === 403) {
        message.error('没有权限')
      } else if (status === 404) {
        message.error('请求的资源不存在')
      } else if (status >= 500) {
        message.error('服务器错误')
      } else {
        message.error(error.response.data?.message || '请求失败')
      }
    } else if (error.request) {
      message.error('网络错误，请检查网络连接')
    } else {
      message.error(error.message || '请求失败')
    }
    return Promise.reject(error)
  }
)

export default api
