import client from './client'
import type { Account } from '@/types'

export const accountsApi = {
  async getAccounts(): Promise<Account[]> {
    const { data } = await client.get<{ accounts: Account[] }>('/accounts')
    return data.accounts
  },

  async createAccount(accountType: string = 'checking'): Promise<Account> {
    const { data } = await client.post<{ account: Account }>('/accounts', { account_type: accountType })
    return data.account
  },

  async getAccount(accountId: string): Promise<Account> {
    const { data } = await client.get<{ account: Account }>(`/accounts/${accountId}`)
    return data.account
  },

  async getBalance(accountId: string): Promise<{ balance: string; currency: string }> {
    const { data } = await client.get<{ balance: string; currency: string }>(`/accounts/${accountId}/balance`)
    return data
  },
}
