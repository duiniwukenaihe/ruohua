package middleware

import (
	"fmt"
	"github.com/chenhqchn/ruohua/server/database/mysql"
	"github.com/chenhqchn/ruohua/server/models"
	"github.com/chenhqchn/ruohua/server/utils/resp"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 判断api请求的 rolename 是否存在
func CheckRoleIsExist() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rolename := ctx.Param("rolename")
		var roledao mysql.RoleDao
		var roleInfo models.SysRole
		roleInfo.RoleName = &rolename
		res, err := roledao.GetRole(roleInfo)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, resp.ApiInternalError)
			ctx.Abort()
		}
		if res.RoleID == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": 400, "msg": fmt.Sprintf("The rolename: %s does not exist", rolename)})
			ctx.Abort()
		}
		ctx.Next()
	}
}
