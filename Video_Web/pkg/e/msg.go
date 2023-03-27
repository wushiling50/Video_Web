package e

var MsgFlags = map[int]string{
	SUCCESS:       "操作成功",
	ERROR:         "操作失败",
	InvalidParams: "请求参数错误",

	ErrorExistUser:    "用户已存在",
	ErrorNotExistUser: "用户不存在",

	ErrorAuthCheckTokenFail:      "Token鉴权失败",
	ErrorAuthCheckTokenTimeout:   "Token已超时",
	ErrorAuthToken:               "Token生成失败",
	ErrorAuth:                    "Token错误",
	ErrorNotCompare:              "不匹配",
	ErrorPasswordSame:            "应当更改不相同的密码",
	ErrorDatabase:                "数据库操作出错,请重试",
	ErrorMsgChange:               "更改信息有误",
	ErrorRegexpParse:             "正则表达式解析错误",
	ErrorImgOpen:                 "头像文件打开错误",
	ErrorImgRead:                 "头像文件读取错误",
	ErrorImgUpload:               "头像文件上传错误",
	ErrorVideoOpen:               "视频文件打开错误",
	ErrorVideoRead:               "视频文件读取错误",
	ErrorShowCollections:         "查看文件夹概况失败",
	ErrorNotExistVideo:           "视频不存在",
	ErrorCollectAction:           "错误的收藏行为",
	ErrorNotExistCollect:         "此收藏不存在",
	ErrorHasExistCollect:         "此收藏已存在",
	ErrorShowCollect:             "收藏夹内容展示失败",
	ErrorLikedAction:             "点赞操作失败",
	ErrorShowLiked:               "点赞内容展示失败",
	ErrorNotExistComment:         "该评论不存在",
	ErrorNotExistLevel:           "该楼层不存在",
	ErrorCommentAction:           "评论方式错误",
	ErrorShowComment:             "展示评论列表失败",
	ErrorDanmu:                   "弹幕发送格式错误",
	ErrorFileType:                "文件格式错误",
	ErrorShowDanmu:               "弹幕池展示失败",
	ErrorVideoNotCheck:           "视频待审核",
	ErrorUserInBlackList:         "用户已被拉黑",
	ErrorNotRight:                "没有权限拉黑管理员",
	ErrorUserHasBanned:           "该用户已被封禁",
	ErrorBlackListOrBannedAction: "错误的拉黑或封禁操作",
	ErrorDeleteComment:           "删除评论错误",
	ErrorDeleteFailed:            "删除评论或回复失败",
	ErrorSendEmail:               "发送邮件错误",
	ErrorRedis:                   "redis数据库错误",
	ErrorSearchUser:              "用户查找失败",
	ErrorSearchVideo:             "视频查找失败",
}

// GetMsg 获取状态码对应信息
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}
