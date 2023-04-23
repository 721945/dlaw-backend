package controllers

import (
	"github.com/721945/dlaw-backend/api/dtos"
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type FolderController struct {
	logger        *libs.Logger
	folderService services.FolderService
}

func NewFolderController(logger *libs.Logger, s services.FolderService) FolderController {
	return FolderController{logger: logger, folderService: s}
}

func (t FolderController) GetFolders(c *gin.Context) {
	folders, err := t.folderService.GetFolders()

	if err != nil {
		t.logger.Error(err)
		_ = c.Error(err)
		return
	}

	tagsDto := dtos.ToFolderDtos(folders)

	c.JSON(200, gin.H{"data": tagsDto})
}

func (t FolderController) GetFolder(c *gin.Context) {
	paramId := c.Param("id")

	id, err := uuid.Parse(paramId)

	userId, isExist := c.Get("id")

	if !isExist {
		t.logger.Error(err)
		_ = c.Error(err)
		return
	}

	if err != nil {
		t.logger.Error(err)
		_ = c.Error(err)
		return
	}

	folder, err := t.folderService.GetFolder(id, userId.(uuid.UUID))

	if err != nil {
		t.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(200, gin.H{"data": folder})
}

func (t FolderController) CreateFolder(c *gin.Context) {
	var input dtos.CreateFolderDto
	//
	if err := c.ShouldBindJSON(&input); err != nil {
		t.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	folder, err := t.folderService.CreateFolder(input)
	//
	if err != nil {
		t.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": folder})
}

func (t FolderController) UpdateFolder(c *gin.Context) {
	var input dtos.UpdateFolderDto
	//
	var paramId = c.Param("id")
	//
	id, err := uuid.Parse(paramId)

	if err != nil {
		t.logger.Error(err)
		_ = c.Error(err)
		return
	}
	//
	if err := c.ShouldBindJSON(&input); err != nil {
		t.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, isExist := c.Get("id")

	if !isExist {
		t.logger.Error(err)
		_ = c.Error(err)
		return
	}

	err = t.folderService.UpdateFolder(id, input, userId.(uuid.UUID))

	if err != nil {
		t.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (t FolderController) DeleteFolder(c *gin.Context) {
	var paramId = c.Param("id")

	id, err := uuid.Parse(paramId)

	if err != nil {
		t.logger.Error(err)
		_ = c.Error(err)
		return
	}

	userId, isExist := c.Get("id")

	if !isExist {
		t.logger.Error(err)
		_ = c.Error(err)
		return
	}

	err = t.folderService.DeleteFolder(id, userId.(uuid.UUID))

	if err != nil {
		t.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": true})
}

func (t FolderController) ArchiveFolder(c *gin.Context) {
	var paramId = c.Param("id")

	id, err := uuid.Parse(paramId)

	if err != nil {
		t.logger.Error(err)
		_ = c.Error(err)
		return
	}

	userId, isExist := c.Get("id")

	if !isExist {
		t.logger.Error(err)
		_ = c.Error(err)
		return
	}

	err = t.folderService.DeleteFolder(id, userId.(uuid.UUID))

	if err != nil {
		t.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": true})
}

// GetFolderLog
func (t FolderController) GetFolderLog(c *gin.Context) {
	userId, isExist := c.Get("id")

	if !isExist {
		t.logger.Error("User not found")
		_ = c.Error(libs.ErrUnauthorized)
		return
	}

	var paramId = c.Param("id")

	id, err := uuid.Parse(paramId)

	if err != nil {
		t.logger.Error(err)
		_ = c.Error(err)
		return
	}

	folderLog, err := t.folderService.GetFolderLog(id)

	if err != nil {
		t.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": folderLog})
}
