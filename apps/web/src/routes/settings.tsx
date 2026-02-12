import { useState } from 'react'
import { authApi } from '@/api/auth'
import { useAuth } from '@/contexts/AuthContext'
import Navbar from '@/components/layout/Navbar'
import Card from '@/components/common/Card'
import Input from '@/components/common/Input'
import Button from '@/components/common/Button'

export default function SettingsPage() {
  const { logout } = useAuth()

  const [currentPassword, setCurrentPassword] = useState('')
  const [newPassword, setNewPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState<string | null>(null)

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError(null)
    setSuccess(null)

    if (newPassword !== confirmPassword) {
      setError('New passwords do not match')
      return
    }

    if (newPassword.length < 8) {
      setError('New password must be at least 8 characters')
      return
    }

    setIsLoading(true)

    try {
      await authApi.changePassword({
        current_password: currentPassword,
        new_password: newPassword,
      })

      setSuccess('Password changed successfully. You will be logged out shortly.')
      setCurrentPassword('')
      setNewPassword('')
      setConfirmPassword('')

      setTimeout(() => {
        logout()
        window.location.href = '/login'
      }, 2000)
    } catch (err: unknown) {
      const apiErr = err as { response?: { data?: { error?: string } } }
      setError(apiErr.response?.data?.error || 'Failed to change password')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="min-h-screen">
      <Navbar />

      <main className="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
        <div className="mb-8">
          <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">Settings</h1>
          <p className="text-gray-600 dark:text-gray-400">Manage your account settings</p>
        </div>

        <div className="max-w-md">
          <Card title="Change Password">
            <form onSubmit={handleSubmit} className="space-y-6">
              <Input
                value={currentPassword}
                onChange={setCurrentPassword}
                label="Current Password"
                type="password"
                placeholder="Enter your current password"
                required
              />
              <Input
                value={newPassword}
                onChange={setNewPassword}
                label="New Password"
                type="password"
                placeholder="Enter your new password (min 8 characters)"
                required
              />
              <Input
                value={confirmPassword}
                onChange={setConfirmPassword}
                label="Confirm New Password"
                type="password"
                placeholder="Confirm your new password"
                required
              />

              {error && (
                <div className="bg-red-50 dark:bg-red-900/30 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 px-4 py-3 rounded-lg text-sm">
                  {error}
                </div>
              )}

              {success && (
                <div className="bg-green-50 dark:bg-green-900/30 border border-green-200 dark:border-green-800 text-green-700 dark:text-green-400 px-4 py-3 rounded-lg text-sm">
                  {success}
                </div>
              )}

              <Button type="submit" loading={isLoading} className="w-full">
                Change Password
              </Button>
            </form>
          </Card>
        </div>
      </main>
    </div>
  )
}
