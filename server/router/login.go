package router

import (
	"github.com/chenhqchn/ruohua/server/api"
	"github.com/chenhqchn/ruohua/server/utils/config"
	"github.com/gin-gonic/gin"
)

func initLoginRouter(router *gin.RouterGroup)  {
	config.L().Debug("Init the login api group")
	loginApi := new(api.LoginApi)
	loginRouter := router.Group("/")
	{
		config.L().Debug("uri: /login, func name: loginApi.Login")
		loginRouter.POST("/login", loginApi.Login)
		//loginRouter.POST("/logout", loginApi.Logout)
	}
}