<script setup lang="ts">
import type { Account } from '@/types'

interface Props {
  account: Account
}

defineProps<Props>()

function formatCurrency(value: string, currency: string): string {
  const num = parseFloat(value)
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: currency,
  }).format(num)
}
</script>

<template>
  <div class="bg-gradient-to-br from-primary-600 to-primary-800 rounded-xl p-6 text-white">
    <div class="flex justify-between items-start mb-4">
      <div>
        <p class="text-primary-100 text-sm">Available Balance</p>
        <p class="text-3xl font-bold mt-1">{{ formatCurrency(account.balance, account.currency) }}</p>
      </div>
      <div class="text-right">
        <p class="text-primary-100 text-xs">Account</p>
        <p class="font-mono text-sm">{{ account.account_number }}</p>
      </div>
    </div>
    <div class="flex justify-between text-sm text-primary-100">
      <span>Agency: {{ account.agency }}</span>
      <span class="capitalize">{{ account.account_type }}</span>
    </div>
  </div>
</template>
