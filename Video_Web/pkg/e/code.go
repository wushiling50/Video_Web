package e

import (
	"encoding/json"
	"fmt"
	"main/Video_Web/serializer"
)

const (
	SUCCESS       = 200
	ERROR         = 500
	InvalidParams = 400

	//成员错误
	ErrorExistUser      = 10002
	ErrorNotExistUser   = 10003
	ErrorPasswordSame   = 10004
	ErrorFailEncryption = 10006
	ErrorNotCompare     = 10007
	ErrorMsgChange      = 10009
	ErrorRegexpParse    = 10010
	ErrorImgOpen        = 10011
	ErrorImgRead        = 10012
	ErrorImgUpload      = 10013

	//视频错误
	ErrorVideoOpen       = 10014
	ErrorVideoRead       = 10015
	ErrorShowCollections = 10016
	ErrorNotExistVideo   = 10017
	ErrorCollectAction   = 10018
	ErrorNotExistCollect = 10019
	ErrorHasExistCollect = 10020
	ErrorShowCollect     = 10021
	ErrorLikedAction     = 10022
	ErrorShowLiked       = 10023
	ErrorNotExistComment = 10024
	ErrorNotExistLevel   = 10025
	ErrorCommentAction   = 10026
	ErrorShowComment     = 10027
	ErrorDanmu           = 10028
	ErrorFileType        = 10029
	ErrorShowDanmu       = 10030

	//审核
	ErrorVideoNotCheck           = 10031
	ErrorUserInBlackList         = 10032
	ErrorNotRight                = 10033
	ErrorUserHasBanned           = 10034
	ErrorBlackListOrBannedAction = 10035
	ErrorDeleteComment           = 10036
	ErrorDeleteFailed            = 10037

	//搜索
	ErrorSearchUser  = 10038
	ErrorSearchVideo = 10039

	ErrorAuthCheckTokenFail    = 30001 //token 错误
	ErrorAuthCheckTokenTimeout = 30002 //token 过期
	ErrorAuthToken             = 30003
	ErrorAuth                  = 30004
	ErrorDatabase              = 40001
	ErrorSendEmail             = 40101
	ErrorRedis                 = 40201
)

// 返回错误信息
func ErrorResponse(err error) serializer.Response {
	if _, ok := err.(*json.UnmarshalTypeError); ok {
		return serializer.Response{
			Status: 40001,
			Msg:    "JSON类型不匹配",
			Error:  fmt.Sprint(err),
		}
	}
	return serializer.Response{
		Status: 40001,
		Msg:    "参数错误",
		Error:  fmt.Sprint(err),
	}
}
