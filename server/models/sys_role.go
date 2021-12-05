package models

import (
	"reflect"
	"strings"
	"time"
)

type SysRole struct {
	RoleID      *int64     `json:"role_id" mysql:"role_id"`                        // 角色ID
	RoleName    *string    `json:"role_name" mysql:"role_name" binding:"required"` // 角色名
	State       *int       `json:"state" mysql:"state" binding:"required"`         // 角色状态，0 正常，1 禁用
	Deleted     *int       `json:"deleted" mysql:"deleted"`                        // 是否被删除 0 正常，1 被删除
	Description *string    `json:"description" mysql:"description"`                // 角色描述
	CreateTime  *time.Time `json:"create_time" mysql:"create_time"`                // 创建时间
	UpdateTime  *time.Time `json:"update_time" mysql:"update_time"`                // 更新时间
}

func (sr *SysRole) FieldToDB() (sqlFields string, scanFields []interface{}) {
	tmp := []string{}
	suType := reflect.TypeOf(sr).Elem()
	suValue := reflect.ValueOf(sr).Elem()
	for i := 0; i < suType.NumField(); i++ {
		tmp = append(tmp, suType.Field(i).Tag.Get("mysql"))
		scanFields = append(scanFields, suValue.Field(i).Addr().Interface())
	}
	sqlFields = "`" + strings.Join(tmp, "`,`") + "`"
	return
}
