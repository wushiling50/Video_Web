package service

import (
	"context"
	"errors"
	"main/Video_Web/cache"
	"main/Video_Web/model"
	"main/Video_Web/pkg/e"
	"main/Video_Web/serializer"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type VideoUploadService struct {
	// Resource string `form:"resource" json:"resource"`
	Title string `form:"title" json:"title"` //视频标题
	Desc  string `form:"desc" json:"desc"`   //视频简介
}

type CollectionCreateService struct {
	Name string `form:"name" json:"name"` //收藏夹名称
	Open uint   `form:"open" json:"open"` //是否公开
}

type CollectionsShowService struct {
	PageNum  int `form:"page_num" json:"page_num" `
	PageSize int `form:"page_size" json:"page_size"`
}
type CollectService struct {
	Vid   string `form:"vid" json:"vid"`     //视频的ID
	Cid   uint   `form:"cid" json:"cid"`     //收藏夹ID
	State uint   `form:"state" json:"state"` //0为添加 ， 1为删除
}
type CollectShowService struct {
	Cid      uint `form:"cid" json:"cid"` //收藏夹ID
	PageNum  int  `form:"page_num" json:"page_num" `
	PageSize int  `form:"page_size" json:"page_size"`
}

type LikedService struct {
	Vid     string `form:"vid" json:"vid"`
	IsLiked uint   `form:"is_liked" json:"is_liked"` //0为不点赞，1为点赞
}

type LikedShowService struct {
	PageNum  int `form:"page_num" json:"page_num" `
	PageSize int `form:"page_size" json:"page_size"`
}
type CommentService struct {
	Vid     string `form:"vid" json:"vid"`         //视频
	Content string `form:"content" json:"content"` //内容

	ParentID uint `form:"parent_id" json:"parent_id"` //回复的评论的ID
	Level    uint `form:"level" json:"level"`         //回复的评论的所在楼层
}
type CommentShowService struct {
	Vid      string `form:"vid" json:"vid"` //视频
	PageNum  int    `form:"page_num" json:"page_num" `
	PageSize int    `form:"page_size" json:"page_size"`
}

type DanmuService struct {
	Vid  string `form:"vid" json:"vid"`
	Type uint   `form:"type" json:"type"` //类型0滚动;1顶部;2底部
	Text string `form:"text" json:"text"`
}

type DanmuShowService struct {
	Vid      string `form:"vid" json:"vid"` //视频
	PageNum  int    `form:"page_num" json:"page_num" `
	PageSize int    `form:"page_size" json:"page_size"`
}

type TransmitService struct {
	Vid     string `form:"vid" json:"vid"`
	Comment string `form:"comment" json:"comment"`
}

type ViewService struct {
	Vid string `form:"vid" json:"vid"`
}

type RankListService struct {
	PageNum  int64 `form:"page_num" json:"page_num" `
	PageSize int64 `form:"page_size" json:"page_size"`
}

func (service *VideoUploadService) VideoUpload(uid uint, resource string) serializer.Response {
	code := e.SUCCESS
	var video model.Video
	var user model.User
	var state uint
	model.DB.Table("user").Select("state").Where("id=?", uid).Scan(&state)
	if state == 1 {
		code = e.ErrorUserInBlackList
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	} else if state == 2 {
		code = e.ErrorUserHasBanned
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	//查
	if err := model.DB.Where("id=?", uid).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = e.ErrorNotExistUser
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}

		//如果是其他错误
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	/////////成功后分发av号
	rand.Seed(time.Now().Unix())
	avNum := rand.Int63()
	av := "av" + strconv.FormatInt(avNum, 10)

	var count int64
	for {
		model.DB.Model(&model.User{}).Where("vid=?", av).First(&user).Count(&count)
		if count == 1 {
			rand.Seed(time.Now().UnixNano())
			av = "av" + strconv.FormatInt(rand.Int63(), 10)
		} else {
			break
		}
	}

	videoModel := model.Video{
		Uid: uid,
		Vid: av,
	}
	err := model.DB.Create(&videoModel).Error
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	/////////视频上传

	err = model.DB.Model(&model.Video{}).Where("uid=?", uid).Where("vid=?", av).
		Update("resource", resource).Find(&video).
		Error
	if err != nil {
		code := e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	/////////视频标题
	err = model.DB.Model(&model.Video{}).Where("uid=?", uid).Where("vid=?", av).
		Update("title", service.Title).Find(&video).
		Error
	if err != nil {
		code := e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	////////视频简介
	err = model.DB.Model(&model.Video{}).Where("uid=?", uid).Where("vid=?", av).
		Update("desc", service.Desc).Find(&video).
		Error
	if err != nil {
		code := e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	return serializer.Response{
		Status: code,
		Data: serializer.VideoUp{
			Vid:      video.Vid,
			Resource: video.Resource,
			Title:    video.Title,
			Desc:     video.Desc,
		},
		Msg: e.GetMsg(code),
	}
}

// 文件夹生成
func (service *CollectionCreateService) CollectionCreate(uid uint) serializer.Response {
	code := e.SUCCESS
	var user model.User
	var collection model.Collection
	var cid int64

	//查
	if err := model.DB.Where("id=?", uid).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = e.ErrorNotExistUser
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}

		//如果是其他错误
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	model.DB.Where("uid=?", uid).Find(&collection).Limit(-1).Count(&cid)
	collectionModel := model.Collection{
		Uid:          uid,
		CollectionID: uint(cid),
		Name:         service.Name,
		Open:         service.Open,
	}
	err := model.DB.Create(&collectionModel).Error
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	return serializer.Response{
		Status: code,
		Data: serializer.CollectionCreate{
			Uid:  uid,
			Cid:  uint(cid),
			Name: service.Name,
			Open: service.Open,
		},
		Msg: e.GetMsg(code),
	}
}

// 全部文件夹展示
func (service *CollectionsShowService) CollectionsShow(uid uint) serializer.Response {
	code := e.SUCCESS
	var user model.User
	var collection []model.Collection
	var count int64 = 0
	//查
	if err := model.DB.Where("id=?", uid).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = e.ErrorNotExistUser
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}

		//如果是其他错误
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	if service.PageSize == 0 {
		service.PageSize = 10
	}
	err := model.DB.Model(&model.Collection{}).Where("uid=?", uid).
		Count(&count).Limit(service.PageSize).Offset((service.PageNum - 1) * service.PageSize).
		Find(&collection).Error
	if err != nil {
		code = e.ErrorShowCollections
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.BuildListResponse(serializer.ShowCollections(collection), uint(count))
}

func (service *CollectService) Collect(uid uint) serializer.Response {
	code := e.SUCCESS
	var user model.User
	var video model.Video
	var collect model.Collect
	var collection model.Collection
	var info string
	var count int64
	//查
	if err := model.DB.Where("id=?", uid).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = e.ErrorNotExistUser
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}

		//如果是其他错误
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	if err := model.DB.Where("vid=?", service.Vid).Where("review=?", 1).First(&video).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = e.ErrorNotExistVideo
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}

		//如果是其他错误
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	if service.State == 0 {
		model.DB.Where("uid=?", uid).Where("vid=?", service.Vid).Find(&collect).Count(&count)
		if count != 0 {
			code = e.ErrorHasExistCollect
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
			}
		}
		collectModel := model.Collect{
			Uid: uid,
			Vid: service.Vid,
			Cid: service.Cid,
		}
		// Set("gorm:foreign_key_check", 0)
		if err := model.DB.Create(&collectModel).Error; err != nil {
			code = e.ErrorDatabase
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}
		info = "收藏成功"
	} else if service.State == 1 {
		if err := model.DB.Where("uid=?", uid).Where("vid=?", service.Vid).Where("cid=?", service.Cid).First(&collect).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				code = e.ErrorNotExistCollect
				return serializer.Response{
					Status: code,
					Msg:    e.GetMsg(code),
					Error:  err.Error(),
				}
			}

			//如果是其他错误
			code = e.ErrorDatabase
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}
		err := model.DB.Where("uid=?", uid).Where("vid=?", service.Vid).Delete(&collect).Error
		if err != nil {
			code = e.ErrorDatabase
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}
		info = "取消收藏成功"
	} else {
		code = e.ErrorCollectAction
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	err := model.DB.Where("uid=?", uid).Where("collection_id=?", service.Cid).Find(&collection).Error
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	name := collection.Name
	return serializer.Response{
		Status: code,
		Data: serializer.Collect{
			Uid:  uid,
			Cid:  service.Cid,
			Name: name,
			Vid:  service.Vid,
			Info: info,
		},
		Msg: e.GetMsg(code),
	}
}

