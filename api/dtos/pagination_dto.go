package dtos

type Pagination struct {
	Page  int `json:"page" binding:"required"`
	Limit int `json:"limit" binding:"required"`
}
