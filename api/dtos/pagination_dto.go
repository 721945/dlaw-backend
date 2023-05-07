package dtos

type Pagination struct {
	Page  int `json:"page" binding:"required"`
	Limit int `json:"limit" binding:"required"`
}

type PaginationResponse struct {
	Total int64 `json:"total"`
	Page  int64 `json:"page"`
	Limit int64 `json:"limit"`
}
