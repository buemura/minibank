import { Link } from '@tanstack/react-router'
import { usePrimaryAccount } from '@/hooks/useAccounts'
import Navbar from '@/components/layout/Navbar'
import TransferForm from '@/components/transaction/TransferForm'
import Card from '@/components/common/Card'

export default function TransferPage() {
  const { primaryAccount, isLoading } = usePrimaryAccount()

  return (
    <div className="min-h-screen">
      <Navbar />

      <main className="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
        <div className="mb-8">
          <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">Transfer Money</h1>
          <p className="text-gray-600 dark:text-gray-400">Send money to another account instantly.</p>
        </div>

        <div className="max-w-lg">
          <Card title="New Transfer">
            {isLoading ? (
              <div className="flex justify-center py-8">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600" />
              </div>
            ) : !primaryAccount ? (
              <div className="text-center py-8 text-gray-500 dark:text-gray-400">
                You need an account to make transfers.{' '}
                <Link to="/dashboard" className="text-primary-600 hover:text-primary-700">
                  Go to Dashboard
                </Link>
              </div>
            ) : (
              <TransferForm />
            )}
          </Card>
        </div>
      </main>
    </div>
  )
}
