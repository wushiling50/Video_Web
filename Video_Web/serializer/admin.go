package serializer

type Check struct {
	Vid    string `json:"vid"`
	Review string `json:"review"`
}

type BlackList struct {
	Uid   uint   `json:"uid"`
	State string `json:"state"`
}

type DeleteComment struct {
	Vid     string `json:"vid"`
	Uid     uint   `json:"uid"`
	Level   uint   `json:"level"`
	Content string `json:"content"`
	State   uint   `json:"state"`
	Msg     string `json:"msg"`
}
