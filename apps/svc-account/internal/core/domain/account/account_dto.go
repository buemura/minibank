package account

type CreateAccountIn struct {
	OwnerName     string `json:"owner_name"`
	OwnerDocument string `json:"owner_document"`
}

type UpdateBalanceIn struct {
	ID         string `json:"id"`
	NewBalance int    `json:"new_balance"`
}
