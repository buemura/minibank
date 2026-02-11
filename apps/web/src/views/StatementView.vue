<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useAccountStore } from '@/stores/account'
import { useTransactionStore } from '@/stores/transaction'
import Navbar from '@/components/layout/Navbar.vue'
import TransactionList from '@/components/transaction/TransactionList.vue'
import Card from '@/components/common/Card.vue'
import Button from '@/components/common/Button.vue'
import Input from '@/components/common/Input.vue'

const accountStore = useAccountStore()
const transactionStore = useTransactionStore()

const startDate = ref(getDefaultStartDate())
const endDate = ref(getDefaultEndDate())

function getDefaultStartDate(): string {
  const date = new Date()
  date.setMonth(date.getMonth() - 1)
  return date.toISOString().split('T')[0]
}

function getDefaultEndDate(): string {
  return new Date().toISOString().split('T')[0]
}

const statement = computed(() => transactionStore.statement)

onMounted(async () => {
  if (accountStore.accounts.length === 0) {
    await accountStore.fetchAccounts()
  }
  if (accountStore.primaryAccount) {
    await fetchStatement()
  }
})

async function fetchStatement() {
  if (!accountStore.primaryAccount) return
  await transactionStore.fetchStatement(
    accountStore.primaryAccount.id,
    startDate.value,
    endDate.value
  )
}

function formatCurrency(value: string, currency: string = 'BRL'): string {
  const num = parseFloat(value)
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: currency,
  }).format(num)
}
</script>

<template>
  <div class="min-h-screen">
    <Navbar />

    <main class="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
      <div class="mb-8">
        <h1 class="text-2xl font-bold text-gray-900 dark:text-gray-100">Account Statement</h1>
        <p class="text-gray-600 dark:text-gray-400">View your transaction history.</p>
      </div>

      <div v-if="!accountStore.primaryAccount" class="text-center py-12">
        <Card>
          <div class="py-8 text-gray-500 dark:text-gray-400">
            You need an account to view statements.
            <router-link to="/dashboard" class="text-primary-600 hover:text-primary-700">
              Go to Dashboard
            </router-link>
          </div>
        </Card>
      </div>

      <div v-else class="space-y-6">
        <Card title="Filter">
          <form @submit.prevent="fetchStatement" class="flex flex-wrap gap-4 items-end">
            <div class="flex-1 min-w-[200px]">
              <Input
                v-model="startDate"
                label="Start Date"
                type="date"
              />
            </div>
            <div class="flex-1 min-w-[200px]">
              <Input
                v-model="endDate"
                label="End Date"
                type="date"
              />
            </div>
            <Button type="submit" :loading="transactionStore.isLoading">
              Apply
            </Button>
          </form>
        </Card>

        <div v-if="statement" class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <Card>
            <div class="text-center">
              <p class="text-sm text-gray-600 dark:text-gray-400">Opening Balance</p>
              <p class="text-xl font-bold text-gray-900 dark:text-gray-100">
                {{ formatCurrency(statement.opening_balance) }}
              </p>
            </div>
          </Card>
          <Card>
            <div class="text-center">
              <p class="text-sm text-gray-600 dark:text-gray-400">Closing Balance</p>
              <p class="text-xl font-bold text-primary-600">
                {{ formatCurrency(statement.closing_balance) }}
              </p>
            </div>
          </Card>
          <Card>
            <div class="text-center">
              <p class="text-sm text-gray-600 dark:text-gray-400">Total Transactions</p>
              <p class="text-xl font-bold text-gray-900 dark:text-gray-100">
                {{ statement.total_count }}
              </p>
            </div>
          </Card>
        </div>

        <Card title="Transactions">
          <div v-if="transactionStore.isLoading" class="flex justify-center py-8">
            <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
          </div>

          <TransactionList
            v-else-if="statement"
            :transactions="statement.transactions"
            :current-account-id="accountStore.primaryAccount?.id"
          />
        </Card>
      </div>
    </main>
  </div>
</template>
