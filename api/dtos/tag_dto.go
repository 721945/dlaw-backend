package dtos

import (
	"github.com/721945/dlaw-backend/models"
)

type TagDto struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreateTagDto struct {
	Name string `json:"name" binding:"required"`
}

type UpdateTagDto struct {
	Name string `json:"name" binding:"required"`
}

func ToTagDto(tag models.Tag) TagDto {
	return TagDto{
		ID:   tag.ID.String(),
		Name: tag.Name,
	}
}

func ToTagDtos(tags []models.Tag) []TagDto {
	var tagDtos = make([]TagDto, 0)

	for _, tag := range tags {
		tagDtos = append(tagDtos, ToTagDto(tag))
	}

	return tagDtos
}

func (dto CreateTagDto) ToModel() models.Tag {
	return models.Tag{
		Name: dto.Name,
	}
}

func (dto UpdateTagDto) ToModel() models.Tag {
	return models.Tag{
		Name: dto.Name,
	}
}
