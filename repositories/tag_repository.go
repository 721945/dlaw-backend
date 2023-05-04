package repositories

import (
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
	"github.com/google/uuid"
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

func (r *TagRepository) GetShowMenuTags() (tags []models.Tag, err error) {
	return tags, r.db.DB.Where("show_menu = true").Find(&tags).Error
}

func (r *TagRepository) GetTag(id uuid.UUID) (tag *models.Tag, err error) {
	return tag, r.db.DB.First(&tag, id).Error
}

func (r *TagRepository) CreateTag(tag models.Tag) (models.Tag, error) {
	return tag, r.db.DB.Create(&tag).Error
}

func (r *TagRepository) UpdateTag(id uuid.UUID, tag models.Tag) error {
	return r.db.DB.Model(&models.Tag{}).Where("id = ?", id).Updates(tag).Error
}

func (r *TagRepository) DeleteTag(id uuid.UUID) error {
	return r.db.DB.Delete(&models.Tag{}, id).Error
}

func (r *TagRepository) GetTagByNames(names []string) (tags []models.Tag, err error) {
	return tags, r.db.DB.Where("name IN (?)", names).Find(&tags).Error
}

func (r *TagRepository) GetTagByName(name string) (tag *models.Tag, err error) {
	return tag, r.db.DB.Where("name = ?", name).First(&tag).Error
}

func (r *TagRepository) CountFilesInTags(fileIds []uuid.UUID) (tags []models.TagCount, err error) {
	//return tags, r.db.DB.Table("file_tags").Select("tag_id as id, COUNT(file_id) as count, tags.name as name").Joins("LEFT JOIN tags ON tags.id = file_tags.tag_id").Group("tag_id").Group("name").Where("file_id IN (?) AND tags.show_menu = TRUE", fileIds).Find(&tags).Error
	return tags, r.db.DB.Table("tags").Select("id, COALESCE(COUNT(file_tags.file_id), 0) as count, name").Joins("LEFT JOIN file_tags ON tags.id = file_tags.tag_id").Group("id").Group("name").Where("tags.show_menu = TRUE").Where("file_tags.file_id IN (?) OR file_tags.file_id IS NULL", fileIds).Find(&tags).Error
	//return tags, r.db.DB.Where("show_menu = true").Find(&tags).Error
	//return tags, r.db.DB.Where("show_menu = true").Find(&tags).Error
	//return count, r.db.DB.Table("file_tags").Select("tag_id, COUNT(file_id) as count").Where("file_id IN (?)", fileIds).Group("tag_id").Find(&count).Error
	//return count, r.db.DB.Table("file_tags").Select("tag_id, COUNT(file_id) as count, show_menu").Where("file_id IN (?) AND show_menu = TRUE", fileIds).Group("tag_id").Find(&count).Error
}
