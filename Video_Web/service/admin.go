package service

import (
	"errors"
	"main/Video_Web/model"
	"main/Video_Web/pkg/e"
	"main/Video_Web/pkg/utils"
	"main/Video_Web/serializer"

	"gorm.io/gorm"
)

type AdminService struct {
	UserName string `form:"user_name" json:"user_name" binding:"required,min=3,max=15" `
	Password string `form:"password" json:"password" binding:"required,min=5,max=16"`
}

type CheckService struct {
	Vid string `form:"vid" json:"vid"`
}

type BlackListService struct {
	Uid   uint `form:"uid" json:"uid"`
	State uint `form:"state" json:"state"`
}

type DeleteService struct {
	Vid     string `form:"vid" json:"vid"`
	Uid     uint   `form:"uid" json:"uid"`
	Level   uint   `form:"level" json:"level"`
	Content string `form:"content" json:"content"`
	State   uint   `form:"state" json:"state"`
}

func (service *AdminService) AdminRegister() serializer.Response {
	code := e.SUCCESS
	var user model.User
	var count int64
	model.DB.Model(&model.User{}).Where("user_name=?", service.UserName).First(&user).Count(&count)
	if count == 1 {
		code = e.ErrorExistUser
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	user.UserName = service.UserName
	user.Role = 1
	//密码加密
	if err := user.SetPassword(service.Password); err != nil {
		code = e.ErrorFailEncryption
		return serializer.Response{
			Status: code,
			Msg:    err.Error(),
			Error:  err.Error(),
		}
	}

	// 创建用户
	if err := model.DB.Create(&user).Error; err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}

}

func (service *AdminService) AdminLogin() serializer.Response {
	code := e.SUCCESS
	var user model.User

	//查找数据库中是否存在该用户
	if err := model.DB.Where("user_name=?", service.UserName).Where("role=?", 1).First(&user).Error; err != nil {
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
	if !user.CheckPassword(service.Password) {
		code = e.ErrorNotCompare
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	//发送token，为了其他功能需要身份验证所给前端存储
	token, err := utils.GenerateToken(user.ID, service.UserName, service.Password)
	if err != nil {
		code = e.ErrorAuthToken
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Data: serializer.TokenData{
			User:  serializer.BuildUser(user),
			Token: token,
		},
		Msg: e.GetMsg(code),
	}
}

func (service *CheckService) Check(uid uint) serializer.Response {
	code := e.SUCCESS
	var user model.User
	var video model.Video

	if err := model.DB.Where("id=?", uid).Where("role=?", 1).First(&user).Error; err != nil {
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

	err := model.DB.Model(&model.Video{}).Where("vid=?", service.Vid).Where("review=?", 0).
		Update("review", 1).Find(&video).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = e.ErrorVideoNotCheck
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}

		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	var review string = "已过审"

	return serializer.Response{
		Status: code,
		Data: serializer.Check{
			Vid:    service.Vid,
			Review: review,
		},
		Msg: e.GetMsg(code),
	}
}

func (service *BlackListService) BlackList(uid uint) serializer.Response {
	code := e.SUCCESS
	var user model.User
	var state uint
	var msg string
	var role int
	if err := model.DB.Where("id=?", uid).Where("role=?", 1).First(&user).Error; err != nil {
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

	model.DB.Table("user").Select("role").Where("id=?", service.Uid).Scan(&role)
	if role == 1 {
		code = e.ErrorNotRight
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	model.DB.Table("user").Select("state").Where("id=?", service.Uid).Scan(&state)

	if state == service.State {
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
		} else if state == 0 {
			code = e.ErrorBlackListOrBannedAction
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
			}
		}
	} else {

		err := model.DB.Model(&model.User{}).Where("id=?", service.Uid).Update("state", service.State).Find(&user).Error
		if err != nil {
			code = e.ErrorDatabase
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}

		model.DB.Table("user").Select("state").Where("id=?", service.Uid).Scan(&state)

		if state == 1 {
			msg = "该用户拉黑成功"
		} else if state == 2 {
			msg = "该用户封禁中"
		} else if state == 0 {
			msg = "该用户解禁"
		} else {
			code = e.ErrorBlackListOrBannedAction
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
			}
		}

	}

	return serializer.Response{
		Status: code,
		Data: serializer.BlackList{
			Uid:   service.Uid,
			State: msg,
		},
		Msg: e.GetMsg(code),
	}
}

func (service *DeleteService) Delete(uid uint) serializer.Response {
	code := e.SUCCESS
	var user model.User
	var video model.Video
	var comment model.Comment
	var msg string

	if err := model.DB.Where("id=?", uid).Where("role=?", 1).First(&user).Error; err != nil {
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

	if service.State == 2 {
		err := model.DB.Model(&model.Comment{}).Where("vid=?", service.Vid).Where("uid=?", service.Uid).
			Where("level=?", service.Level).Where("content=?", service.Content).First(&comment).
			Delete(&comment).Error
		if err != nil {
			code = e.ErrorDeleteFailed
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}
		msg = "删除回复成功"
	} else if service.State == 1 {
		if service.Uid != 0 && service.Content != "" {
			code = e.ErrorDeleteComment
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
			}
		}
		err := model.DB.Model(&model.Comment{}).Where("vid=?", service.Vid).Where("level=?", service.Level).
			Delete(&comment).Error
		if err != nil {
			code = e.ErrorDeleteFailed
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}
		msg = "删除评论成功"
	} else {
		code = e.ErrorDeleteComment
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	return serializer.Response{
		Status: code,
		Data: serializer.DeleteComment{
			Vid:     service.Vid,
			Uid:     service.Uid,
			Level:   service.Level,
			Content: service.Content,
			State:   service.State,
			Msg:     msg,
		},
		Msg: e.GetMsg(code),
	}
}
