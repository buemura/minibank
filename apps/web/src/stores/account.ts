import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { accountsApi } from '@/api/accounts'
import type { Account } from '@/types'

export const useAccountStore = defineStore('account', () => {
  const accounts = ref<Account[]>([])
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  const primaryAccount = computed(() => accounts.value.find((a) => a.account_type === 'checking') || accounts.value[0])

  async function fetchAccounts(): Promise<void> {
    isLoading.value = true
    error.value = null

    try {
      accounts.value = await accountsApi.getAccounts()
    } catch (e: unknown) {
      const err = e as { response?: { data?: { error?: string } } }
      error.value = err.response?.data?.error || 'Failed to fetch accounts'
    } finally {
      isLoading.value = false
    }
  }

  async function createAccount(accountType: string = 'checking'): Promise<Account | null> {
    isLoading.value = true
    error.value = null

    try {
      const account = await accountsApi.createAccount(accountType)
      accounts.value.push(account)
      return account
    } catch (e: unknown) {
      const err = e as { response?: { data?: { error?: string } } }
      error.value = err.response?.data?.error || 'Failed to create account'
      return null
    } finally {
      isLoading.value = false
    }
  }

  async function refreshBalance(accountId: string): Promise<void> {
    try {
      const { balance } = await accountsApi.getBalance(accountId)
      const account = accounts.value.find((a) => a.id === accountId)
      if (account) {
        account.balance = balance
      }
    } catch {
      // Silent fail - balance will refresh on next full fetch
    }
  }

  function clearAccounts() {
    accounts.value = []
  }

  return {
    accounts,
    isLoading,
    error,
    primaryAccount,
    fetchAccounts,
    createAccount,
    refreshBalance,
    clearAccounts,
  }
})
