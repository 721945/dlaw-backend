package services

import (
	"github.com/721945/dlaw-backend/api/dtos"
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
	"github.com/721945/dlaw-backend/repositories"
	"github.com/google/uuid"
)

type TagService struct {
	logger  *libs.Logger
	tagRepo repositories.TagRepository
}

func NewTagService(logger *libs.Logger, r repositories.TagRepository) TagService {
	return TagService{logger: logger, tagRepo: r}
}

func (s *TagService) GetTags() ([]dtos.TagDto, error) {
	tags, err := s.tagRepo.GetTags()

	if err != nil {
		return []dtos.TagDto{}, err
	}

	return dtos.ToTagDtos(tags), nil
}

func (s *TagService) GetTagMenus() ([]dtos.TagDto, error) {
	tags, err := s.tagRepo.GetShowMenuTags()

	if err != nil {
		return []dtos.TagDto{}, err
	}

	return dtos.ToTagDtos(tags), nil
}

func (s *TagService) GetTag(id uuid.UUID) (*dtos.TagDto, error) {

	tagModel, err := s.tagRepo.GetTag(id)

	if err != nil {
		return nil, err
	}
	tag := dtos.ToTagDto(*tagModel)
	return &tag, nil
}

func (s *TagService) CreateTag(tag models.Tag) (string, error) {
	tag, err := s.tagRepo.CreateTag(tag)

	if err != nil {
		return "", err
	}

	return tag.ID.String(), err
}

func (s *TagService) UpdateTag(id uuid.UUID, tag models.Tag) error {
	return s.tagRepo.UpdateTag(id, tag)
}

func (s *TagService) DeleteTag(id uuid.UUID) error {
	return s.tagRepo.DeleteTag(id)
}

//func (s *TagService) GetTagByName(name string) (tag *models.Tag, err error) {
//	return s.fileTypeRepository.GetTagByName(name)
//}
