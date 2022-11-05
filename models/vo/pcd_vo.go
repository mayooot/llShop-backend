package vo

// PCDDicVO 展示所有可选的收货地址
type PCDDicVO struct {
	ID      int         `json:"id"`
	Name    string      `json:"name"`
	PicList []*PCDDicVO `json:"picList"`
}
