<script setup lang="ts">
import { ref, computed } from 'vue'
import { v4 as uuidv4 } from 'uuid'
import { useTransactionStore } from '@/stores/transaction'
import { useAccountStore } from '@/stores/account'
import Button from '@/components/common/Button.vue'
import Input from '@/components/common/Input.vue'

const transactionStore = useTransactionStore()
const accountStore = useAccountStore()

const destinationAccount = ref('')
const amount = ref('')
const description = ref('')
const isSubmitting = ref(false)
const errorMessage = ref<string | null>(null)
const successMessage = ref<string | null>(null)

const currentBalance = computed(() => accountStore.primaryAccount?.balance ?? '0.00')
const accountId = computed(() => accountStore.primaryAccount?.id ?? '')

const isValidAmount = computed(() => {
  const numAmount = parseFloat(amount.value)
  const numBalance = parseFloat(currentBalance.value)
  return !isNaN(numAmount) && numAmount > 0 && numAmount <= numBalance
})

const canSubmit = computed(() => {
  return destinationAccount.value.trim() !== '' && isValidAmount.value && !isSubmitting.value
})

async function handleSubmit() {
  if (!canSubmit.value) {
    errorMessage.value = 'Please fill in all required fields correctly'
    return
  }

  isSubmitting.value = true
  errorMessage.value = null
  successMessage.value = null

  const idempotencyKey = uuidv4()

  const result = await transactionStore.transfer(accountId.value, {
    idempotency_key: idempotencyKey,
    destination_account_number: destinationAccount.value,
    amount: amount.value,
    description: description.value,
  })

  if (result.success) {
    successMessage.value = `Transfer of ${formatCurrency(amount.value)} completed successfully!`
    destinationAccount.value = ''
    amount.value = ''
    description.value = ''
    await accountStore.fetchAccounts()
  } else {
    errorMessage.value = result.error_message || 'Transfer failed'
  }

  isSubmitting.value = false
}

function formatCurrency(value: string): string {
  const num = parseFloat(value)
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: accountStore.primaryAccount?.currency || 'BRL',
  }).format(num)
}
</script>

<template>
  <form @submit.prevent="handleSubmit" class="space-y-6">
    <div class="bg-gray-50 dark:bg-gray-700 rounded-lg p-4">
      <p class="text-sm text-gray-600 dark:text-gray-400">Your Balance</p>
      <p class="text-2xl font-bold text-primary-600">{{ formatCurrency(currentBalance) }}</p>
    </div>

    <Input
      v-model="destinationAccount"
      label="Destination Account Number"
      placeholder="Enter the recipient's account number"
      required
    />

    <Input
      v-model="amount"
      label="Amount"
      type="number"
      step="0.01"
      min="0.01"
      placeholder="0.00"
      required
      :error="amount && !isValidAmount ? 'Invalid amount or insufficient balance' : undefined"
    />

    <Input
      v-model="description"
      label="Description (optional)"
      placeholder="What's this transfer for?"
    />

    <div v-if="errorMessage" class="bg-red-50 dark:bg-red-900/30 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 px-4 py-3 rounded-lg text-sm">
      {{ errorMessage }}
    </div>

    <div v-if="successMessage" class="bg-green-50 dark:bg-green-900/30 border border-green-200 dark:border-green-800 text-green-700 dark:text-green-400 px-4 py-3 rounded-lg text-sm">
      {{ successMessage }}
    </div>

    <Button type="submit" :disabled="!canSubmit" :loading="isSubmitting" class="w-full">
      Transfer
    </Button>
  </form>
</template>
