package repositories

import (
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
)

type ActionRepository struct {
	logger *libs.Logger
	db     libs.Database
}

func NewActionRepository(logger *libs.Logger, db libs.Database) ActionRepository {
	return ActionRepository{logger: logger, db: db}
}

func (r *ActionRepository) GetActions() (actions []models.Action, err error) {
	return actions, r.db.DB.Find(&actions).Error
}

func (r *ActionRepository) GetAction(id uint) (action *models.Action, err error) {
	return action, r.db.DB.First(&action, id).Error
}

func (r *ActionRepository) CreateAction(action models.Action) (models.Action, error) {
	return action, r.db.DB.Create(&action).Error
}

func (r *ActionRepository) UpdateAction(id uint, action models.Action) error {
	return r.db.DB.Model(&models.Action{}).Where("id = ?", id).Updates(action).Error
}

func (r *ActionRepository) DeleteAction(id uint) error {
	return r.db.DB.Delete(&models.Action{}, id).Error
}
