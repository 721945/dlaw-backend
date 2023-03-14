package controllers

import (
	"github.com/721945/dlaw-backend/api/dtos"
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/services"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"strconv"
)

type FileController struct {
	logger      *libs.Logger
	fileService services.FileService
}

func NewFileController(logger *libs.Logger, s services.FileService) FileController {
	return FileController{logger: logger, fileService: s}
}

func (f FileController) GetFiles(c *gin.Context) {
	actions, err := f.fileService.GetFiles()
	if err != nil {
		f.logger.Error(err)
		_ = c.Error(err)
		return
	}

	actionsDto := make([]dtos.FileDto, len(actions))

	for i, action := range actions {
		actionsDto[i] = dtos.ToFileDto(action)
	}

	c.JSON(http.StatusOK, gin.H{"data": actionsDto})

}

func (f FileController) GetFile(c *gin.Context) {
	paramId := c.Param("id")

	id, err := strconv.Atoi(paramId)

	if err != nil {
		f.logger.Error(err)
		_ = c.Error(err)
		return
	}

	action, err := f.fileService.GetFile(uint(id))

	if err != nil {
		f.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": dtos.ToFileDto(*action)})
}

func (f FileController) CreateFile(c *gin.Context) {
	var input dtos.CreateFileDto
	if err := c.ShouldBindJSON(&input); err != nil {
		f.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	action, err := f.fileService.CreateFile(input.ToModel())
	if err != nil {
		f.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": dtos.ToFileDto(action)})
}

func (f FileController) UpdateFile(c *gin.Context) {
	var input dtos.UpdateFileDto
	if err := c.ShouldBindJSON(&input); err != nil {
		f.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	paramId := c.Param("id")
	id, err := strconv.Atoi(paramId)

	if err != nil {
		f.logger.Error(err)
		_ = c.Error(err)
		return
	}

	err = f.fileService.UpdateFile(uint(id), input.ToModel())

	if err != nil {
		f.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (f FileController) DeleteFile(c *gin.Context) {
	paramId := c.Param("id")
	id, err := strconv.Atoi(paramId)

	if err != nil {
		f.logger.Error(err)
		_ = c.Error(err)
		return
	}

	err = f.fileService.DeleteFile(uint(id))
	if err != nil {
		f.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (f FileController) GetSignedUrl(c *gin.Context) {
	var input dtos.GetSingleFileDto
	if err := c.ShouldBindJSON(&input); err != nil {
		f.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	url, err := f.fileService.GetSignedUrl(input.Amount)

	if err != nil {
		f.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": url})
}

func (f FileController) UploadFile(c *gin.Context) {
	//var input dtos.UploadFileDto
	//if err := c.ShouldBindJSON(&input); err != nil {
	//	f.logger.Error(err)
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	//}

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

	url, err := f.fileService.UploadFile(file, fileHeader.Filename)

	if err != nil {
		f.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": url,
	})
}
