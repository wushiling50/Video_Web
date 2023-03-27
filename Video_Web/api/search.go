package api

import (
	"main/Video_Web/pkg/e"
	"main/Video_Web/pkg/utils"
	"main/Video_Web/service"

	"github.com/gin-gonic/gin"
	logging "github.com/sirupsen/logrus"
)

func SearchUser(c *gin.Context) {
	var searchUser service.SearchUserService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&searchUser); err != nil {
		logging.Error(err)
		c.JSON(400, e.ErrorResponse(err))
	} else {
		res := searchUser.SearchUser(claim.Id)
		c.JSON(200, res)
	}
}

func SearchVideo(c *gin.Context) {
	var searchVideo service.SearchVideoService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&searchVideo); err != nil {
		logging.Error(err)
		c.JSON(400, e.ErrorResponse(err))
	} else {
		res := searchVideo.SearchVideo(claim.Id)
		c.JSON(200, res)
	}
}
