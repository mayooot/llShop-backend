package vo

// Attribute 商品属性名
type Attribute struct {
	KeyID           int64             `json:"keyID,string"`
	KeyName         string            `json:"keyName"`
	AttributeValues []*AttributeValue `json:"attributeValues"`
}

// AttributeValue 商品属性名下的规格条目
type AttributeValue struct {
	ValueID   int64  `json:"valueID"`
	ValueName string `json:"valueName"`
}
