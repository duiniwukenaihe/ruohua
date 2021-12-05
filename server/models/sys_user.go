package models

import (
	"reflect"
	"strings"
	"time"
)

type SysUser struct {
	UserID     int64     `json:"user_id" mysql:"user_id"`         // 用户ID
	DeptID     int64     `json:"dept_id" mysql:"dept_id"`         // 用户所属部门ID
	UserName   string    `json:"username" mysql:"username"`       // 用户名
	NickName   string    `json:"nickname" mysql:"nickname"`       // 用户别名
	Email      string    `json:"email" mysql:"email"`             // 邮箱地址
	Phone      string    `json:"phone" mysql:"phone"`             // 手机号
	Sex        int       `json:"sex" mysql:"sex"`                 // 性别 0 女性，1 男性
	Password   string    `json:"-" mysql:"password"`              // 密码
	State      int       `json:"state" mysql:"state"`             // 用户状态，0 正常，1 禁用
	Deleted    int       `json:"-" mysql:"deleted"`               // 是否已被删除，0正常，1 被删除
	LoginIP    string    `json:"login_ip" mysql:"login_ip"`       // 登录IP
	LastLogin  string    `json:"last_login" mysql:"last_login"`   // 上一次登录时间
	CreateTime time.Time `json:"create_time" mysql:"create_time"` // 创建时间
	UpdateTime time.Time `json:"update_time" mysql:"update_time"` // 更新时间
}

// Solve the problem of inconsistency between strcut attribute and DB internal field order
// sqlFields is used to construct sql statement fields，For example: select sqlFields from table
// scanFields is used to construct rows.scan's parameter
func (su *SysUser) FieldToDB() (sqlFields string, scanFields []interface{}) {
	tmp := []string{}
	suType := reflect.TypeOf(su).Elem()
	suValue := reflect.ValueOf(su).Elem()
	for i := 0; i < suType.NumField(); i++ {
		tmp = append(tmp, suType.Field(i).Tag.Get("mysql"))
		scanFields = append(scanFields, suValue.Field(i).Addr().Interface())
	}
	sqlFields = "`" + strings.Join(tmp, "`,`") + "`"
	return
}
