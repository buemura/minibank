import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { accountsApi } from '@/api/accounts'

export function useAccounts() {
  return useQuery({
    queryKey: ['accounts'],
    queryFn: accountsApi.getAccounts,
  })
}

export function usePrimaryAccount() {
  const { data: accounts, ...rest } = useAccounts()
  const primaryAccount = accounts?.find((a) => a.account_type === 'checking') ?? accounts?.[0] ?? null
  return { primaryAccount, accounts, ...rest }
}

export function useCreateAccount() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (accountType: string = 'checking') => accountsApi.createAccount(accountType),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['accounts'] })
    },
  })
}

export function useAccountLookup(accountNumber: string) {
  return useQuery({
    queryKey: ['account-lookup', accountNumber],
    queryFn: () => accountsApi.lookupByNumber(accountNumber),
    enabled: accountNumber.trim().length > 0,
    retry: false,
  })
}

export function useRefreshBalance() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (accountId: string) => accountsApi.getBalance(accountId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['accounts'] })
    },
  })
}
