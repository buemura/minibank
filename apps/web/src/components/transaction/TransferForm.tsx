import { useState, useEffect } from 'react'
import { useNavigate } from '@tanstack/react-router'
import { v4 as uuidv4 } from 'uuid'
import { usePrimaryAccount, useAccountLookup } from '@/hooks/useAccounts'
import { useTransfer } from '@/hooks/useTransactions'
import Button from '@/components/common/Button'
import Input from '@/components/common/Input'

export default function TransferForm() {
  const navigate = useNavigate()
  const { primaryAccount } = usePrimaryAccount()
  const transferMutation = useTransfer()

  const [destinationAccount, setDestinationAccount] = useState('')
  const [debouncedAccount, setDebouncedAccount] = useState('')
  const [amount, setAmount] = useState('')
  const [description, setDescription] = useState('')
  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const [successMessage, setSuccessMessage] = useState<string | null>(null)

  useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedAccount(destinationAccount.trim())
    }, 500)
    return () => clearTimeout(timer)
  }, [destinationAccount])

  const accountLookup = useAccountLookup(debouncedAccount)

  const currentBalance = primaryAccount?.balance ?? '0.00'
  const accountId = primaryAccount?.id ?? ''

  const numAmount = parseFloat(amount)
  const numBalance = parseFloat(currentBalance)
  const isValidAmount = !isNaN(numAmount) && numAmount > 0 && numAmount <= numBalance
  const canSubmit =
    accountLookup.isSuccess &&
    !!accountLookup.data &&
    isValidAmount &&
    !transferMutation.isPending

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
      setErrorMessage('Please fill in all required fields correctly')
      return
    }

    setErrorMessage(null)
    setSuccessMessage(null)

    try {
      const result = await transferMutation.mutateAsync({
        accountId,
        request: {
          idempotency_key: uuidv4(),
          destination_account_number: destinationAccount,
          amount,
          description,
        },
      })

      if (result.success) {
        navigate({ to: '/dashboard' })
      } else {
        setErrorMessage(result.error_message || 'Transfer failed')
      }
    } catch (e: unknown) {
      const err = e as { response?: { data?: { error_message?: string; error?: string } } }
      setErrorMessage(err.response?.data?.error_message || err.response?.data?.error || 'Transfer failed')
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div className="bg-gray-50 dark:bg-gray-700 rounded-lg p-4">
        <p className="text-sm text-gray-600 dark:text-gray-400">Your Balance</p>
        <p className="text-2xl font-bold text-primary-600">{formatCurrency(currentBalance)}</p>
      </div>

      <Input
        value={destinationAccount}
        onChange={setDestinationAccount}
        label="Destination Account Number"
        placeholder="Enter the recipient's account number"
        required
      />

      {accountLookup.isFetching && (
        <div className="text-sm text-gray-500 dark:text-gray-400">Looking up account...</div>
      )}

      {accountLookup.isSuccess && accountLookup.data && (
        <div className="bg-blue-50 dark:bg-blue-900/30 border border-blue-200 dark:border-blue-800 rounded-lg p-4">
          <p className="text-sm text-gray-600 dark:text-gray-400 mb-1">Recipient</p>
          <p className="font-semibold text-gray-900 dark:text-gray-100">{accountLookup.data.owner_name}</p>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            Agency {accountLookup.data.agency} &middot; Account {accountLookup.data.account_number}
          </p>
        </div>
      )}

      {accountLookup.isError && debouncedAccount && (
        <div className="text-sm text-red-500">Account not found. Please check the account number.</div>
      )}

      <Input
        value={amount}
        onChange={setAmount}
        label="Amount"
        type="number"
        step="0.01"
        min="0.01"
        placeholder="0.00"
        required
        error={amount && !isValidAmount ? 'Invalid amount or insufficient balance' : undefined}
      />

      <Input
        value={description}
        onChange={setDescription}
        label="Description (optional)"
        placeholder="What's this transfer for?"
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

      <Button type="submit" disabled={!canSubmit} loading={transferMutation.isPending} className="w-full">
        Transfer
      </Button>
    </form>
  )
}
