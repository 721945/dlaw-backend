package dtos

import "github.com/721945/dlaw-backend/models"

type CreateActionDto struct {
	Name string `json:"name"`
}

func (c CreateActionDto) ToModel() models.Action {
	return models.Action{
		Name: c.Name,
	}
}
