export interface User {
  id: string
  email: string
  full_name: string
  phone: string
  status: string
  created_at: string
}

export interface Account {
  id: string
  user_id: string
  account_number: string
  agency: string
  account_type: string
  balance: string
  currency: string
  status: string
  created_at: string
  updated_at: string
}

export interface Transaction {
  id: string
  idempotency_key: string
  type: TransactionType
  status: TransactionStatus
  source_account_id: string
  source_account_number: string
  destination_account_id: string
  destination_account_number: string
  amount: string
  currency: string
  description: string
  processed_at: string
  created_at: string
}

export type TransactionType = 'TRANSFER' | 'DEPOSIT' | 'WITHDRAWAL' | 'FEE'
export type TransactionStatus = 'PENDING' | 'COMPLETED' | 'FAILED' | 'REVERSED'

export interface LoginCredentials {
  email: string
  password: string
}

export interface RegisterData {
  email: string
  password: string
  full_name: string
  phone: string
}

export interface AuthResponse {
  user: User
  access_token: string
  refresh_token: string
  expires_in: number
}

export interface TransferRequest {
  idempotency_key: string
  destination_account_number: string
  amount: string
  description?: string
}

export interface TransferResponse {
  success: boolean
  transaction?: Transaction
  error_code?: string
  error_message?: string
}

export interface DepositRequest {
  idempotency_key: string
  amount: string
  description?: string
}

export interface DepositResponse {
  success: boolean
  transaction?: Transaction
  error_code?: string
  error_message?: string
}

export interface WithdrawalRequest {
  idempotency_key: string
  amount: string
  description?: string
}

export interface WithdrawalResponse {
  success: boolean
  transaction?: Transaction
  error_code?: string
  error_message?: string
}

export interface StatementResponse {
  account_id: string
  opening_balance: string
  closing_balance: string
  transactions: Transaction[]
  total_count: number
  start_date: string
  end_date: string
}

export interface AccountLookup {
  account_number: string
  agency: string
  owner_name: string
}

export interface ChangePasswordRequest {
  current_password: string
  new_password: string
}

export interface ApiError {
  error: string
  error_code?: string
}
