package controllers

import (
	"github.com/721945/dlaw-backend/api/dtos"
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ActionController struct {
	logger        *libs.Logger
	actionService services.ActionService
}

func NewActionController(logger *libs.Logger, s services.ActionService) ActionController {
	return ActionController{logger: logger, actionService: s}
}

func (a ActionController) GetActions(ctx *gin.Context) {
	actions, err := a.actionService.GetActions()
	if err != nil {
		a.logger.Error(err)
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": actions})

}

func (a ActionController) GetAction(c *gin.Context) {
	paramId := c.Param("id")

	id, err := strconv.Atoi(paramId)

	if err != nil {
		a.logger.Error(err)
		_ = c.Error(err)
		return
	}

	action, err := a.actionService.GetAction(uint(id))
	if err != nil {
		a.logger.Error(err)
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": action})
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

	c.JSON(http.StatusOK, gin.H{"data": action})
}


