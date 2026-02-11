<script setup lang="ts">
import type { Transaction } from '@/types'

interface Props {
  transactions: Transaction[]
  currentAccountId?: string
}

const props = defineProps<Props>()

function formatCurrency(value: string, currency: string): string {
  const num = parseFloat(value)
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: currency,
  }).format(num)
}

function formatDate(dateString: string): string {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

function isOutgoing(transaction: Transaction): boolean {
  return transaction.source_account_id === props.currentAccountId
}

function getTransactionIcon(transaction: Transaction): string {
  if (transaction.type === 'TRANSFER') {
    return isOutgoing(transaction) ? 'arrow-up' : 'arrow-down'
  }
  return 'circle'
}

function getStatusColor(status: string): string {
  switch (status) {
    case 'COMPLETED':
      return 'text-green-600 bg-green-100 dark:bg-green-900/30 dark:text-green-400'
    case 'PENDING':
      return 'text-yellow-600 bg-yellow-100 dark:bg-yellow-900/30 dark:text-yellow-400'
    case 'FAILED':
      return 'text-red-600 bg-red-100 dark:bg-red-900/30 dark:text-red-400'
    default:
      return 'text-gray-600 bg-gray-100 dark:bg-gray-700 dark:text-gray-400'
  }
}
</script>

<template>
  <div class="space-y-4">
    <div v-if="transactions.length === 0" class="text-center py-8 text-gray-500 dark:text-gray-400">
      No transactions found
    </div>

    <div
      v-for="transaction in transactions"
      :key="transaction.id"
      class="flex items-center justify-between p-4 bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 hover:shadow-sm transition-shadow"
    >
      <div class="flex items-center space-x-4">
        <div
          :class="[
            'w-10 h-10 rounded-full flex items-center justify-center',
            isOutgoing(transaction) ? 'bg-red-100 dark:bg-red-900/30 text-red-600 dark:text-red-400' : 'bg-green-100 dark:bg-green-900/30 text-green-600 dark:text-green-400',
          ]"
        >
          <svg v-if="getTransactionIcon(transaction) === 'arrow-up'" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 10l7-7m0 0l7 7m-7-7v18" />
          </svg>
          <svg v-else-if="getTransactionIcon(transaction) === 'arrow-down'" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 14l-7 7m0 0l-7-7m7 7V3" />
          </svg>
        </div>

        <div>
          <p class="font-medium text-gray-900 dark:text-gray-100">
            {{ isOutgoing(transaction) ? 'Sent to' : 'Received from' }}
            {{ isOutgoing(transaction) ? transaction.destination_account_number : transaction.source_account_number }}
          </p>
          <p v-if="transaction.description" class="text-sm text-gray-500 dark:text-gray-400">{{ transaction.description }}</p>
          <p class="text-xs text-gray-400 dark:text-gray-500">{{ formatDate(transaction.created_at) }}</p>
        </div>
      </div>

      <div class="text-right">
        <p
          :class="[
            'font-semibold',
            isOutgoing(transaction) ? 'text-red-600 dark:text-red-400' : 'text-green-600 dark:text-green-400',
          ]"
        >
          {{ isOutgoing(transaction) ? '-' : '+' }}{{ formatCurrency(transaction.amount, transaction.currency) }}
        </p>
        <span
          :class="[
            'inline-block px-2 py-1 text-xs rounded-full',
            getStatusColor(transaction.status),
          ]"
        >
          {{ transaction.status }}
        </span>
      </div>
    </div>
  </div>
</template>
