package routes

import (
	"main/Video_Web/api"
	"main/Video_Web/middleware"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	store := cookie.NewStore([]byte("something-very-secret"))
	r.Use(sessions.Sessions("mysession", store))
	r.StaticFS("/static", http.Dir("./static"))
	v1 := r.Group("api/v1")
	{
		//用户操作
		v1.POST("user/register", api.UserRegister) // 注册
		v1.POST("user/login", api.UserLogin)       //登录

		authed := v1.Group("/")
		authed.Use(middleware.JWT())
		{
			//用户模块
			v1.POST("user/sendingemail", api.SendingEmail)   //发送邮件
			v1.POST("user/handlingemail", api.HandlingEmail) //绑定 ， 解绑 ，更换密码
			v1.PUT("user/change", api.UserChange)            //更改密码
			v1.GET("user/msg/show", api.ShowMsg)             //展示个人信息
			v1.PUT("user/msg/change", api.ChangeMsg)         //更改个人信息
			v1.POST("user/img", api.UploadImg)               //更换头像

			//视频模块
			v1.POST("video/upload", api.VideoUpload)                 //上传视频及其信息
			v1.POST("video/create/collection", api.CollectionCreate) //创建收藏夹
			v1.GET("video/show/collections", api.CollectionsShow)    //展示收藏夹概况
			v1.POST("video/collect", api.Collect)                    //收藏与取消收藏
			v1.GET("video/collect/show", api.CollectShow)            //展示收藏夹内容
			v1.POST("video/liked", api.Liked)                        //点赞
			v1.GET("video/liked/show", api.LikedShow)                //展示点赞列表
			v1.POST("video/comment", api.Comment)                    //评论
			v1.GET("video/comment/show", api.CommentShow)            //展示视频的评论列表
			v1.POST("video/transmit", api.Transmit)                  //转发
			v1.POST("video/danmu", api.Danmu)                        //发送弹幕
			v1.GET("video/danmu/show", api.DanmuShow)                //展示视频的弹幕池
			v1.POST("video/view", api.View)                          //增加视频点击量
			v1.GET("video/rankList", api.RankList)                   //点击量排行榜

			//管理员模块
			v1.POST("admin/register", api.AdminRegister) //管理员注册
			v1.POST("admin/login", api.AdminLogin)       //管理员登录
			v1.PUT("admin/check", api.Check)             //视频审核
			v1.PUT("admin/blacklist", api.BlackList)     //用户管理
			v1.DELETE("admin/delete", api.Delete)        //删除评论

			//搜索模块
			v1.POST("search/user", api.SearchUser)   //搜索用户
			v1.POST("search/video", api.SearchVideo) //搜索视频

		}
	}
	return r
}
