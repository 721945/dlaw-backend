package controllers

import (
	"github.com/721945/dlaw-backend/api/dtos"
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type TypeController struct {
	logger      *libs.Logger
	fileService services.FileTypeService
}

func NewTypeController(logger *libs.Logger, s services.FileTypeService) TypeController {
	return TypeController{logger: logger, fileService: s}
}

func (t TypeController) GetTypes(c *gin.Context) {
	tags, err := t.fileService.GetFileTypes()

	if err != nil {
		t.logger.Error(err)
		_ = c.Error(err)
		return
	}

	//tagsDto := dtos.ToTypeDtos(tags)

	c.JSON(200, gin.H{"data": tags})
}

func (t TypeController) GetType(c *gin.Context) {
	paramId := c.Param("id")

	id, err := uuid.Parse(paramId)

	if err != nil {
		t.logger.Error(err)
		_ = c.Error(err)
		return
	}

	tag, err := t.fileService.GetFileType(id)

	if err != nil {
		t.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(200, gin.H{"data": tag})
}

func (t TypeController) CreateType(c *gin.Context) {
	var input dtos.CreateFileTypeDto

	if err := c.ShouldBindJSON(&input); err != nil {
		t.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileType, err := t.fileService.CreateFileType(input)

	if err != nil {
		t.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": fileType})
}
func (t TypeController) CreateTypes(c *gin.Context) {
	var input dtos.CreateFileTypesDto

	if err := c.ShouldBindJSON(&input); err != nil {
		t.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ids := make([]string, len(input.NameList))

	for i, name := range input.NameList {
		id, err := t.fileService.CreateFileType(dtos.CreateFileTypeDto{Name: name})

		if err != nil {
			t.logger.Error(err)
			_ = c.Error(err)
			return
		}

		ids[i] = id
	}

	c.JSON(http.StatusCreated, gin.H{"data": ids})
}

func (t TypeController) UpdateType(c *gin.Context) {
	var input dtos.UpdateFileTypeDto

	var paramId = c.Param("id")

	id, err := uuid.Parse(paramId)

	if err := c.ShouldBindJSON(&input); err != nil {
		t.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = t.fileService.UpdateFileType(id, input)

	if err != nil {
		t.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (t TypeController) DeleteType(c *gin.Context) {
	var paramId = c.Param("id")

	id, err := uuid.Parse(paramId)

	if err != nil {
		t.logger.Error(err)
		_ = c.Error(err)
		return
	}

	err = t.fileService.DeleteFileType(id)

	if err != nil {
		t.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": true})
}
