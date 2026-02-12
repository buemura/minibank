import type { Transaction } from '@/types'

interface TransactionListProps {
  transactions: Transaction[]
  currentAccountId?: string
}

function formatCurrency(value: string, currency: string): string {
  const num = parseFloat(value)
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: currency,
  }).format(num)
}

function formatDate(dateString: string): string {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

function getStatusColor(status: string): string {
  switch (status) {
    case 'COMPLETED':
      return 'text-green-600 bg-green-100 dark:bg-green-900/30 dark:text-green-400'
    case 'PENDING':
      return 'text-yellow-600 bg-yellow-100 dark:bg-yellow-900/30 dark:text-yellow-400'
    case 'FAILED':
      return 'text-red-600 bg-red-100 dark:bg-red-900/30 dark:text-red-400'
    default:
      return 'text-gray-600 bg-gray-100 dark:bg-gray-700 dark:text-gray-400'
  }
}

export default function TransactionList({ transactions, currentAccountId }: TransactionListProps) {
  function isOutgoing(transaction: Transaction): boolean {
    return transaction.source_account_id === currentAccountId
  }

  function getTransactionIcon(transaction: Transaction): string {
    if (transaction.type === 'WITHDRAWAL') return 'arrow-up'
    if (transaction.type === 'DEPOSIT') return 'arrow-down'
    return isOutgoing(transaction) ? 'arrow-up' : 'arrow-down'
  }

  function getTransactionLabel(transaction: Transaction): string {
    if (transaction.type === 'WITHDRAWAL') return 'Withdrawal'
    if (transaction.type === 'DEPOSIT') return 'Deposit'
    if (transaction.type === 'FEE') return 'Withdrawal fee'
    if (isOutgoing(transaction)) {
      return `Sent to ${transaction.destination_account_number}`
    }
    return `Received from ${transaction.source_account_number}`
  }

  if (transactions.length === 0) {
    return (
      <div className="text-center py-8 text-gray-500 dark:text-gray-400">
        No transactions found
      </div>
    )
  }

  return (
    <div className="space-y-4">
      {transactions.map((transaction) => (
        <div
          key={transaction.id}
          className="flex items-center justify-between p-4 bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 hover:shadow-sm transition-shadow"
        >
          <div className="flex items-center space-x-4">
            <div
              className={`w-10 h-10 rounded-full flex items-center justify-center ${
                transaction.type === 'DEPOSIT'
                  ? 'bg-green-100 dark:bg-green-900/30 text-green-600 dark:text-green-400'
                  : transaction.type === 'WITHDRAWAL'
                    ? 'bg-red-100 dark:bg-red-900/30 text-red-600 dark:text-red-400'
                    : isOutgoing(transaction)
                      ? 'bg-red-100 dark:bg-red-900/30 text-red-600 dark:text-red-400'
                      : 'bg-green-100 dark:bg-green-900/30 text-green-600 dark:text-green-400'
              }`}
            >
              {getTransactionIcon(transaction) === 'arrow-up' ? (
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 10l7-7m0 0l7 7m-7-7v18" />
                </svg>
              ) : (
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 14l-7 7m0 0l-7-7m7 7V3" />
                </svg>
              )}
            </div>

            <div>
              <p className="font-medium text-gray-900 dark:text-gray-100">
                {getTransactionLabel(transaction)}
              </p>
              {transaction.description && (
                <p className="text-sm text-gray-500 dark:text-gray-400">{transaction.description}</p>
              )}
              <p className="text-xs text-gray-400 dark:text-gray-500">{formatDate(transaction.created_at)}</p>
            </div>
          </div>

          <div className="text-right">
            <p
              className={`font-semibold ${
                transaction.type === 'DEPOSIT'
                  ? 'text-green-600 dark:text-green-400'
                  : transaction.type === 'WITHDRAWAL'
                    ? 'text-red-600 dark:text-red-400'
                    : isOutgoing(transaction)
                      ? 'text-red-600 dark:text-red-400'
                      : 'text-green-600 dark:text-green-400'
              }`}
            >
              {transaction.type === 'DEPOSIT' ? '+' : transaction.type === 'WITHDRAWAL' ? '-' : isOutgoing(transaction) ? '-' : '+'}
              {formatCurrency(transaction.amount, transaction.currency)}
            </p>
            <span className={`inline-block px-2 py-1 text-xs rounded-full ${getStatusColor(transaction.status)}`}>
              {transaction.status}
            </span>
          </div>
        </div>
      ))}
    </div>
  )
}
