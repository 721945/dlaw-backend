package controllers

import (
	"github.com/721945/dlaw-backend/api/dtos"
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type ActionController struct {
	logger        *libs.Logger
	actionService services.ActionService
}

func NewActionController(logger *libs.Logger, s services.ActionService) ActionController {
	return ActionController{logger: logger, actionService: s}
}

// GetActions godoc
// @summary Get all actions
// @description Get all actions in database
// @tags actions
// @security ApiKeyAuth
// @id GetActions
// @accept json
// @produce json
// @response 200 {string} string "OK"
// @response 400 {string}  "Bad Request"
// @response 401 {string}  "Unauthorized"
// @response 500 {string}  "Internal Server Error"
// @Router /actions [get]
func (a ActionController) GetActions(ctx *gin.Context) {
	actions, err := a.actionService.GetActions()
	if err != nil {
		a.logger.Error(err)
		_ = ctx.Error(err)
		return
	}

	actionsDto := make([]dtos.ActionDto, len(actions))

	for i, action := range actions {
		actionsDto[i] = dtos.ToActionDto(action)
	}

	ctx.JSON(http.StatusOK, gin.H{"data": actionsDto})

}

func (a ActionController) GetAction(c *gin.Context) {
	paramId := c.Param("id")

	id, err := uuid.Parse(paramId)

	if err != nil {
		a.logger.Error(err)
		_ = c.Error(err)
		return
	}

	action, err := a.actionService.GetAction(id)

	if err != nil {
		a.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": dtos.ToActionDto(*action)})
}

func (a ActionController) CreateAction(c *gin.Context) {
	var input dtos.CreateActionDto
	if err := c.ShouldBindJSON(&input); err != nil {
		a.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	action, err := a.actionService.CreateAction(input.ToModel())
	if err != nil {
		a.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": dtos.ToActionDto(action)})
}

func (a ActionController) UpdateAction(c *gin.Context) {
	var input dtos.UpdateActionDto
	if err := c.ShouldBindJSON(&input); err != nil {
		a.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	paramId := c.Param("id")
	id, err := uuid.Parse(paramId)

	if err != nil {
		a.logger.Error(err)
		_ = c.Error(err)
		return
	}

	err = a.actionService.UpdateAction(id, input.ToModel())

	if err != nil {
		a.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (a ActionController) DeleteAction(c *gin.Context) {
	paramId := c.Param("id")
	id, err := uuid.Parse(paramId)

	if err != nil {
		a.logger.Error(err)
		_ = c.Error(err)
		return
	}

	err = a.actionService.DeleteAction(id)
	if err != nil {
		a.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
