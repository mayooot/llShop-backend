package dto

// SignUp 封装用户注册的请求体
type SignUp struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
	Code     string `json:"code" binding:"required"`
}

// Login 封装用户登录的请求体
type Login struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Infos 封装用户信息的请求体
type Infos struct {
	ID string `json:"id"`
	// 用户名，有长度限制
	Username string `json:"username"`
	// 手机号，必须为11位
	Phone string `json:"phone"`
	// 密码，同时包含大小写英文、数字、特殊字符
	Password string `json:"password"`
	// 邮箱
	Email string `json:"email"`
	// 性别，范围[0, 10]，1代表男，2代表女.....
	Gender uint8 `json:"gender"`
}
