import { createContext, useContext, useState, useCallback, useEffect, type ReactNode } from 'react'
import { authApi } from '@/api/auth'
import type { User, LoginCredentials, RegisterData } from '@/types'

interface AuthContextType {
  user: User | null
  accessToken: string | null
  isAuthenticated: boolean
  isLoading: boolean
  error: string | null
  login: (credentials: LoginCredentials) => Promise<boolean>
  register: (data: RegisterData) => Promise<boolean>
  logout: () => void
  refreshAccessToken: () => Promise<boolean>
  fetchCurrentUser: () => Promise<void>
  clearError: () => void
}

const AuthContext = createContext<AuthContextType | null>(null)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [accessToken, setAccessToken] = useState<string | null>(localStorage.getItem('access_token'))
  const [refreshToken, setRefreshToken] = useState<string | null>(localStorage.getItem('refresh_token'))
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const isAuthenticated = !!accessToken

  const clearError = useCallback(() => setError(null), [])

  const logout = useCallback(() => {
    const rt = localStorage.getItem('refresh_token')
    if (rt) {
      authApi.logout(rt).catch(() => {})
    }

    setUser(null)
    setAccessToken(null)
    setRefreshToken(null)

    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
  }, [])

  const refreshAccessToken = useCallback(async (): Promise<boolean> => {
    const rt = refreshToken ?? localStorage.getItem('refresh_token')
    if (!rt) return false

    try {
      const response = await authApi.refresh(rt)

      setAccessToken(response.access_token)
      setRefreshToken(response.refresh_token)

      localStorage.setItem('access_token', response.access_token)
      localStorage.setItem('refresh_token', response.refresh_token)

      return true
    } catch {
      logout()
      return false
    }
  }, [refreshToken, logout])

  const login = useCallback(async (credentials: LoginCredentials): Promise<boolean> => {
    setIsLoading(true)
    setError(null)

    try {
      const response = await authApi.login(credentials)

      setAccessToken(response.access_token)
      setRefreshToken(response.refresh_token)
      setUser(response.user)

      localStorage.setItem('access_token', response.access_token)
      localStorage.setItem('refresh_token', response.refresh_token)

      return true
    } catch (e: unknown) {
      const err = e as { response?: { data?: { error?: string | { message?: string } } } }
      const apiError = err.response?.data?.error
      setError(typeof apiError === 'object' ? apiError?.message || 'Login failed' : apiError || 'Login failed')
      return false
    } finally {
      setIsLoading(false)
    }
  }, [])

  const register = useCallback(async (data: RegisterData): Promise<boolean> => {
    setIsLoading(true)
    setError(null)

    try {
      const response = await authApi.register(data)

      setAccessToken(response.access_token)
      setRefreshToken(response.refresh_token)
      setUser(response.user)

      localStorage.setItem('access_token', response.access_token)
      localStorage.setItem('refresh_token', response.refresh_token)

      return true
    } catch (e: unknown) {
      const err = e as { response?: { data?: { error?: string | { message?: string } } } }
      const apiError = err.response?.data?.error
      setError(typeof apiError === 'object' ? apiError?.message || 'Registration failed' : apiError || 'Registration failed')
      return false
    } finally {
      setIsLoading(false)
    }
  }, [])

  const fetchCurrentUser = useCallback(async () => {
    const token = localStorage.getItem('access_token')
    if (!token) return

    try {
      const userData = await authApi.me()
      setUser(userData)
    } catch {
      await refreshAccessToken()
    }
  }, [refreshAccessToken])

  useEffect(() => {
    if (accessToken) {
      fetchCurrentUser()
    }
  }, []) // eslint-disable-line react-hooks/exhaustive-deps

  return (
    <AuthContext.Provider
      value={{
        user,
        accessToken,
        isAuthenticated,
        isLoading,
        error,
        login,
        register,
        logout,
        refreshAccessToken,
        fetchCurrentUser,
        clearError,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth(): AuthContextType {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
