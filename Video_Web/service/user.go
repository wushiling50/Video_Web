package service

import (
	"errors"
	"fmt"
	"main/Video_Web/conf"
	"main/Video_Web/model"
	"main/Video_Web/pkg/e"
	"main/Video_Web/pkg/utils"
	"main/Video_Web/serializer"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mail.v2"
	"gorm.io/gorm"
)

type UserService struct {
	UserName string `form:"user_name" json:"user_name" binding:"required,min=3,max=15" `
	Password string `form:"password" json:"password" binding:"required,min=5,max=16"`
}
type UserChangeService struct {
	UserName    string `form:"user_name" json:"user_name" binding:"required,min=3,max=15"`
	PrePassword string `form:"prepassword" json:"pre_password" binding:"required,min=3,max=15"`
	NewPassword string `form:"newpassword" json:"new_password" binding:"required,min=3,max=15"`
}
type UserShowService struct {
}
type UserMsgService struct {
	UserName string `form:"user_name" json:"user_name" binding:"required,min=3,max=15" `
	Email    string `form:"email" json:"email" binding:"email"`
	Gender   int    `form:"gender" json:"gender"  ` //0为男性 ， 1为女性
	Birthday string `form:"birthday" json:"birthday"`
	Sign     string `form:"sign" json:"sign" binding:"required,min=1,max=30"`
}

type UserImgUploadService struct {
}

