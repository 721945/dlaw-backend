package services

import (
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
	"github.com/721945/dlaw-backend/repositories"
)

type TagService struct {
	logger  *libs.Logger
	tagRepo repositories.TagRepository
}

func NewTagService(logger *libs.Logger, r repositories.TagRepository) TagService {
	return TagService{logger: logger, tagRepo: r}
}

func (s *TagService) GetTags() (tags []models.Tag, err error) {
	return s.tagRepo.GetTags()
}

func (s *TagService) GetTag(id uint) (tag *models.Tag, err error) {
	return s.tagRepo.GetTag(id)
}

func (s *TagService) CreateTag(tag models.Tag) (models.Tag, error) {
	return s.tagRepo.CreateTag(tag)
}

func (s *TagService) UpdateTag(id uint, tag models.Tag) error {
	return s.tagRepo.UpdateTag(id, tag)
}

func (s *TagService) DeleteTag(id uint) error {
	return s.tagRepo.DeleteTag(id)
}

//func (s *TagService) GetTagByName(name string) (tag *models.Tag, err error) {
//	return s.tagRepo.GetTagByName(name)
//}
