package request

type LoginBody struct {
	UserName string `json:"username" binding:"required"` //用户名
	Password string `json:"password" binding:"required"` //密码
}
