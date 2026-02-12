import { useState } from 'react'
import { useNavigate } from '@tanstack/react-router'
import { v4 as uuidv4 } from 'uuid'
import { usePrimaryAccount } from '@/hooks/useAccounts'
import { useDeposit } from '@/hooks/useTransactions'
import Button from '@/components/common/Button'
import Input from '@/components/common/Input'

export default function DepositForm() {
  const navigate = useNavigate()
  const { primaryAccount } = usePrimaryAccount()
  const depositMutation = useDeposit()

  const [amount, setAmount] = useState('')
  const [description, setDescription] = useState('')
  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const [successMessage, setSuccessMessage] = useState<string | null>(null)

  const currentBalance = primaryAccount?.balance ?? '0.00'
  const accountId = primaryAccount?.id ?? ''

  const numAmount = parseFloat(amount)
  const isValidAmount = !isNaN(numAmount) && numAmount > 0
  const canSubmit = isValidAmount && !depositMutation.isPending

  function formatCurrency(value: string): string {
    const num = parseFloat(value)
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: primaryAccount?.currency || 'BRL',
    }).format(num)
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()

    if (!canSubmit) {
      setErrorMessage('Please enter a valid amount')
      return
    }

    setErrorMessage(null)
    setSuccessMessage(null)

    try {
      const result = await depositMutation.mutateAsync({
        accountId,
        request: {
          idempotency_key: uuidv4(),
          amount,
          description,
        },
      })

      if (result.success) {
        navigate({ to: '/dashboard' })
      } else {
        setErrorMessage(result.error_message || 'Deposit failed')
      }
    } catch (e: unknown) {
      const err = e as { response?: { data?: { error_message?: string; error?: string } } }
      setErrorMessage(err.response?.data?.error_message || err.response?.data?.error || 'Deposit failed')
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div className="bg-gray-50 dark:bg-gray-700 rounded-lg p-4">
        <p className="text-sm text-gray-600 dark:text-gray-400">Your Balance</p>
        <p className="text-2xl font-bold text-primary-600">{formatCurrency(currentBalance)}</p>
      </div>

      <Input
        value={amount}
        onChange={setAmount}
        label="Amount"
        type="number"
        step="0.01"
        min="0.01"
        placeholder="0.00"
        required
        error={amount && !isValidAmount ? 'Please enter a valid amount greater than 0' : undefined}
      />

      <Input
        value={description}
        onChange={setDescription}
        label="Description (optional)"
        placeholder="What's this deposit for?"
      />

      {errorMessage && (
        <div className="bg-red-50 dark:bg-red-900/30 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 px-4 py-3 rounded-lg text-sm">
          {errorMessage}
        </div>
      )}

      {successMessage && (
        <div className="bg-green-50 dark:bg-green-900/30 border border-green-200 dark:border-green-800 text-green-700 dark:text-green-400 px-4 py-3 rounded-lg text-sm">
          {successMessage}
        </div>
      )}

      <Button type="submit" disabled={!canSubmit} loading={depositMutation.isPending} className="w-full">
        Deposit
      </Button>
    </form>
  )
}
