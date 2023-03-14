package controllers

import (
	"github.com/721945/dlaw-backend/api/dtos"
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/services"
	"github.com/gin-gonic/gin"
	"strconv"
)

type PermissionController struct {
	logger            *libs.Logger
	permissionService services.PermissionService
}

func NewPermissionController(logger *libs.Logger, permissionService services.PermissionService) PermissionController {
	return PermissionController{logger: logger, permissionService: permissionService}
}

func (p PermissionController) GetPermissions(c *gin.Context) {
	permissions, err := p.permissionService.GetPermissions()
	if err != nil {
		p.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(200, gin.H{"data": permissions})
}

func (p PermissionController) GetPermission(c *gin.Context) {
	paramId := c.Param("id")

	id, err := strconv.Atoi(paramId)
	if err != nil {
		p.logger.Error(err)
		_ = c.Error(err)
		return
	}

	permission, err := p.permissionService.GetPermission(uint(id))

	if err != nil {
		p.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(200, gin.H{"data": dtos.ToPermissionDto(*permission)})
}

func (p PermissionController) CreatePermission(c *gin.Context) {

	var input dtos.CreatePermissionDto
	if err := c.ShouldBindJSON(&input); err != nil {
		p.logger.Error(err)
		_ = c.Error(err)
		return
	}

	permission, err := p.permissionService.CreatePermission(input.ToModel())

	if err != nil {
		p.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(201, gin.H{"data": dtos.ToPermissionDto(permission)})
}

func (p PermissionController) UpdatePermission(c *gin.Context) {
	paramId := c.Param("id")

	id, err := strconv.Atoi(paramId)

	if err != nil {
		p.logger.Error(err)
		_ = c.Error(err)
		return
	}

	var input dtos.UpdatePermissionDto
	if err := c.ShouldBindJSON(&input); err != nil {
		p.logger.Error(err)
		_ = c.Error(err)
		return
	}

	err = p.permissionService.UpdatePermission(uint(id), input.ToModel())

	if err != nil {
		p.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(200, gin.H{})

}

func (p PermissionController) DeletePermission(c *gin.Context) {
	paramId := c.Param("id")

	id, err := strconv.Atoi(paramId)

	if err != nil {
		p.logger.Error(err)
		_ = c.Error(err)
		return
	}

	err = p.permissionService.DeletePermission(uint(id))

	if err != nil {
		p.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(200, gin.H{})
}