// 文件夹内容展示
func (service *CollectShowService) CollectShow(uid uint) serializer.Response {
	code := e.SUCCESS
	var user model.User
	var collect []model.Collect
	var count int64 = 0
	//查
	if err := model.DB.Where("id=?", uid).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = e.ErrorNotExistUser
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}

		//如果是其他错误
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	if service.PageSize == 0 {
		service.PageSize = 10
	}
	err := model.DB.Model(&model.Collect{}).Where("uid=?", uid).Where("cid=?", service.Cid).
		Count(&count).Limit(service.PageSize).Offset((service.PageNum - 1) * service.PageSize).
		Find(&collect).Error
	if err != nil {
		code = e.ErrorShowCollect
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.BuildListResponse(serializer.ShowCollects(collect), uint(count))
}

func (service *LikedService) Liked(uid uint) serializer.Response {
	code := e.SUCCESS
	var user model.User
	var video model.Video
	var liked model.Liked
	var count int64
	//查
	if err := model.DB.Where("id=?", uid).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = e.ErrorNotExistUser
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}

		//如果是其他错误
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	if err := model.DB.Where("vid=?", service.Vid).Where("review=?", 1).First(&video).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = e.ErrorNotExistVideo
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}

		//如果是其他错误
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	if err := model.DB.Where("vid=?", service.Vid).Find(&liked).Count(&count).Error; err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}

	}
	if count == 0 {
		likedModel := model.Liked{
			Uid: uid,
			Vid: service.Vid,
		}

		if err := model.DB.Create(&likedModel).Error; err != nil {
			code = e.ErrorDatabase
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}
	}
	var isLiked bool
	if service.IsLiked == 0 {
		isLiked = false
	} else if service.IsLiked == 1 {
		isLiked = true
	} else {
		code = e.ErrorLikedAction
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	err := model.DB.Model(&model.Liked{}).Where("uid=?", uid).Where("vid=?", service.Vid).
		Update("is_liked", isLiked).Find(&liked).
		Error
	if err != nil {
		code := e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	return serializer.Response{
		Status: code,
		Data: serializer.Liked{
			Uid:     liked.Uid,
			Vid:     liked.Vid,
			IsLiked: liked.IsLiked,
		},
		Msg: e.GetMsg(code),
	}
}

func (service *LikedShowService) LikedShow(uid uint) serializer.Response {
	code := e.SUCCESS
	var user model.User

	var liked []model.Liked
	var count int64 = 0
	//查
	if err := model.DB.Where("id=?", uid).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = e.ErrorNotExistUser
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}

		//如果是其他错误
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	if service.PageSize == 0 {
		service.PageSize = 10
	}
	err := model.DB.Model(&model.Liked{}).Where("uid=?", uid).
		Count(&count).Limit(service.PageSize).Offset((service.PageNum - 1) * service.PageSize).
		Find(&liked).Error
	if err != nil {
		code = e.ErrorShowLiked
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.BuildListResponse(serializer.ShowLikeds(liked), uint(count))
}

func (service *CommentService) Comment(uid uint) serializer.Response {
	code := e.SUCCESS
	var user model.User
	var video model.Video
	var comment model.Comment
	var count int64 = 0
	var state string
	var content string

	var userState uint
	model.DB.Table("user").Select("state").Where("id=?", uid).Scan(&userState)
	if userState == 1 {
		code = e.ErrorUserInBlackList
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	} else if userState == 2 {
		code = e.ErrorUserHasBanned
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	//查
	if err := model.DB.Where("id=?", uid).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = e.ErrorNotExistUser
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}

		//如果是其他错误
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	if err := model.DB.Where("vid=?", service.Vid).Where("review=?", 1).First(&video).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = e.ErrorNotExistVideo
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}

		//如果是其他错误
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	//对视频评论
	if service.Level == 0 && service.ParentID == 0 {
		content = service.Content
		commentModel := model.Comment{
			Vid:     service.Vid,
			Uid:     uid,
			Content: content,
			State:   1,
		}
		if err := model.DB.Create(&commentModel).Error; err != nil {

			code = e.ErrorDatabase
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}
		if err := model.DB.Model(&comment).Where("vid=?", service.Vid).Count(&count).Error; err != nil {
			code = e.ErrorDatabase
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}
		err := model.DB.Model(&model.Comment{}).Where("uid=?", uid).Where("vid=?", service.Vid).Where("content=?", service.Content).
			Where("state=?", 1).Update("level", count).Find(&comment).Error
		if err != nil {
			code := e.ErrorDatabase
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}
		state = "评论"
	} else if service.Level != 0 && service.ParentID != 0 {
		//回复
		if err := model.DB.Where("uid=?", service.ParentID).First(&comment).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				code = e.ErrorNotExistComment
				return serializer.Response{
					Status: code,
					Msg:    e.GetMsg(code),
					Error:  err.Error(),
				}
			}

			//如果是其他错误
			code = e.ErrorDatabase
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}
		if err := model.DB.Where("level=?", service.Level).First(&comment).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				code = e.ErrorNotExistLevel
				return serializer.Response{
					Status: code,
					Msg:    e.GetMsg(code),
					Error:  err.Error(),
				}
			}

			//如果是其他错误
			code = e.ErrorDatabase
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}
		var name string

		model.DB.Table("user").Select("user_name").Where("id=?", service.ParentID).Scan(&name)
		content = "@" + name + ":" + service.Content

		replyModel := model.Comment{
			Vid:      service.Vid,
			Uid:      uid,
			Level:    service.Level,
			Content:  content,
			ParentID: service.ParentID,
			State:    2,
		}
		if err := model.DB.Create(&replyModel).Error; err != nil {
			code = e.ErrorDatabase
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}
		state = "回复"
	} else {
		code = e.ErrorCommentAction
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	return serializer.Response{
		Status: code,
		Data: serializer.Comment{
			Vid:      service.Vid,
			Uid:      uid,
			Level:    service.Level,
			Content:  content,
			ParentID: service.ParentID,
			State:    state,
		},
		Msg: e.GetMsg(code),
	}

}

