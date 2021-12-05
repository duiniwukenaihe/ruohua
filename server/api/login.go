package api

import (
	"fmt"
	"github.com/chenhqchn/ruohua/server/models/request"
	"github.com/chenhqchn/ruohua/server/service"
	"github.com/chenhqchn/ruohua/server/utils/config"
	"github.com/chenhqchn/ruohua/server/utils/resp"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginApi struct {
	loginService service.LoginService
}

func (a *LoginApi) Login(ctx *gin.Context) {
	loginBody := request.LoginBody{}

	config.L().Debug("Login parameter binding check")
	if err := ctx.BindJSON(&loginBody); err == nil {
		token, err := a.loginService.Login(ctx, loginBody.UserName, loginBody.Password)
		fmt.Println(err)
		if err != nil {
			switch err {
			case resp.UserNotExistError, resp.UserStateError, resp.UserPasswordNotMatchError:
				ctx.JSON(http.StatusBadRequest, gin.H{"status": 400, "msg": err.Error()})
				return
			default:
				ctx.JSON(http.StatusInternalServerError, resp.ApiInternalError)
				return
			}
		}
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "token": token, "msg": "User login successfully"})
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "400", "msg": resp.UserLoginParaError.Error(),
			"data": loginBody})
	}
}
