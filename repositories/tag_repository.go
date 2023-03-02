package repositories

import (
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
)

type TagRepository struct {
	logger *libs.Logger
	db     libs.Database
}

func NewTagRepository(logger *libs.Logger, db libs.Database) TagRepository {
	return TagRepository{logger: logger, db: db}
}

func (r *TagRepository) GetTags() (tags []models.Tag, err error) {
	return tags, r.db.DB.Find(&tags).Error
}

func (r *TagRepository) GetTag(id uint) (tag *models.Tag, err error) {
	return tag, r.db.DB.First(&tag, id).Error
}

func (r *TagRepository) CreateTag(tag models.Tag) (models.Tag, error) {
	return tag, r.db.DB.Create(&tag).Error
}

func (r *TagRepository) UpdateTag(id uint, tag models.Tag) error {
	return r.db.DB.Model(&models.Tag{}).Where("id = ?", id).Updates(tag).Error
}

func (r *TagRepository) DeleteTag(id uint) error {
	return r.db.DB.Delete(&models.Tag{}, id).Error
}