func (service *CommentShowService) CommentShow(uid uint) serializer.Response {
	code := e.SUCCESS
	var user model.User
	var video model.Video
	var comment []model.Comment
	var count int64

	//查
	if err := model.DB.Where("id=?", uid).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = e.ErrorNotExistUser
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}

		//如果是其他错误
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	if err := model.DB.Where("vid=?", service.Vid).First(&video).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = e.ErrorNotExistVideo
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}

		//如果是其他错误
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	if service.PageSize == 0 {
		service.PageSize = 10
	}
	err := model.DB.Model(&model.Comment{}).Where("uid=?", uid).Where("vid=?", service.Vid).
		Count(&count).Limit(service.PageSize).Offset((service.PageNum - 1) * service.PageSize).
		Find(&comment).Error
	if err != nil {
		code = e.ErrorShowComment
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.BuildListResponse(serializer.ShowComments(comment), uint(count))

}

func (service *DanmuService) Danmu(uid uint) serializer.Response {
	code := e.SUCCESS
	var user model.User
	var video model.Video

	var state uint
	model.DB.Table("user").Select("state").Where("id=?", uid).Scan(&state)
	if state == 1 {
		code = e.ErrorUserInBlackList
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	} else if state == 2 {
		code = e.ErrorUserHasBanned
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	//查
	if err := model.DB.Where("id=?", uid).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = e.ErrorNotExistUser
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}

		//如果是其他错误
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	if err := model.DB.Where("vid=?", service.Vid).Where("review=?", 1).First(&video).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = e.ErrorNotExistVideo
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}

		//如果是其他错误
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	if service.Type > 3 || service.Type < 1 {
		code = e.ErrorDanmu
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	danmuModel := model.Danmu{
		Vid:  service.Vid,
		Uid:  uid,
		Text: service.Text,
		Type: service.Type,
	}
	if err := model.DB.Create(&danmuModel).Error; err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Data: serializer.Danmu{
			Vid:  service.Vid,
			Type: service.Type,
			Text: service.Text,
			Uid:  uid,
		},
		Msg: e.GetMsg(code),
	}

}

