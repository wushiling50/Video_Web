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

func VideoUpload(c *gin.Context) {
	var videoUpload service.VideoUploadService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))

	resource, err := c.FormFile("resource")
	if err != nil {
		logging.Error(err)
		c.JSON(400, e.ErrorResponse(err))
	}

	suffix := path.Ext(resource.Filename)
	if suffix != ".mp4" && suffix != ".avi" {
		logging.Error(err)
		c.JSON(400, e.ErrorResponse(err))
	}

	dir := "E:/VSCODE/Gocode/goproject/Video_Web/static/videos/"
	filename := strconv.FormatInt(time.Now().UnixNano(), 10) + claim.UserName + suffix
	dst := dir + filename
	c.SaveUploadedFile(resource, dst)

	if err := c.ShouldBind(&videoUpload); err != nil {
		logging.Error(err)
		c.JSON(400, e.ErrorResponse(err))
	} else {
		res := videoUpload.VideoUpload(claim.Id, dst)
		c.JSON(200, res)
	}
}

func CollectionCreate(c *gin.Context) {
	var collectionCreate service.CollectionCreateService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&collectionCreate); err != nil {
		logging.Error(err)
		c.JSON(400, e.ErrorResponse(err))
	} else {
		res := collectionCreate.CollectionCreate(claim.Id)
		c.JSON(200, res)
	}
}

func CollectionsShow(c *gin.Context) {
	var collectionsShow service.CollectionsShowService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&collectionsShow); err != nil {
		logging.Error(err)
		c.JSON(400, e.ErrorResponse(err))
	} else {
		res := collectionsShow.CollectionsShow(claim.Id)
		c.JSON(200, res)
	}
}

func Collect(c *gin.Context) {
	var collect service.CollectService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&collect); err != nil {
		logging.Error(err)
		c.JSON(400, e.ErrorResponse(err))
	} else {
		res := collect.Collect(claim.Id)
		c.JSON(200, res)
	}
}

func CollectShow(c *gin.Context) {
	var collectShow service.CollectShowService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&collectShow); err != nil {
		logging.Error(err)
		c.JSON(400, e.ErrorResponse(err))
	} else {
		res := collectShow.CollectShow(claim.Id)
		c.JSON(200, res)
	}
}

func Liked(c *gin.Context) {
	var liked service.LikedService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&liked); err != nil {
		logging.Error(err)
		c.JSON(400, e.ErrorResponse(err))
	} else {
		res := liked.Liked(claim.Id)
		c.JSON(200, res)
	}

}

func LikedShow(c *gin.Context) {
	var likedShow service.LikedShowService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&likedShow); err != nil {
		logging.Error(err)
		c.JSON(400, e.ErrorResponse(err))
	} else {
		res := likedShow.LikedShow(claim.Id)
		c.JSON(200, res)
	}

}

func Comment(c *gin.Context) {
	var comment service.CommentService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&comment); err != nil {
		logging.Error(err)
		c.JSON(400, e.ErrorResponse(err))
	} else {
		res := comment.Comment(claim.Id)
		c.JSON(200, res)
	}

}

func CommentShow(c *gin.Context) {
	var commentShow service.CommentShowService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&commentShow); err != nil {
		logging.Error(err)
		c.JSON(400, e.ErrorResponse(err))
	} else {
		res := commentShow.CommentShow(claim.Id)
		c.JSON(200, res)
	}
}

func Transmit(c *gin.Context) {
	var transmit service.TransmitService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&transmit); err != nil {
		logging.Error(err)
		c.JSON(400, e.ErrorResponse(err))
	} else {
		res := transmit.Transmit(claim.Id)
		c.JSON(200, res)
	}

}

func Danmu(c *gin.Context) {
	var danmu service.DanmuService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&danmu); err != nil {
		logging.Error(err)
		c.JSON(400, e.ErrorResponse(err))
	} else {
		res := danmu.Danmu(claim.Id)
		c.JSON(200, res)
	}

}

func DanmuShow(c *gin.Context) {
	var danmuShow service.DanmuShowService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&danmuShow); err != nil {
		logging.Error(err)
		c.JSON(400, e.ErrorResponse(err))
	} else {
		res := danmuShow.DanmuShow(claim.Id)
		c.JSON(200, res)
	}
}

func View(c *gin.Context) {
	var view service.ViewService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&view); err != nil {
		logging.Error(err)
		c.JSON(400, e.ErrorResponse(err))
	} else {
		res := view.View(claim.Id)
		c.JSON(200, res)
	}
}

func RankList(c *gin.Context) {
	var rankList service.RankListService
	claim, _ := utils.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&rankList); err != nil {
		logging.Error(err)
		c.JSON(400, e.ErrorResponse(err))
	} else {
		res := rankList.RankList(claim.Id)
		c.JSON(200, res)
	}
}