type SendingEmailService struct {
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`

	Operation uint `form:"operation" json:"operation"` //1为绑定邮箱 2为解绑邮箱 3为改密码
}

type HandlingEmailService struct {
}

func (service *UserService) Register() serializer.Response {
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

func (service *UserService) Login() serializer.Response {
	code := e.SUCCESS
	var user model.User
	var uid uint
	var state uint
	model.DB.Table("user").Select("id").Where("user_name=?", service.UserName).Scan(&uid)

	model.DB.Table("user").Select("state").Where("id=?", uid).Scan(&state)

	if state == 2 {
		code = e.ErrorUserHasBanned
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	//查找数据库中是否存在该用户
	if err := model.DB.Where("user_name=?", service.UserName).First(&user).Error; err != nil {
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

// 更改密码
func (service *UserChangeService) UserChange(uid uint) serializer.Response {
	var user model.User
	code := e.SUCCESS
	//查
	if err := model.DB.Where("user_name=?", service.UserName).First(&user).Error; err != nil {
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
	username := user.UserName
	//比较之前的密码是否与库中密文相同
	if !user.CheckPassword(service.PrePassword) {
		code = e.ErrorNotCompare
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	prepassword := service.PrePassword
	//前后密码不应相同
	if (service.NewPassword) == (prepassword) {
		code = e.ErrorPasswordSame
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	newpassword := service.NewPassword

	//密码加密
	bytes2, err := bcrypt.GenerateFromPassword([]byte(service.NewPassword), 12)
	if err != nil {
		code = e.ErrorFailEncryption
		return serializer.Response{
			Status: code,
			Msg:    err.Error(),
			Error:  err.Error(),
		}
	}
	PD2 := string(bytes2)
	//更改密码
	err = model.DB.Model(&model.User{}).Where("id=?", uid).
		Update("password_digest", PD2).Find(&user).
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
		Data: serializer.ChangePassword{
			UserName:    username,
			PrePassword: prepassword,
			NewPassword: newpassword,
		},
		Msg: e.GetMsg(code),
	}
}

// 展示用户个人信息
func (service *UserShowService) ShowMsg(uid uint) serializer.Response {
	var user model.User
	code := e.SUCCESS
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

	return serializer.Response{
		Status: code,
		Data: serializer.UserMsg{
			Uid:      uid,
			Img:      user.Img,
			UserName: user.UserName,
			Email:    user.Email,
			Gender:   user.Gender,
			Birthday: user.Birthday,
			Sign:     user.Sign,
		},
		Msg: e.GetMsg(code),
	}
}

// 修改用户个人信息
func (service *UserMsgService) MsgChange(uid uint) serializer.Response {
	var user model.User
	code := e.SUCCESS
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

	//更改用户名
	err := model.DB.Model(&model.User{}).Where("id=?", uid).
		Update("user_name", service.UserName).Find(&user).
		Error
	if err != nil {
		code := e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	//更改邮箱
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
	if b1 := emailRe.MatchString(service.Email); !b1 {
		code := e.ErrorMsgChange
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	email := string(emailRe.Find([]byte(service.Email)))
	err = model.DB.Model(&model.User{}).Where("id=?", uid).
		Update("email", email).Find(&user).
		Error
	if err != nil {
		code := e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	//性别
	if service.Gender < 0 || service.Gender > 1 {
		code := e.ErrorMsgChange
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	err = model.DB.Model(&model.User{}).Where("id=?", uid).
		Update("gender", service.Gender).Find(&user).
		Error
	if err != nil {
		code := e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	//生日
	s2 := "[1-2][0-9]{3}年([1-9]|1[0-2])月([0-9]|[1-2][0-9]|3[0-1])日"
	birthdayRe, err := regexp.Compile(s2)
	if err != nil {
		code := e.ErrorRegexpParse
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	if b2 := birthdayRe.MatchString(service.Birthday); !b2 {
		code := e.ErrorMsgChange
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	birthday := string(birthdayRe.Find([]byte(service.Birthday)))
	err = model.DB.Model(&model.User{}).Where("id=?", uid).
		Update("birthday", birthday).Find(&user).
		Error
	if err != nil {
		code := e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	//签名
	err = model.DB.Model(&model.User{}).Where("id=?", uid).
		Update("sign", service.Sign).Find(&user).
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
		Data: serializer.UserMsg{
			Uid:      uid,
			UserName: service.UserName,
			Email:    service.Email,
			Gender:   service.Gender,
			Birthday: service.Birthday,
			Sign:     service.Sign,
		},
		Msg: e.GetMsg(code),
	}
}

func (service *UserImgUploadService) UploadImg(uid uint, path string) serializer.Response {
	code := e.SUCCESS
	var user model.User
	var err error

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

	err = model.DB.Model(&model.User{}).Where("id=?", uid).
		Update("img", path).Find(&user).Error
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
		Data: serializer.UserMsg{
			Uid:      uid,
			Img:      user.Img,
			UserName: user.UserName,
			Email:    user.Email,
			Gender:   user.Gender,
			Birthday: user.Birthday,
			Sign:     user.Sign,
		},
		Msg: e.GetMsg(code),
	}
}

func (service *SendingEmailService) SendingEmail(uid uint) serializer.Response {
	code := e.SUCCESS
	var address string
	var user model.User
	var mailStr string

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

	token, err := utils.GenerateEmailToken(uid, service.Operation, service.Email, service.Password)
	if err != nil {
		code = e.ErrorAuthToken
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	model.DB.Table("notice").Select("text").Where("id=?", service.Operation).Scan(&mailStr)

	// sp := base64.StdEncoding.EncodeToString([]byte(conf.SmtpPass))
	// se := base64.StdEncoding.EncodeToString([]byte(conf.SmtpEmail))

	address = conf.ValidEmail + token //发送方
	mailText := mailStr + "\n" + address
	m := mail.NewMessage()
	m.SetHeader("From", conf.SmtpEmail)
	m.SetHeader("To", service.Email)
	m.SetHeader("Subject", "hello")
	m.SetBody("text/html", mailText)
	d := mail.NewDialer(conf.SmtpHost, 465, conf.SmtpEmail, conf.SmtpPass)
	d.StartTLSPolicy = mail.MandatoryStartTLS
	if err := d.DialAndSend(m); err != nil {
		code = e.ErrorSendEmail
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	// em := email.NewEmail()
	// em.From = conf.SmtpEmail
	// em.To = []string{service.Email}
	// em.Subject = "hello"
	// em.Text = []byte(mailText)

	// if err := em.Send("smtp.qq.com:25", smtp.PlainAuth("", conf.SmtpEmail, conf.SmtpPass, conf.SmtpHost)); err != nil {
	// 	code = e.ErrorSendEmail
	// 	return serializer.Response{
	// 		Status: code,
	// 		Msg:    e.GetMsg(code),
	// 		Error:  err.Error(),
	// 	}
	// }
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

func (service *HandlingEmailService) HandlingEmail(token string) serializer.Response {
	code := e.SUCCESS
	var uid uint
	var email string
	var password string
	var operation uint
	var user *model.User

	if token == "" {
		code = e.InvalidParams
	} else {
		claims, err := utils.ParseEmailToken(token)
		if err != nil {
			code = e.ErrorAuthCheckTokenFail
		} else if time.Now().Unix() > claims.ExpiresAt {
			code = e.ErrorAuthCheckTokenTimeout
		} else {
			uid = claims.Uid
			email = claims.Email
			password = claims.Password
			operation = claims.Operation
			fmt.Println(uid, " ", email, " ", operation)
		}

	}
	if code != e.SUCCESS {
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
	// model.DB.Model(&model.User{}).Where("id=?", uid).Update("email", email).First(&user)
	if operation == 1 {
		//1:绑定邮箱
		user.Email = email
	} else if operation == 2 {
		//2：解绑邮箱
		user.Email = ""
	} else if operation == 3 {
		//3：修改密码
		err := user.SetPassword(password)
		if err != nil {
			code = e.ErrorDatabase
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
			}
		}
	}
	err := model.DB.Model(&model.User{}).Where("id=?", uid).
		Updates(&user).Error
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	return serializer.Response{
		Status: code,
		Data: serializer.UserMsg{
			Uid:      uid,
			Img:      user.Img,
			UserName: user.UserName,
			Email:    user.Email,
			Gender:   user.Gender,
			Birthday: user.Birthday,
			Sign:     user.Sign,
		},
		Msg: e.GetMsg(code),
	}
}
