package serializer

import (
	"main/Video_Web/model"

	"github.com/redis/go-redis/v9"
)

type VideoUp struct {
	Resource string `json:"resource"`
	Vid      string `json:"vid"`
	Title    string `json:"title"`
	Desc     string `json:"desc"`
}

type CollectionCreate struct {
	Uid  uint   `json:"uid"`
	Cid  uint   `json:"cid"`
	Name string `json:"name"`
	Open uint   `json:"open"`
}

type CollectionsShow struct {
	Uid  uint   `json:"uid"`
	Cid  uint   `json:"cid"`
	Name string `json:"name"`
	Open string `json:"open"`
}

type Collect struct {
	Uid  uint   `json:"uid"`
	Cid  uint   `json:"cid"`
	Name string `json:"name"`
	Vid  string `json:"vid"`
	Info string `json:"info"`
}

type CollectsShow struct {
	Uid uint   `json:"uid"`
	Cid uint   `json:"cid"`
	Vid string `json:"vid"`
}

type Liked struct {
	Uid     uint   `json:"uid"`
	Vid     string `json:"vid"`
	IsLiked bool   `json:"is_liked"`
}

type LikedShow struct {
	Vid     string `json:"vid"`
	IsLiked bool   `json:"is_liked"`
}

type Comment struct {
	Vid      string `json:"vid"`
	Uid      uint   `json:"uid"`
	Level    uint   `json:"level"`
	Content  string `json:"content"`
	ParentID uint   `json:"parent_id"`
	State    string `json:"state"` //mark
}

type CommentShow struct {
	Vid     string `json:"vid"`
	Uid     uint   `json:"uid"`
	Level   uint   `json:"lecel"`
	Content string `json:"content"`

	ParentID uint   `json:"parent_id"`
	State    string `json:"state"`
}

type Danmu struct {
	Vid  string `json:"vid"`
	Type uint   `json:"type"` //类型0滚动;1顶部;2底部
	Text string `json:"text"`
	Uid  uint   `json:"uid"` //发送人的id
}

type DanmuShow struct {
	Vid  string `json:"vid"`
	Type uint   `json:"type"` //类型0滚动;1顶部;2底部
	Text string `json:"text"`
	Uid  uint   `json:"uid"` //发送人的id
}

type Transmit struct {
	Uid     uint   `json:"uid"`
	Vid     string `json:"vid"`
	Path    string `json:"path"`
	Comment string `json:"comment"`
}

type View struct {
	Vid    string `json:"vid"`
	Clicks uint   `json:"clicks"`
}

type RankList struct {
	Vid    interface{} `json:"vid"`
	Clicks uint64      `json:"clicks"`
}

func ShowCollection(item model.Collection) CollectionsShow {
	var open string
	if item.Open == 0 {
		open = "私密"
	} else if item.Open == 1 {
		open = "公开"
	}
	return CollectionsShow{
		Uid:  item.Uid,
		Cid:  item.CollectionID,
		Name: item.Name,
		Open: open,
	}
}

func ShowCollections(items []model.Collection) (collections []CollectionsShow) {
	for _, item := range items {
		collection := ShowCollection(item)
		collections = append(collections, collection)
	}
	return collections
}

func ShowCollect(item model.Collect) CollectsShow {
	return CollectsShow{
		Uid: item.Uid,
		Cid: item.Cid,
		Vid: item.Vid,
	}
}

func ShowCollects(items []model.Collect) (collects []CollectsShow) {
	for _, item := range items {
		collect := ShowCollect(item)
		collects = append(collects, collect)
	}
	return collects
}

func ShowLiked(item model.Liked) LikedShow {
	return LikedShow{
		Vid:     item.Vid,
		IsLiked: item.IsLiked,
	}
}

func ShowLikeds(items []model.Liked) (likeds []LikedShow) {
	for _, item := range items {
		liked := ShowLiked(item)
		likeds = append(likeds, liked)
	}
	return likeds
}

func ShowComment(item model.Comment) CommentShow {
	var state string
	if item.State == 1 {
		state = "评论"
	} else if item.State == 2 {
		state = "回复"
	}
	return CommentShow{
		Vid:      item.Vid,
		Uid:      item.Uid,
		Level:    item.Level,
		Content:  item.Content,
		ParentID: item.ParentID,
		State:    state,
	}
}

func ShowComments(items []model.Comment) (comments []CommentShow) {
	for _, item := range items {
		comment := ShowComment(item)
		comments = append(comments, comment)
	}
	return comments
}

func ShowDanmu(item model.Danmu) DanmuShow {
	return DanmuShow{
		Vid:  item.Vid,
		Type: item.Type,
		Text: item.Text,
		Uid:  item.Uid,
	}
}

func ShowDanmus(items []model.Danmu) (danmus []DanmuShow) {
	for _, item := range items {
		danmu := ShowDanmu(item)
		danmus = append(danmus, danmu)
	}
	return danmus
}

func RanksList(item redis.Z) RankList {
	return RankList{
		Vid:    item.Member,
		Clicks: uint64(item.Score),
	}
}

func RankLists(items []redis.Z) (lists []RankList) {
	for _, item := range items {
		list := RanksList(item)
		lists = append(lists, list)
	}
	return lists
}

type rankList struct {
	Item interface{} `json:"item"`
}

// 排行榜专属序列化器
func BuildRankListResponse(items interface{}) Response {
	return Response{
		Status: 200,
		Data: rankList{
			Item: items,
		},
		Msg: "ok",
	}
}
