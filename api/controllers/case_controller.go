package controllers

import (
	"github.com/721945/dlaw-backend/api/dtos"
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CaseController struct {
	logger      *libs.Logger
	caseService services.CaseService
}

func NewCaseController(logger *libs.Logger, s services.CaseService) CaseController {
	return CaseController{logger: logger, caseService: s}
}

func (ctrl CaseController) GetCases(c *gin.Context) {
	cases, err := ctrl.caseService.GetCases()

	if err != nil {
		ctrl.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(200, cases)
}

func (ctrl CaseController) GetCase(c *gin.Context) {

	userId, isExist := c.Get("id")

	if !isExist {
		_ = c.Error(libs.ErrUnauthorized)
		return
	}

	paramId := c.Param("id")

	id, err := uuid.Parse(paramId)

	if err != nil {
		ctrl.logger.Error(err)
		_ = c.Error(err)
		return
	}

	mCase, err := ctrl.caseService.GetCase(id, userId.(uuid.UUID))

	if err != nil {
		ctrl.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(200, gin.H{
		"data": mCase,
	})
}

func (ctrl CaseController) CreateCase(c *gin.Context) {

	id, isExisted := c.Get("id")

	if !isExisted {
		_ = c.Error(libs.ErrUnauthorized)
		return
	}

	var input dtos.CreateCaseDto

	if err := c.ShouldBindJSON(&input); err != nil {
		ctrl.logger.Error(err)
		_ = c.Error(err)
		return
	}

	mCase, err := ctrl.caseService.CreateCase(input, id.(uuid.UUID))

	if err != nil {
		ctrl.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(200, gin.H{
		"data": mCase,
	})
}

func (ctrl CaseController) GetOwnCases(c *gin.Context) {

	id, isExisted := c.Get("id")

	ctrl.logger.Info(id)

	if !isExisted {

		_ = c.Error(libs.ErrUnauthorized)
		return
	}

	mCase, err := ctrl.caseService.GetOwnCases(id.(uuid.UUID))

	if err != nil {
		ctrl.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(200, gin.H{
		"data": mCase,
	})
}

func (ctrl CaseController) UpdateCase(c *gin.Context) {

	userId, isExisted := c.Get("id")

	if !isExisted {
		_ = c.Error(libs.ErrUnauthorized)
		return
	}

	var input dtos.UpdateCaseDto

	if err := c.ShouldBindJSON(&input); err != nil {
		ctrl.logger.Error(err)
		_ = c.Error(err)
		return
	}

	paramId := c.Param("id")

	id, err := uuid.Parse(paramId)

	if err != nil {
		ctrl.logger.Error(err)
		_ = c.Error(err)
		return
	}

	err = ctrl.caseService.UpdateCase(id, input, userId.(uuid.UUID))

	if err != nil {
		ctrl.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(200, gin.H{
		"data": "ok",
	})
}
