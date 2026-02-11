import client from './client'
import type { Transaction, TransferRequest, TransferResponse, DepositRequest, DepositResponse, StatementResponse } from '@/types'

export const transactionsApi = {
  async transfer(accountId: string, request: TransferRequest): Promise<TransferResponse> {
    const { data } = await client.post<TransferResponse>(`/accounts/${accountId}/transfers`, request)
    return data
  },

  async deposit(accountId: string, request: DepositRequest): Promise<DepositResponse> {
    const { data } = await client.post<DepositResponse>(`/accounts/${accountId}/deposits`, request)
    return data
  },

  async getTransaction(transactionId: string): Promise<Transaction> {
    const { data } = await client.get<{ transaction: Transaction }>(`/transactions/${transactionId}`)
    return data.transaction
  },

  async getTransactionHistory(
    accountId: string,
    page: number = 1,
    pageSize: number = 20
  ): Promise<{ transactions: Transaction[]; total_count: number }> {
    const { data } = await client.get<{ transactions: Transaction[]; total_count: number }>(
      `/accounts/${accountId}/transactions`,
      { params: { page, page_size: pageSize } }
    )
    return data
  },

  async getStatement(
    accountId: string,
    startDate: string,
    endDate: string,
    page: number = 1,
    pageSize: number = 50
  ): Promise<StatementResponse> {
    const { data } = await client.get<StatementResponse>(`/accounts/${accountId}/statement`, {
      params: { start_date: startDate, end_date: endDate, page, page_size: pageSize },
    })
    return data
  },
}