func (service *DanmuShowService) DanmuShow(uid uint) serializer.Response {
	code := e.SUCCESS
	var user model.User
	var video model.Video
	var danmu []model.Danmu
	var count int64

	//查
	if err := model.DB.Where("id=?", uid).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = e.ErrorNotExistUser
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}

		//如果是其他错误
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	if err := model.DB.Where("vid=?", service.Vid).First(&video).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = e.ErrorNotExistVideo
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}

		//如果是其他错误
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	if service.PageSize == 0 {
		service.PageSize = 10
	}
	err := model.DB.Model(&model.Danmu{}).Where("uid=?", uid).Where("vid=?", service.Vid).
		Count(&count).Limit(service.PageSize).Offset((service.PageNum - 1) * service.PageSize).
		Find(&danmu).Error
	if err != nil {
		code = e.ErrorShowDanmu
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.BuildListResponse(serializer.ShowDanmus(danmu), uint(count))

}

func (service *TransmitService) Transmit(uid uint) serializer.Response {
	code := e.SUCCESS
	var user model.User
	var video model.Video

	//查
	if err := model.DB.Where("id=?", uid).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = e.ErrorNotExistUser
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}

		//如果是其他错误
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	if err := model.DB.Where("vid=?", service.Vid).Where("review=?", 1).First(&video).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = e.ErrorNotExistVideo
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}

		//如果是其他错误
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	var path string

	model.DB.Table("video").Select("resource").Where("vid=?", service.Vid).Scan(&path)

	transmitModel := model.Transmit{
		Uid:     uid,
		Vid:     service.Vid,
		Comment: service.Comment,
		Path:    path,
	}
	if err := model.DB.Create(&transmitModel).Error; err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Data: serializer.Transmit{
			Vid:     service.Vid,
			Uid:     uid,
			Comment: service.Comment,
			Path:    path,
		},
		Msg: e.GetMsg(code),
	}

}

