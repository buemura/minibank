package common

type PaginationMetaOut struct {
	Page       int `json:"page"`
	Items      int `json:"items"`
	TotalPages int `json:"total_pages"`
	TotalItems int `json:"total_items"`
}
