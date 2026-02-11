<script setup lang="ts">
import { onMounted } from 'vue'
import { useAccountStore } from '@/stores/account'
import { useTransactionStore } from '@/stores/transaction'
import Navbar from '@/components/layout/Navbar.vue'
import BalanceDisplay from '@/components/account/BalanceDisplay.vue'
import TransactionList from '@/components/transaction/TransactionList.vue'
import Card from '@/components/common/Card.vue'
import Button from '@/components/common/Button.vue'

const accountStore = useAccountStore()
const transactionStore = useTransactionStore()

onMounted(async () => {
  await accountStore.fetchAccounts()

  if (accountStore.primaryAccount) {
    await transactionStore.fetchTransactionHistory(accountStore.primaryAccount.id)
  }
})

async function handleCreateAccount() {
  await accountStore.createAccount('checking')
}
</script>

<template>
  <div class="min-h-screen">
    <Navbar />

    <main class="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
      <div class="mb-8">
        <h1 class="text-2xl font-bold text-gray-900 dark:text-gray-100">Dashboard</h1>
        <p class="text-gray-600 dark:text-gray-400">Welcome back! Here's your account overview.</p>
      </div>

      <div v-if="accountStore.isLoading" class="flex justify-center py-12">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
      </div>

      <div v-else-if="!accountStore.primaryAccount" class="text-center py-12">
        <Card>
          <div class="py-8">
            <h3 class="text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">No Account Yet</h3>
            <p class="text-gray-600 dark:text-gray-400 mb-6">Create your first account to start banking.</p>
            <Button @click="handleCreateAccount" :loading="accountStore.isLoading">
              Create Account
            </Button>
          </div>
        </Card>
      </div>

      <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div class="lg:col-span-1">
          <BalanceDisplay :account="accountStore.primaryAccount" />

          <div class="mt-6 grid grid-cols-3 gap-4">
            <router-link to="/transfer">
              <Button variant="primary" class="w-full">
                Transfer
              </Button>
            </router-link>
            <router-link to="/deposit">
              <Button variant="primary" class="w-full">
                Deposit
              </Button>
            </router-link>
            <router-link to="/statement">
              <Button variant="secondary" class="w-full">
                Statement
              </Button>
            </router-link>
          </div>
        </div>

        <div class="lg:col-span-2">
          <Card title="Recent Transactions">
            <TransactionList
              :transactions="transactionStore.transactions.slice(0, 5)"
              :current-account-id="accountStore.primaryAccount?.id"
            />

            <div v-if="transactionStore.transactions.length > 5" class="mt-4 text-center">
              <router-link to="/statement" class="text-primary-600 hover:text-primary-700 text-sm font-medium">
                View all transactions
              </router-link>
            </div>
          </Card>
        </div>
      </div>
    </main>
  </div>
</template>
