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
	return tags, r.db.DB.Table("tags").
		Select("tags.id, tags.name, COALESCE(file_counts.count, 0) as count").
		Joins("LEFT JOIN (SELECT file_tags.tag_id, COUNT(DISTINCT file_tags.file_id) as count FROM file_tags WHERE file_tags.file_id IN (?) GROUP BY file_tags.tag_id) as file_counts ON tags.id = file_counts.tag_id", fileIds).
		Group("tags.id").
		Group("tags.name").
		Group("file_counts.count").
		Where("tags.show_menu = TRUE").
		Find(&tags).Error
	//Select("id, COALESCE(COUNT(file_tags.file_id), 0) as count, name").
	//Joins("LEFT JOIN file_tags ON tags.id = file_tags.tag_id").
	//Group("id").Group("name").
	//Where("tags.show_menu = TRUE").
	//Where("file_tags.file_id IN (?) OR file_tags.file_id IS NULL", fileIds).
	//Find(&tags).Error
}

func (r *TagRepository) CountFilesInTagsByFolderId(folderId uuid.UUID) (tags []models.TagCount, err error) {
	return tags, r.db.DB.Table("tags").
		Select("tags.id, tags.name, COALESCE(file_counts.count, 0) as count").
		Joins("LEFT JOIN (SELECT file_tags.tag_id, COUNT(DISTINCT file_tags.file_id) as count FROM file_tags JOIN files ON file_tags.file_id = files.id WHERE files.folder_id = ? GROUP BY file_tags.tag_id) as file_counts ON tags.id = file_counts.tag_id", folderId).
		Where("tags.show_menu = TRUE").
		Find(&tags).Error
	//Select("tags.id, tags.name, COUNT(DISTINCT file_tags.file_id) as count").
	//Joins("LEFT JOIN file_tags ON tags.id = file_tags.tag_id").
	//Joins("LEFT JOIN files ON file_tags.file_id = files.id").
	//Where("tags.show_menu = TRUE").
	//Where("files.folder_id = ?", folderId).
	//Group("tags.id").
	//Group("tags.name").
	//Find(&tags).Error
	//Select("id, COALESCE(COUNT(file_tags.file_id), 0) as count, name").
	//Joins("LEFT JOIN file_tags ON tags.id = file_tags.tag_id").
	//Joins("LEFT JOIN files ON file_tags.file_id = files.id").
	//Where("tags.show_menu = TRUE").
	//Where("file_tags.file_id IN (SELECT id FROM files WHERE folder_id = ?) OR file_tags.file_id IS NULL", folderId).
	//Group("id").Group("name").
	//Find(&tags).Error
}
