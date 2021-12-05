package middleware

import (
	"fmt"
	"github.com/chenhqchn/ruohua/server/database/mysql"
	"github.com/chenhqchn/ruohua/server/models"
	"github.com/chenhqchn/ruohua/server/utils/resp"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 判断api请求的 username 是否存在
func CheckUserNameIsExist() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username := ctx.Param("username")
		var userdao mysql.UserDao
		var userinfo models.SysUser

		userinfo.UserName = username
		res, err := userdao.GetUser(userinfo)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, resp.ApiInternalError)
			ctx.Abort()
		}
		if res.UserID == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": 400, "msg": fmt.Sprintf("The username: %s does not exist", username)})
			ctx.Abort()
		}
		ctx.Next()
	}
}
