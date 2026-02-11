<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import Button from '@/components/common/Button.vue'
import Input from '@/components/common/Input.vue'
import Card from '@/components/common/Card.vue'

const router = useRouter()
const authStore = useAuthStore()

const email = ref('')
const password = ref('')

async function handleSubmit() {
  const success = await authStore.login({
    email: email.value,
    password: password.value,
  })

  if (success) {
    router.push('/dashboard')
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <div class="text-center">
        <h1 class="text-4xl font-bold">
          <span class="text-primary-600">Vibe</span>
          <span class="text-gray-600 dark:text-gray-400 font-light">Bank</span>
        </h1>
        <p class="mt-2 text-gray-600 dark:text-gray-400">Sign in to your account</p>
      </div>

      <Card>
        <form @submit.prevent="handleSubmit" class="space-y-6">
          <Input
            v-model="email"
            label="Email"
            type="email"
            placeholder="you@example.com"
            required
          />

          <Input
            v-model="password"
            label="Password"
            type="password"
            placeholder="Enter your password"
            required
          />

          <div v-if="authStore.error" class="bg-red-50 dark:bg-red-900/30 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 px-4 py-3 rounded-lg text-sm">
            {{ authStore.error }}
          </div>

          <Button type="submit" :loading="authStore.isLoading" class="w-full">
            Sign In
          </Button>

          <p class="text-center text-sm text-gray-600 dark:text-gray-400">
            Don't have an account?
            <router-link to="/register" class="text-primary-600 hover:text-primary-700 font-medium">
              Sign up
            </router-link>
          </p>
        </form>
      </Card>
    </div>
  </div>
</template>
