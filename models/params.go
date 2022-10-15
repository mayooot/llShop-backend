package models

// 封装请求参数的结构体

// ParamSignUp 封装用户注册的请求体
type ParamSignUp struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
	Code     string `json:"code" binding:"required"`
}

// ParamLogin 封装用户登录的请求体
type ParamLogin struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ParamInfos 封装用户信息的请求体
type ParamInfos struct {
	Id       string `json:"id" binding:"required"`
	Username string `json:"username" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Avatar   string `json:"avatar" binding:"required"`
	Gender   string `json:"gender" binding:"required"`
}
