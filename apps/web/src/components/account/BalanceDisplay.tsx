import type { Account } from '@/types'

interface BalanceDisplayProps {
  account: Account
}

function formatCurrency(value: string, currency: string): string {
  const num = parseFloat(value)
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: currency,
  }).format(num)
}

export default function BalanceDisplay({ account }: BalanceDisplayProps) {
  return (
    <div className="bg-gradient-to-br from-primary-600 to-primary-800 rounded-xl p-6 text-white">
      <div className="flex justify-between items-start mb-4">
        <div>
          <p className="text-primary-100 text-sm">Available Balance</p>
          <p className="text-3xl font-bold mt-1">{formatCurrency(account.balance, account.currency)}</p>
        </div>
        <div className="text-right">
          <p className="text-primary-100 text-xs">Account</p>
          <p className="font-mono text-sm">{account.account_number}</p>
        </div>
      </div>
      <div className="flex justify-between text-sm text-primary-100">
        <span>Agency: {account.agency}</span>
        <span className="capitalize">{account.account_type}</span>
      </div>
    </div>
  )
}
