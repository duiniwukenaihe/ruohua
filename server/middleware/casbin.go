package middleware

import (
	selfCasbin "github.com/chenhqchn/ruohua/server/utils/casbin"
	selfJwt "github.com/chenhqchn/ruohua/server/utils/jwt"
	"github.com/chenhqchn/ruohua/server/utils/resp"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckPrivilege() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claim := ctx.MustGet("userClaims").(*selfJwt.JwtClaims)
		path := ctx.Request.URL.Path
		method := ctx.Request.Method

		err := selfCasbin.Enforcer.LoadPolicy()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, resp.ApiInternalError)
			ctx.Abort()
			return
		}

		if claim.Username == "admin" {
			ctx.Next()
			return
		}

		// claim.Role 是角色切片,如：["角色1","角色2"]
		// 依次遍历角色切片，判断每一个角色是否有此接口权限
		allow := false
		for _, role := range claim.Role {
			result, err := selfCasbin.Enforcer.Enforce(role, path, method)
			if err != nil {
				continue
			}

			if !result {
				continue
			} else {
				allow = true
				break
			}
		}
		if !allow {
			ctx.JSON(http.StatusForbidden, gin.H{"status": 403, "msg": "403 Forbidden"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
