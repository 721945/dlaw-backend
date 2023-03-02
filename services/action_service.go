package services

import (
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
	"github.com/721945/dlaw-backend/repositories"
)

type ActionService struct {
	logger     *libs.Logger
	actionRepo repositories.ActionRepository
}

func NewActionService(logger *libs.Logger, r repositories.ActionRepository) ActionService {
	return ActionService{logger: logger, actionRepo: r}
}

func (s *ActionService) GetActions() (actions []models.Action, err error) {
	return s.actionRepo.GetActions()
}

func (s *ActionService) GetAction(id uint) (action *models.Action, err error) {
	return s.actionRepo.GetAction(id)
}

func (s *ActionService) CreateAction(action models.Action) (models.Action, error) {
	return s.actionRepo.CreateAction(action)
}

func (s *ActionService) UpdateAction(id uint, action models.Action) error {
	return s.actionRepo.UpdateAction(id, action)
}

func (s *ActionService) DeleteAction(id uint) error {
	return s.actionRepo.DeleteAction(id)
}
