package controllers

import (
	"github.com/721945/dlaw-backend/api/dtos"
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	id, err := uuid.Parse(paramId)
	if err != nil {
		p.logger.Error(err)
		_ = c.Error(err)
		return
	}

	permission, err := p.permissionService.GetPermission(id)

	if err != nil {
		p.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(200, gin.H{"data": permission})
}
func (p PermissionController) GetPermissionName(c *gin.Context) {
	name := c.Param("name")

	permission, err := p.permissionService.GetPermissionByName(name)

	if err != nil {
		p.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(200, gin.H{"data": permission})
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

	c.JSON(201, gin.H{"data": permission})
}

func (p PermissionController) UpdatePermission(c *gin.Context) {
	paramId := c.Param("id")

	id, err := uuid.Parse(paramId)

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

	err = p.permissionService.UpdatePermission(id, input.ToModel())

	if err != nil {
		p.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(200, gin.H{})

}

func (p PermissionController) DeletePermission(c *gin.Context) {
	paramId := c.Param("id")

	id, err := uuid.Parse(paramId)

	if err != nil {
		p.logger.Error(err)
		_ = c.Error(err)
		return
	}

	err = p.permissionService.DeletePermission(id)

	if err != nil {
		p.logger.Error(err)
		_ = c.Error(err)
		return
	}

	c.JSON(200, gin.H{})
}
