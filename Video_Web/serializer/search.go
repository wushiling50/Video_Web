package serializer

import "main/Video_Web/model"

type UserSearch struct {
	Uid      uint   `json:"uid"`
	Img      string `json:"img"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Gender   int    `json:"gender"`
	Birthday string `json:"birthday"`
	Sign     string `json:"sign"`
}

type VideoSearch struct {
	Uid    uint   `json:"uid"`
	Vid    string `json:"vid"`
	Title  string `json:"title"`
	Desc   string `json:"desc"`
	Clicks uint   `json:"clicks"`
}

func SearchUser(item model.User) UserSearch {
	return UserSearch{
		Uid:      item.ID,
		Img:      item.Img,
		UserName: item.UserName,
		Email:    item.Email,
		Gender:   item.Gender,
		Birthday: item.Birthday,
		Sign:     item.Sign,
	}
}

func SearchUsers(items []model.User) (users []UserSearch) {
	for _, item := range items {
		user := SearchUser(item)
		users = append(users, user)
	}
	return users
}

func SearchVideo(item model.Video) VideoSearch {
	return VideoSearch{
		Uid:    item.ID,
		Vid:    item.Vid,
		Title:  item.Title,
		Desc:   item.Desc,
		Clicks: uint(item.Clicks),
	}
}

func SearchVideos(items []model.Video) (videos []VideoSearch) {
	for _, item := range items {
		video := SearchVideo(item)
		videos = append(videos, video)
	}
	return videos
}
