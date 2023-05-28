package dtos

import (
	"github.com/721945/dlaw-backend/models"
)

type TagDto struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

type TagCountDto struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Count       int    `json:"count"`
}

type CreateTagDto struct {
	Name        string `json:"name" binding:"required"`
	DisplayName string `json:"displayName" binding:"required"`
}

type UpdateTagDto struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

func ToTagCountDto(tag models.TagCount) TagCountDto {
	return TagCountDto{
		ID:          tag.ID.String(),
		Name:        tag.Name,
		DisplayName: tag.DisplayName,
		Count:       tag.Count,
	}
}

func ToTagDto(tag models.Tag) TagDto {
	return TagDto{
		ID:          tag.ID.String(),
		Name:        tag.Name,
		DisplayName: tag.DisplayName,
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
		Name:        dto.Name,
		DisplayName: dto.DisplayName,
	}
}

func (dto UpdateTagDto) ToModel() models.Tag {
	return models.Tag{
		Name:        dto.Name,
		DisplayName: dto.DisplayName,
	}
}
