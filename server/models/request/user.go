package request

type SysUserBody struct {
	DeptID         *int64  `json:"dept_id" binding:"required"`         // 用户所属部门ID
	UserName       *string `json:"username" binding:"required"`        // 用户名
	NickName       *string `json:"nickname" binding:"required"`        // 用户别名（强制实名）
	Email          *string `json:"email" binding:"required"`           // 邮箱地址
	Phone          *string `json:"phone" binding:"required"`           // 手机号
	Sex            *int    `json:"sex" binding:"required"`             // 性别 0 女性，1 男性
	FirstPassword  *string `json:"first_password" binding:"required"`  // 第一次输入密码
	SecondPassword *string `json:"second_password" binding:"required"` // 第二次输入密码
	State          *int    `json:"state" binding:"required"`           // 用户状态，0 正常，1 禁用
}

type UpdateUserBody struct {
	DeptID         *int64  `json:"dept_id"`         // 用户所属部门ID
	NickName       *string `json:"nickname"`        // 用户别名（强制实名）
	Email          *string `json:"email"`           // 邮箱地址
	Phone          *string `json:"phone"`           // 手机号
	Sex            *int    `json:"sex"`             // 性别 0 女性，1 男性
	FirstPassword  *string `json:"first_password"`  // 第一次输入密码
	SecondPassword *string `json:"second_password"` // 第二次输入密码
	State          *int    `json:"state"`           // 用户状态，0 正常，1 禁用
}
