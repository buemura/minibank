<script setup lang="ts">
import { onMounted } from 'vue'
import { useAccountStore } from '@/stores/account'
import Navbar from '@/components/layout/Navbar.vue'
import TransferForm from '@/components/transaction/TransferForm.vue'
import Card from '@/components/common/Card.vue'

const accountStore = useAccountStore()

onMounted(async () => {
  if (accountStore.accounts.length === 0) {
    await accountStore.fetchAccounts()
  }
})
</script>

<template>
  <div class="min-h-screen">
    <Navbar />

    <main class="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
      <div class="mb-8">
        <h1 class="text-2xl font-bold text-gray-900 dark:text-gray-100">Transfer Money</h1>
        <p class="text-gray-600 dark:text-gray-400">Send money to another account instantly.</p>
      </div>

      <div class="max-w-lg">
        <Card title="New Transfer">
          <div v-if="accountStore.isLoading" class="flex justify-center py-8">
            <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
          </div>

          <div v-else-if="!accountStore.primaryAccount" class="text-center py-8 text-gray-500 dark:text-gray-400">
            You need an account to make transfers.
            <router-link to="/dashboard" class="text-primary-600 hover:text-primary-700">
              Go to Dashboard
            </router-link>
          </div>

          <TransferForm v-else />
        </Card>
      </div>
    </main>
  </div>
</template>
