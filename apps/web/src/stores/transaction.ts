import { defineStore } from 'pinia'
import { ref } from 'vue'
import { transactionsApi } from '@/api/transactions'
import type { Transaction, TransferRequest, TransferResponse, DepositRequest, DepositResponse, StatementResponse } from '@/types'

export const useTransactionStore = defineStore('transaction', () => {
  const transactions = ref<Transaction[]>([])
  const statement = ref<StatementResponse | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  async function transfer(accountId: string, request: TransferRequest): Promise<TransferResponse> {
    isLoading.value = true
    error.value = null

    try {
      const response = await transactionsApi.transfer(accountId, request)
      if (response.success && response.transaction) {
        transactions.value.unshift(response.transaction)
      }
      return response
    } catch (e: unknown) {
      const err = e as { response?: { data?: { error?: string; error_message?: string } } }
      error.value = err.response?.data?.error_message || err.response?.data?.error || 'Transfer failed'
      return {
        success: false,
        error_message: error.value,
      }
    } finally {
      isLoading.value = false
    }
  }

  async function deposit(accountId: string, request: DepositRequest): Promise<DepositResponse> {
    isLoading.value = true
    error.value = null

    try {
      const response = await transactionsApi.deposit(accountId, request)
      if (response.success && response.transaction) {
        transactions.value.unshift(response.transaction)
      }
      return response
    } catch (e: unknown) {
      const err = e as { response?: { data?: { error?: string; error_message?: string } } }
      error.value = err.response?.data?.error_message || err.response?.data?.error || 'Deposit failed'
      return {
        success: false,
        error_message: error.value,
      }
    } finally {
      isLoading.value = false
    }
  }

  async function fetchTransactionHistory(
    accountId: string,
    page: number = 1,
    pageSize: number = 20
  ): Promise<void> {
    isLoading.value = true
    error.value = null

    try {
      const response = await transactionsApi.getTransactionHistory(accountId, page, pageSize)
      transactions.value = response.transactions
    } catch (e: unknown) {
      const err = e as { response?: { data?: { error?: string } } }
      error.value = err.response?.data?.error || 'Failed to fetch transactions'
    } finally {
      isLoading.value = false
    }
  }

  async function fetchStatement(
    accountId: string,
    startDate: string,
    endDate: string
  ): Promise<void> {
    isLoading.value = true
    error.value = null

    try {
      statement.value = await transactionsApi.getStatement(accountId, startDate, endDate)
    } catch (e: unknown) {
      const err = e as { response?: { data?: { error?: string } } }
      error.value = err.response?.data?.error || 'Failed to fetch statement'
    } finally {
      isLoading.value = false
    }
  }

  function clearTransactions() {
    transactions.value = []
    statement.value = null
  }

  return {
    transactions,
    statement,
    isLoading,
    error,
    transfer,
    deposit,
    fetchTransactionHistory,
    fetchStatement,
    clearTransactions,
  }
})
