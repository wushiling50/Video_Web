package service

import (
	"context"
	"errors"
	"main/Video_Web/cache"
	"main/Video_Web/model"
	"main/Video_Web/pkg/e"
	"main/Video_Web/serializer"
	"regexp"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type SearchUserService struct {
	Info     string `json:"info" form:"info"`
	PageNum  int    `form:"page_num" json:"page_num" `
	PageSize int    `form:"page_size" json:"page_size"`
}

type SearchVideoService struct {
	Info     string `json:"info" form:"info"`
	PageNum  int    `form:"page_num" json:"page_num" `
	PageSize int    `form:"page_size" json:"page_size"`
}

func (service *SearchUserService) SearchUser(uid uint) serializer.Response {
	code := e.SUCCESS
	var user model.User
	var users []model.User
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

	if service.PageSize == 0 {
		service.PageSize = 10
	}

	s1 := "[\\w!#$%&'*+/=?^_`{|}~-]+(?:\\.[\\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\\w](?:[\\w-]*[\\w])?\\.)+[\\w](?:[\\w-]*[\\w])?"
	emailRe, err := regexp.Compile(s1)
	if err != nil {
		code := e.ErrorRegexpParse
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	if b1 := emailRe.MatchString(service.Info); !b1 {
		userMsg, _ := strconv.Atoi(service.Info)
		if userMsg == 0 && service.Info != "0" {
			err := model.DB.Model(&model.User{}).
				Where("user_name LIKE ? OR  birthday LIKE ? ", "%"+service.Info+"%", "%"+service.Info+"%").
				Count(&count).Limit(service.PageSize).Offset((service.PageNum - 1) * service.PageSize).Find(&users).Error
			if err != nil {
				code := e.ErrorSearchUser
				return serializer.Response{
					Status: code,
					Msg:    e.GetMsg(code),
					Error:  err.Error(),
				}
			}
		} else {
			err := model.DB.Model(&model.User{}).
				Where("user_name LIKE ? OR birthday LIKE ? OR id=?", "%"+service.Info+"%", "%"+service.Info+"%", userMsg).
				Count(&count).Limit(service.PageSize).Offset((service.PageNum - 1) * service.PageSize).Find(&users).Error
			if err != nil {
				code := e.ErrorSearchUser
				return serializer.Response{
					Status: code,
					Msg:    e.GetMsg(code),
					Error:  err.Error(),
				}
			}
		}
	} else {
		err := model.DB.Model(&model.User{}).
			Where("email LIKE ? ", "%"+service.Info+"%").
			Count(&count).Limit(service.PageSize).Offset((service.PageNum - 1) * service.PageSize).Find(&users).Error
		if err != nil {
			code := e.ErrorSearchUser
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}
	}

	// redis 数据库操作
	r := cache.RedisClient
	ctx := context.Background()

	if err := r.RPush(ctx, cache.SearchUserKey, service.Info).Err(); err != nil {
		code = e.ErrorRedis
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	return serializer.BuildListResponse(serializer.SearchUsers(users), uint(count))

}

func (service *SearchVideoService) SearchVideo(uid uint) serializer.Response {
	code := e.SUCCESS
	var user model.User
	var videos []model.Video
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

	if service.PageSize == 0 {
		service.PageSize = 10
	}

	judge := strings.HasPrefix(service.Info, "av")

	if judge {
		err := model.DB.Model(&model.Video{}).
			Where("vid=?", service.Info).
			Count(&count).Limit(service.PageSize).Offset((service.PageNum - 1) * service.PageSize).Find(&videos).Error
		if err != nil {
			code := e.ErrorSearchVideo
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}
	} else {
		err := model.DB.Model(&model.Video{}).
			Where("uid LIKE ? OR title LIKE ? OR clicks LIKE ?", "%"+service.Info+"%", "%"+service.Info+"%", "%"+service.Info+"%").
			Count(&count).Limit(service.PageSize).Offset((service.PageNum - 1) * service.PageSize).Find(&videos).Error
		if err != nil {
			code := e.ErrorSearchVideo
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}
	}

	// redis 数据库操作
	r := cache.RedisClient
	ctx := context.Background()

	if err := r.RPush(ctx, cache.SearchVideoKey, service.Info).Err(); err != nil {
		code = e.ErrorRedis
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	return serializer.BuildListResponse(serializer.SearchVideos(videos), uint(count))

}
