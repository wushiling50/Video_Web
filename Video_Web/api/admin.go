package api

import (
	"main/Video_Web/pkg/e"
	"main/Video_Web/pkg/utils"
	"main/Video_Web/service"

	"github.com/gin-gonic/gin"
	logging "github.com/sirupsen/logrus"
)

func AdminRegister(c *gin.Context) {
	var adminRegister service.AdminService
	if err := c.ShouldBind(&adminRegister); err != nil {
		logging.Error(err)
		c.JSON(400, e.ErrorResponse(err))
	} else {
		res := adminRegister.AdminRegister()
		c.JSON(200, res)
	}
}

func AdminLogin(c *gin.Context) {
	var adminLogin service.AdminService
	if err := c.ShouldBind(&adminLogin); err == nil {
		res := adminLogin.AdminLogin()
		c.JSON(200, res)
	} else {
		logging.Error(err)
		c.JSON(400, e.ErrorResponse(err))
	}
}

func Check(c *gin.Context) {
	var check service.CheckService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&check); err == nil {
		res := check.Check(claim.Id)
		c.JSON(200, res)
	} else {
		logging.Error(err)
		c.JSON(400, e.ErrorResponse(err))
	}
}

func BlackList(c *gin.Context) {
	var blackList service.BlackListService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&blackList); err == nil {
		res := blackList.BlackList(claim.Id)
		c.JSON(200, res)
	} else {
		logging.Error(err)
		c.JSON(400, e.ErrorResponse(err))
	}
}

func Delete(c *gin.Context) {
	var delete service.DeleteService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&delete); err == nil {
		res := delete.Delete(claim.Id)
		c.JSON(200, res)
	} else {
		logging.Error(err)
		c.JSON(400, e.ErrorResponse(err))
	}
}
