package controllers

import (
	"github.com/721945/dlaw-backend/api/dtos"
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/services"
	"github.com/gin-gonic/gin"
	"strconv"
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

	paramId := c.Param("id")

	id, err := strconv.Atoi(paramId)

	if err != nil {
		ctrl.logger.Error(err)
		_ = c.Error(err)
		return
	}

	mCase, err := ctrl.caseService.GetCase(uint(id))

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

	var input dtos.CreateCaseDto

	if err := c.ShouldBindJSON(&input); err != nil {
		ctrl.logger.Error(err)
		_ = c.Error(err)
		return
	}

	mCase, err := ctrl.caseService.CreateCase(input)

	if err != nil {
		ctrl.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(200, gin.H{
		"data": mCase,
	})
}
