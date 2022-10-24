package vo

// AttributeVO 商品属性名
type AttributeVO struct {
	KeyID           int64               `json:"keyID,string"`
	KeyName         string              `json:"keyName"`
	AttributeValues []*AttributeValueVO `json:"attributeValues"`
}

// AttributeValueVO 商品属性名下的规格条目
type AttributeValueVO struct {
	ValueID   int64  `json:"valueID,string"`
	ValueName string `json:"valueName"`
}