func (service *ViewService) View(uid uint) serializer.Response {
	code := e.SUCCESS
	var user model.User
	var video model.Video

	//查
	if err := model.DB.Where("id=?", uid).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = e.ErrorNotExistUser
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}

		//如果是其他错误
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	if err := model.DB.Where("vid=?", service.Vid).Where("review=?", 1).First(&video).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = e.ErrorNotExistVideo
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}

		//如果是其他错误
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	// redis 数据库操作
	r := cache.RedisClient
	ctx := context.Background()

	click, _ := r.ZScore(ctx, cache.RankKey, service.Vid).Result()
	if click == 0 {
		err := r.ZAdd(ctx, cache.RankKey, redis.Z{Score: 1, Member: service.Vid}).Err()
		if err != nil {
			code = e.ErrorRedis
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}
	} else {
		err := r.ZIncrBy(ctx, cache.RankKey, 1, service.Vid).Err()
		if err != nil {
			code = e.ErrorRedis
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}
	}

	clicks, _ := r.ZScore(ctx, cache.RankKey, service.Vid).Result()

	err := model.DB.Model(&model.Video{}).Where("vid=?", service.Vid).
		Update("clicks", clicks).First(&video).Error
	if err != nil {
		code := e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	return serializer.Response{
		Status: code,
		Data: serializer.View{
			Vid:    service.Vid,
			Clicks: uint(clicks),
		},
		Msg: e.GetMsg(code),
	}

}

func (service *RankListService) RankList(uid uint) serializer.Response {
	code := e.SUCCESS
	var user model.User

	//查
	if err := model.DB.Where("id=?", uid).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = e.ErrorNotExistUser
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}

		//如果是其他错误
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	if service.PageSize == 0 {
		service.PageSize = 10
	}

	if service.PageNum == 0 {
		service.PageNum = 1
	}

	// redis 数据库操作
	r := cache.RedisClient
	ctx := context.Background()

	op := redis.ZRangeBy{
		Min:    "1",
		Max:    strconv.Itoa(math.MaxUint32),
		Offset: (service.PageNum - 1) * service.PageSize,
		Count:  service.PageSize,
	}
	vals, err := r.ZRevRangeByScoreWithScores(ctx, cache.RankKey, &op).Result()
	// fmt.Println(vals)
	if err != nil {
		code = e.ErrorRedis
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	return serializer.BuildRankListResponse(serializer.RankLists(vals))

}
