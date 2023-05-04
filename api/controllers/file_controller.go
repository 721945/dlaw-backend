package controllers

import (
	"github.com/721945/dlaw-backend/api/dtos"
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"mime/multipart"
	"net/http"
)

type FileController struct {
	logger      *libs.Logger
	fileService services.FileService
}

func NewFileController(logger *libs.Logger, s services.FileService) FileController {
	return FileController{logger: logger, fileService: s}
}

func (f FileController) GetFiles(c *gin.Context) {
	dto, err := f.fileService.GetFiles()
	if err != nil {
		f.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": dto})
}

func (f FileController) GetFile(c *gin.Context) {
	paramId := c.Param("id")

	id, err := uuid.Parse(paramId)

	if err != nil {
		f.logger.Error(err)
		_ = c.Error(err)
		return
	}

	dto, err := f.fileService.GetFile(id)

	if err != nil {
		f.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": dto})
}

func (f FileController) CreateFile(c *gin.Context) {
	var input dtos.CreateFileDto
	if err := c.ShouldBindJSON(&input); err != nil {
		f.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dto, err := f.fileService.CreateFile(input.ToModel())

	if err != nil {
		f.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": dto})
}

func (f FileController) UpdateFile(c *gin.Context) {
	var input dtos.UpdateFileDto
	if err := c.ShouldBindJSON(&input); err != nil {
		f.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	paramId := c.Param("id")
	id, err := uuid.Parse(paramId)

	if err != nil {
		f.logger.Error(err)
		_ = c.Error(err)
		return
	}

	err = f.fileService.UpdateFile(id, input.ToModel())

	if err != nil {
		f.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": "ok",
	})
}

func (f FileController) DeleteFile(c *gin.Context) {
	paramId := c.Param("id")
	id, err := uuid.Parse(paramId)

	if err != nil {
		f.logger.Error(err)
		_ = c.Error(err)
		return
	}

	err = f.fileService.DeleteFile(id)
	if err != nil {
		f.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": "ok",
	})
}

func (f FileController) GetSignedUrl(c *gin.Context) {
	var input dtos.GetSingleFileDto
	if err := c.ShouldBindJSON(&input); err != nil {
		f.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//url, err := f.fileService.GetSignedUrl(input.Amount)

	//if err != nil {
	//	f.logger.Error(err)
	//	_ = c.Error(err)
	//	return
	//}

	c.JSON(http.StatusOK, gin.H{"data": []string{}})
}

func (f FileController) UploadFile(c *gin.Context) {

	userId, isExisted := c.Get("id")

	if !isExisted {
		f.logger.Error("User not found")
		_ = c.Error(libs.ErrInternalServerError)
		return
	}

	file, fileHeader, err := c.Request.FormFile("file")

	if err != nil {
		f.logger.Error(err)
		_ = c.Error(libs.ErrInternalServerError)
		return
	}

	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			f.logger.Error(err)
		}
	}(file)

	formFolderId := c.Request.FormValue("folderId")

	folderId, err := uuid.Parse(formFolderId)

	if err != nil {
		f.logger.Error(err)
		_ = c.Error(err)
		return
	}

	url, err := f.fileService.UploadFile(
		file,
		fileHeader.Filename,
		fileHeader.Header.Get("Content-Type"),
		folderId,
		userId.(uuid.UUID),
	)

	if err != nil {
		f.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": url,
	})
}

func (f FileController) CountFileInTags(c *gin.Context) {
	id, isExist := c.Get("id")

	if !isExist {
		f.logger.Error("User not found")
		_ = c.Error(libs.ErrInternalServerError)
		return
	}

	dto, err := f.fileService.CountFilesInTags(id.(uuid.UUID))

	if err != nil {
		f.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": dto})
}
