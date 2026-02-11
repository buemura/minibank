<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import ThemeToggle from '@/components/common/ThemeToggle.vue'

const authStore = useAuthStore()
const router = useRouter()

function handleLogout() {
  authStore.logout()
  router.push('/login')
}
</script>

<template>
  <nav class="bg-white dark:bg-gray-800 shadow-sm border-b border-gray-200 dark:border-gray-700">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between h-16">
        <div class="flex items-center">
          <router-link to="/dashboard" class="flex items-center space-x-2">
            <span class="text-2xl font-bold text-primary-600">Vibe</span>
            <span class="text-2xl font-light text-gray-600 dark:text-gray-400">Bank</span>
          </router-link>
          <div class="hidden sm:ml-10 sm:flex sm:space-x-4">
            <router-link
              to="/dashboard"
              class="px-3 py-2 rounded-md text-sm font-medium text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-100 hover:bg-gray-100 dark:hover:bg-gray-700"
              active-class="text-primary-600 bg-primary-50 dark:bg-primary-900/30"
            >
              Dashboard
            </router-link>
            <router-link
              to="/transfer"
              class="px-3 py-2 rounded-md text-sm font-medium text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-100 hover:bg-gray-100 dark:hover:bg-gray-700"
              active-class="text-primary-600 bg-primary-50 dark:bg-primary-900/30"
            >
              Transfer
            </router-link>
            <router-link
              to="/statement"
              class="px-3 py-2 rounded-md text-sm font-medium text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-100 hover:bg-gray-100 dark:hover:bg-gray-700"
              active-class="text-primary-600 bg-primary-50 dark:bg-primary-900/30"
            >
              Statement
            </router-link>
          </div>
        </div>
        <div class="flex items-center space-x-4">
          <span class="text-sm text-gray-600 dark:text-gray-400">{{ authStore.currentUser?.full_name }}</span>
          <ThemeToggle />
          <button
            @click="handleLogout"
            class="px-3 py-2 rounded-md text-sm font-medium text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-100 hover:bg-gray-100 dark:hover:bg-gray-700"
          >
            Logout
          </button>
        </div>
      </div>
    </div>
  </nav>
</template>
