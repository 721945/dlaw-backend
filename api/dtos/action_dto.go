package dtos

import (
	"github.com/721945/dlaw-backend/models"
	"time"
)

type ActionDto struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateActionDto struct {
	Name     string   `json:"name"`
	NameList []string `json:"nameList"`
}

func ToActionDto(action models.Action) ActionDto {
	return ActionDto{
		ID:        action.ID.String(),
		Name:      action.Name,
		CreatedAt: action.CreatedAt,
		UpdatedAt: action.UpdatedAt,
	}
}

func (c CreateActionDto) ToModel() models.Action {
	return models.Action{
		Name: c.Name,
	}
}
func (c CreateActionDto) ToModelList() []models.Action {
	actions := make([]models.Action, len(c.NameList))

	for i, name := range c.NameList {
		actions[i] = models.Action{
			Name: name,
		}
	}

	return actions
}

type UpdateActionDto struct {
	Name string `json:"name"`
}

func (u UpdateActionDto) ToModel() models.Action {
	return models.Action{
		Name: u.Name,
	}
}
