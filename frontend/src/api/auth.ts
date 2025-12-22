import api from './axios'

export interface LoginParams {
  username: string
  password: string
}

export interface RegisterParams {
  username: string
  email: string
  password: string
  phone?: string
}

export interface User {
  id: number
  username: string
  email: string
  phone?: string
  role: string
  avatar?: string
}

export interface LoginResponse {
  token: string
  expires_at: number
  user: User
}

export const authAPI = {
  // 用户登录
  login: (params: LoginParams) => {
    return api.post<any, LoginResponse>('/auth/login', params)
  },

  // 用户注册
  register: (params: RegisterParams) => {
    return api.post<any, LoginResponse>('/auth/register', params)
  },

  // 获取当前用户信息
  getUserInfo: () => {
    return api.get<any, User>('/auth/me')
  },
}
