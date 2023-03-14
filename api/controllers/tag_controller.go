package controllers

import (
	"github.com/721945/dlaw-backend/api/dtos"
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type TagController struct {
	logger     *libs.Logger
	tagService services.TagService
}

func NewTagController(logger *libs.Logger, s services.TagService) TagController {
	return TagController{logger: logger, tagService: s}
}

func (t TagController) GetTags(c *gin.Context) {
	tags, err := t.tagService.GetTags()

	if err != nil {
		t.logger.Error(err)
		_ = c.Error(err)
		return
	}

	tagsDto := dtos.ToTagDtos(tags)

	c.JSON(200, gin.H{"data": tagsDto})
}

func (t TagController) GetTag(c *gin.Context) {
	paramId := c.Param("id")

	id, err := strconv.Atoi(paramId)

	if err != nil {
		t.logger.Error(err)
		_ = c.Error(err)
		return
	}

	tag, err := t.tagService.GetTag(uint(id))

	if err != nil {
		t.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(200, gin.H{"data": tag})
}

func (t TagController) CreateTag(c *gin.Context) {
	var input dtos.CreateTagDto

	if err := c.ShouldBindJSON(&input); err != nil {
		t.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag, err := t.tagService.CreateTag(input.ToModel())

	if err != nil {
		t.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": dtos.ToTagDto(tag)})
}

func (t TagController) UpdateTag(c *gin.Context) {
	var input dtos.UpdateTagDto

	var paramId = c.Param("id")

	id, err := strconv.Atoi(paramId)

	if err := c.ShouldBindJSON(&input); err != nil {
		t.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = t.tagService.UpdateTag(uint(id), input.ToModel())

	if err != nil {
		t.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (t TagController) DeleteTag(c *gin.Context) {
	var paramId = c.Param("id")

	id, err := strconv.Atoi(paramId)

	if err != nil {
		t.logger.Error(err)
		_ = c.Error(err)
		return
	}

	err = t.tagService.DeleteTag(uint(id))

	if err != nil {
		t.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": true})
}
