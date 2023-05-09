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

func (f FileController) GetFilesByTags(c *gin.Context) {
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

	userId, _ := c.Get("id")

	user := userId.(uuid.UUID)

	dto, err := f.fileService.GetFile(id, &user)

	if err != nil {
		f.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": dto})
}

func (f FileController) GetPublicFile(c *gin.Context) {
	paramId := c.Param("id")

	id, err := uuid.Parse(paramId)

	if err != nil {
		_ = c.Error(err)
		return
	}

	dto, err := f.fileService.GetFile(id, nil)

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
	paramId := c.Param("id")
	id, err := uuid.Parse(paramId)

	if err != nil {
		f.logger.Error(err)
		_ = c.Error(err)
		return
	}

	var input dtos.UpdateFileDto
	if err := c.ShouldBindJSON(&input); err != nil {
		f.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = f.fileService.UpdateFile(id, input)

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

	if formFolderId == "" {
		url, err := f.fileService.UploadFileNoFolder(
			file,
			fileHeader.Filename,
			fileHeader.Header.Get("Content-Type"),
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

		return
	}

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

func (f FileController) RecentViewedFiles(c *gin.Context) {
	id, isExist := c.Get("id")

	if !isExist {
		f.logger.Error("User not found")
		_ = c.Error(libs.ErrInternalServerError)
		return
	}

	dto, err := f.fileService.GetRecentFileOpened(id.(uuid.UUID))

	if err != nil {
		f.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": dto})
}

func (f FileController) MoveFile(c *gin.Context) {
	userId, isExist := c.Get("id")

	if !isExist {
		f.logger.Error("User not found")
		_ = c.Error(libs.ErrInternalServerError)
		return
	}

	id := c.Param("id")

	fileId, err := uuid.Parse(id)

	if err != nil {
		f.logger.Error(err)
		_ = c.Error(err)
		return
	}

	var input dtos.MoveFileDto

	if err := c.ShouldBindJSON(&input); err != nil {
		f.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = f.fileService.MoveFile(fileId, input, userId.(uuid.UUID))

	if err != nil {
		f.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": "ok"})
}

func (f FileController) SearchFiles(c *gin.Context) {
	userId, isExist := c.Get("id")

	if !isExist {
		f.logger.Error("User not found")
		_ = c.Error(libs.ErrInternalServerError)
		return
	}

	word := c.Param("word")

	folderID := c.DefaultQuery("folderId", "")
	tag := c.DefaultQuery("tag", "")
	caseID := c.DefaultQuery("caseId", "")
	fileType := c.DefaultQuery("type", "")
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "20")

	files, pagination, err := f.fileService.SearchFiles(word, caseID, folderID, tag, fileType, page, limit, userId.(uuid.UUID))

	if err != nil {
		f.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": files, "pagination": pagination})
}

func (f FileController) ShareFile(c *gin.Context) {
	userId, isExist := c.Get("id")

	if !isExist {
		f.logger.Error("User not found")
		_ = c.Error(libs.ErrInternalServerError)
		return
	}

	id := c.Param("id")

	link, err := f.fileService.ShareFile(id, userId.(uuid.UUID))

	if err != nil {
		f.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": link})
}

func (f FileController) PublicFile(c *gin.Context) {
	userId, isExist := c.Get("id")

	if !isExist {
		f.logger.Error("User not found")
		_ = c.Error(libs.ErrInternalServerError)
		return
	}

	id := c.Param("id")

	link, err := f.fileService.PublicFile(id, userId.(uuid.UUID))

	if err != nil {
		f.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": link})
}

func (f FileController) UnShareFile(c *gin.Context) {
	userId, isExist := c.Get("id")

	if !isExist {
		f.logger.Error("User not found")
		_ = c.Error(libs.ErrInternalServerError)
		return
	}

	id := c.Param("id")

	err := f.fileService.RemoveShareFile(id, userId.(uuid.UUID))

	if err != nil {
		f.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": "ok"})
}

func (f FileController) UnPublicFile(c *gin.Context) {
	userId, isExist := c.Get("id")

	if !isExist {
		f.logger.Error("User not found")
		_ = c.Error(libs.ErrInternalServerError)
		return
	}

	id := c.Param("id")

	link, err := f.fileService.PublicFile(id, userId.(uuid.UUID))

	if err != nil {
		f.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": link})
}
