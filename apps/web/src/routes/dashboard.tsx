import { Link } from '@tanstack/react-router'
import { usePrimaryAccount, useCreateAccount } from '@/hooks/useAccounts'
import { useTransactionHistory } from '@/hooks/useTransactions'
import Navbar from '@/components/layout/Navbar'
import BalanceDisplay from '@/components/account/BalanceDisplay'
import TransactionList from '@/components/transaction/TransactionList'
import Card from '@/components/common/Card'
import Button from '@/components/common/Button'

export default function DashboardPage() {
  const { primaryAccount, isLoading: accountsLoading } = usePrimaryAccount()
  const createAccountMutation = useCreateAccount()
  const { data: transactionsData } = useTransactionHistory(primaryAccount?.id)

  const transactions = transactionsData?.transactions ?? []

  async function handleCreateAccount() {
    await createAccountMutation.mutateAsync('checking')
  }

  return (
    <div className="min-h-screen">
      <Navbar />

      <main className="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
        <div className="mb-8">
          <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">Dashboard</h1>
          <p className="text-gray-600 dark:text-gray-400">Welcome back! Here's your account overview.</p>
        </div>

        {accountsLoading ? (
          <div className="flex justify-center py-12">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600" />
          </div>
        ) : !primaryAccount ? (
          <div className="text-center py-12">
            <Card>
              <div className="py-8">
                <h3 className="text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">No Account Yet</h3>
                <p className="text-gray-600 dark:text-gray-400 mb-6">Create your first account to start banking.</p>
                <Button onClick={handleCreateAccount} loading={createAccountMutation.isPending}>
                  Create Account
                </Button>
              </div>
            </Card>
          </div>
        ) : (
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            <div className="lg:col-span-1">
              <BalanceDisplay account={primaryAccount} />

              <div className="mt-6 grid grid-cols-3 gap-4">
                <Link to="/transfer">
                  <Button variant="primary" className="w-full">
                    Transfer
                  </Button>
                </Link>
                <Link to="/deposit">
                  <Button variant="primary" className="w-full">
                    Deposit
                  </Button>
                </Link>
                <Link to="/statement">
                  <Button variant="secondary" className="w-full">
                    Statement
                  </Button>
                </Link>
              </div>
            </div>

            <div className="lg:col-span-2">
              <Card title="Recent Transactions">
                <TransactionList
                  transactions={transactions.slice(0, 5)}
                  currentAccountId={primaryAccount.id}
                />

                {transactions.length > 5 && (
                  <div className="mt-4 text-center">
                    <Link to="/statement" className="text-primary-600 hover:text-primary-700 text-sm font-medium">
                      View all transactions
                    </Link>
                  </div>
                )}
              </Card>
            </div>
          </div>
        )}
      </main>
    </div>
  )
}
