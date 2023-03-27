package api

import (
	"main/Video_Web/pkg/e"
	"main/Video_Web/pkg/utils"
	"main/Video_Web/service"
	"path"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	logging "github.com/sirupsen/logrus"
)

func UserRegister(c *gin.Context) {
	var userRegister service.UserService
	if err := c.ShouldBind(&userRegister); err == nil {
		res := userRegister.Register()
		c.JSON(200, res)
	} else {
		logging.Error(err)
		c.JSON(400, e.ErrorResponse(err))
	}
}

func UserLogin(c *gin.Context) {
	var userLogin service.UserService
	if err := c.ShouldBind(&userLogin); err == nil {
		res := userLogin.Login()
		c.JSON(200, res)
	} else {
		logging.Error(err)
		c.JSON(400, e.ErrorResponse(err))
	}
}

func UserChange(c *gin.Context) {
	var userChange service.UserChangeService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&userChange); err != nil {
		logging.Error(err)
		c.JSON(400, e.ErrorResponse(err))
	} else {
		res := userChange.UserChange(claim.Id)
		c.JSON(200, res)
	}

}

func ShowMsg(c *gin.Context) {
	var showMsg service.UserShowService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&showMsg); err != nil {
		logging.Error(err)
		c.JSON(400, e.ErrorResponse(err))
	} else {
		res := showMsg.ShowMsg(claim.Id)
		c.JSON(200, res)
	}

}

func ChangeMsg(c *gin.Context) {
	var msgchange service.UserMsgService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&msgchange); err != nil {
		logging.Error(err)
		c.JSON(400, e.ErrorResponse(err))
	} else {
		res := msgchange.MsgChange(claim.Id)
		c.JSON(200, res)
	}
}

func UploadImg(c *gin.Context) {
	var uploadImg service.UserImgUploadService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))

	img, err := c.FormFile("img")
	if err != nil {
		logging.Error(err)
		c.JSON(400, e.ErrorResponse(err))
	}

	suffix := path.Ext(img.Filename)
	if suffix != ".jpg" && suffix != ".jpeg" && suffix != ".png" {
		logging.Error(err)
		c.JSON(400, e.ErrorResponse(err))
	}

	dir := "E:/VSCODE/Gocode/goproject/Video_Web/static/imgs/"
	filename := strconv.FormatInt(time.Now().UnixNano(), 10) + claim.UserName + suffix
	dst := dir + filename
	c.SaveUploadedFile(img, dst)

	if err := c.ShouldBind(&uploadImg); err != nil {
		logging.Error(err)
		c.JSON(400, e.ErrorResponse(err))
	} else {
		res := uploadImg.UploadImg(claim.Id, dst)
		c.JSON(200, res)
	}
}

func SendingEmail(c *gin.Context) {
	var sendingEmail service.SendingEmailService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&sendingEmail); err != nil {
		logging.Error(err)
		c.JSON(400, e.ErrorResponse(err))
	} else {
		res := sendingEmail.SendingEmail(claim.Id)
		c.JSON(200, res)
	}

}

func HandlingEmail(c *gin.Context) {
	var handlingEmail service.HandlingEmailService

	if err := c.ShouldBind(&handlingEmail); err != nil {
		logging.Error(err)
		c.JSON(400, e.ErrorResponse(err))
	} else {
		res := handlingEmail.HandlingEmail(c.GetHeader("Authorization"))
		c.JSON(200, res)
	}

}
