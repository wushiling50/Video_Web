package model

func migration() {

	err := DB.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(
			&User{},
			&Notice{},

			&Video{},
			&Liked{},      //点赞表
			&Collect{},    //收藏表
			&Collection{}, //收藏夹表
			&Comment{},    //评论表
			&Danmu{},      //弹幕表
			&Transmit{},   //转发
		)

	if err != nil {
		panic("数据库迁移失败")
	}

}
