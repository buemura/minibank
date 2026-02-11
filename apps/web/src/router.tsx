import { createRouter, createRoute, createRootRoute, redirect, Outlet } from '@tanstack/react-router'
import LoginPage from '@/routes/login'
import RegisterPage from '@/routes/register'
import DashboardPage from '@/routes/dashboard'
import TransferPage from '@/routes/transfer'
import DepositPage from '@/routes/deposit'
import StatementPage from '@/routes/statement'

const rootRoute = createRootRoute({
  component: () => (
    <div className="min-h-screen">
      <Outlet />
    </div>
  ),
})

const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/',
  beforeLoad: () => {
    throw redirect({ to: '/dashboard' })
  },
})

const loginRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/login',
  beforeLoad: () => {
    const token = localStorage.getItem('access_token')
    if (token) {
      throw redirect({ to: '/dashboard' })
    }
  },
  component: LoginPage,
})

const registerRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/register',
  beforeLoad: () => {
    const token = localStorage.getItem('access_token')
    if (token) {
      throw redirect({ to: '/dashboard' })
    }
  },
  component: RegisterPage,
})

const dashboardRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/dashboard',
  beforeLoad: () => {
    const token = localStorage.getItem('access_token')
    if (!token) {
      throw redirect({ to: '/login' })
    }
  },
  component: DashboardPage,
})

const transferRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/transfer',
  beforeLoad: () => {
    const token = localStorage.getItem('access_token')
    if (!token) {
      throw redirect({ to: '/login' })
    }
  },
  component: TransferPage,
})

const depositRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/deposit',
  beforeLoad: () => {
    const token = localStorage.getItem('access_token')
    if (!token) {
      throw redirect({ to: '/login' })
    }
  },
  component: DepositPage,
})

const statementRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/statement',
  beforeLoad: () => {
    const token = localStorage.getItem('access_token')
    if (!token) {
      throw redirect({ to: '/login' })
    }
  },
  component: StatementPage,
})

const routeTree = rootRoute.addChildren([
  indexRoute,
  loginRoute,
  registerRoute,
  dashboardRoute,
  transferRoute,
  depositRoute,
  statementRoute,
])

export const router = createRouter({ routeTree })

declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router
  }
}
