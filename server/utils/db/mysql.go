package db

import (
	"fmt"
	"strings"
)

// 根据要更新的字段，构造 sql 语句
func BuildCommon(sql string, info map[string]interface{}) string {
	tmpSlice := []string{}
	for field, value := range info {
		tmpSlice = append(tmpSlice, fmt.Sprintf("%v='%v'", field, value))
	}
	return fmt.Sprintf(sql, strings.Join(tmpSlice, ","))
}