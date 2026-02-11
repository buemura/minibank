import { useState } from 'react'
import { Link } from '@tanstack/react-router'
import { usePrimaryAccount } from '@/hooks/useAccounts'
import { useStatement } from '@/hooks/useTransactions'
import Navbar from '@/components/layout/Navbar'
import TransactionList from '@/components/transaction/TransactionList'
import Card from '@/components/common/Card'
import Button from '@/components/common/Button'
import Input from '@/components/common/Input'

function getDefaultStartDate(): string {
  const date = new Date()
  date.setMonth(date.getMonth() - 1)
  return date.toISOString().split('T')[0]
}

function getDefaultEndDate(): string {
  return new Date().toISOString().split('T')[0]
}

function formatCurrency(value: string, currency: string = 'BRL'): string {
  const num = parseFloat(value)
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: currency,
  }).format(num)
}

export default function StatementPage() {
  const { primaryAccount } = usePrimaryAccount()

  const [startDate, setStartDate] = useState(getDefaultStartDate())
  const [endDate, setEndDate] = useState(getDefaultEndDate())
  const [appliedStart, setAppliedStart] = useState(getDefaultStartDate())
  const [appliedEnd, setAppliedEnd] = useState(getDefaultEndDate())

  const { data: statement, isLoading } = useStatement(primaryAccount?.id, appliedStart, appliedEnd)

  function handleFilter(e: React.FormEvent) {
    e.preventDefault()
    setAppliedStart(startDate)
    setAppliedEnd(endDate)
  }

  if (!primaryAccount) {
    return (
      <div className="min-h-screen">
        <Navbar />
        <main className="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
          <div className="text-center py-12">
            <Card>
              <div className="py-8 text-gray-500 dark:text-gray-400">
                You need an account to view statements.{' '}
                <Link to="/dashboard" className="text-primary-600 hover:text-primary-700">
                  Go to Dashboard
                </Link>
              </div>
            </Card>
          </div>
        </main>
      </div>
    )
  }

  return (
    <div className="min-h-screen">
      <Navbar />

      <main className="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
        <div className="mb-8">
          <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">Account Statement</h1>
          <p className="text-gray-600 dark:text-gray-400">View your transaction history.</p>
        </div>

        <div className="space-y-6">
          <Card title="Filter">
            <form onSubmit={handleFilter} className="flex flex-wrap gap-4 items-end">
              <div className="flex-1 min-w-[200px]">
                <Input
                  value={startDate}
                  onChange={setStartDate}
                  label="Start Date"
                  type="date"
                />
              </div>
              <div className="flex-1 min-w-[200px]">
                <Input
                  value={endDate}
                  onChange={setEndDate}
                  label="End Date"
                  type="date"
                />
              </div>
              <Button type="submit" loading={isLoading}>
                Apply
              </Button>
            </form>
          </Card>

          {statement && (
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <Card>
                <div className="text-center">
                  <p className="text-sm text-gray-600 dark:text-gray-400">Opening Balance</p>
                  <p className="text-xl font-bold text-gray-900 dark:text-gray-100">
                    {formatCurrency(statement.opening_balance)}
                  </p>
                </div>
              </Card>
              <Card>
                <div className="text-center">
                  <p className="text-sm text-gray-600 dark:text-gray-400">Closing Balance</p>
                  <p className="text-xl font-bold text-primary-600">
                    {formatCurrency(statement.closing_balance)}
                  </p>
                </div>
              </Card>
              <Card>
                <div className="text-center">
                  <p className="text-sm text-gray-600 dark:text-gray-400">Total Transactions</p>
                  <p className="text-xl font-bold text-gray-900 dark:text-gray-100">
                    {statement.total_count}
                  </p>
                </div>
              </Card>
            </div>
          )}

          <Card title="Transactions">
            {isLoading ? (
              <div className="flex justify-center py-8">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600" />
              </div>
            ) : statement ? (
              <TransactionList
                transactions={statement.transactions}
                currentAccountId={primaryAccount.id}
              />
            ) : null}
          </Card>
        </div>
      </main>
    </div>
  )
}
