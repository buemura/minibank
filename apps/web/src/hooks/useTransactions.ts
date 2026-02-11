import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { transactionsApi } from '@/api/transactions'
import type { TransferRequest, DepositRequest } from '@/types'

export function useTransactionHistory(accountId: string | undefined, page: number = 1, pageSize: number = 20) {
  return useQuery({
    queryKey: ['transactions', accountId, page, pageSize],
    queryFn: () => transactionsApi.getTransactionHistory(accountId!, page, pageSize),
    enabled: !!accountId,
  })
}

export function useStatement(accountId: string | undefined, startDate: string, endDate: string, enabled: boolean = true) {
  return useQuery({
    queryKey: ['statement', accountId, startDate, endDate],
    queryFn: () => transactionsApi.getStatement(accountId!, startDate, endDate),
    enabled: !!accountId && enabled,
  })
}

export function useTransfer() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ accountId, request }: { accountId: string; request: TransferRequest }) =>
      transactionsApi.transfer(accountId, request),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['accounts'] })
      queryClient.invalidateQueries({ queryKey: ['transactions'] })
      queryClient.invalidateQueries({ queryKey: ['statement'] })
    },
  })
}

export function useDeposit() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ accountId, request }: { accountId: string; request: DepositRequest }) =>
      transactionsApi.deposit(accountId, request),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['accounts'] })
      queryClient.invalidateQueries({ queryKey: ['transactions'] })
      queryClient.invalidateQueries({ queryKey: ['statement'] })
    },
  })
}
