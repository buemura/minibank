import client from './client'
import type { AuthResponse, ChangePasswordRequest, LoginCredentials, RegisterData, User } from '@/types'

export const authApi = {
  async login(credentials: LoginCredentials): Promise<AuthResponse> {
    const { data } = await client.post<AuthResponse>('/auth/login', credentials)
    return data
  },

  async register(userData: RegisterData): Promise<AuthResponse> {
    const { data } = await client.post<AuthResponse>('/auth/register', userData)
    return data
  },

  async refresh(refreshToken: string): Promise<AuthResponse> {
    const { data } = await client.post<AuthResponse>('/auth/refresh', { refresh_token: refreshToken })
    return data
  },

  async logout(refreshToken: string): Promise<void> {
    await client.post('/auth/logout', { refresh_token: refreshToken })
  },

  async me(): Promise<User> {
    const { data } = await client.get<User>('/auth/me')
    return data
  },

  async changePassword(data: ChangePasswordRequest): Promise<{ success: boolean; message: string }> {
    const { data: response } = await client.post<{ success: boolean; message: string }>('/auth/change-password', data)
    return response
  },
}
