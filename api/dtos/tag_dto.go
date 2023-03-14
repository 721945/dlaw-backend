package dtos

import (
	"github.com/721945/dlaw-backend/models"
	"time"
)

type TagDto struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateTagDto struct {
	Name string `json:"name" binding:"required"`
}

type UpdateTagDto struct {
	Name string `json:"name" binding:"required"`
}

func ToTagDto(tag models.Tag) TagDto {
	return TagDto{
		ID:        tag.ID,
		Name:      tag.Name,
		CreatedAt: tag.CreatedAt,
		UpdatedAt: tag.UpdatedAt,
	}
}

func ToTagDtos(tags []models.Tag) []TagDto {
	var tagDtos []TagDto

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
