package router

import (
	"github.com/chenhqchn/ruohua/server/utils/config"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	config.L().Debug("Start to initing gin engine")
	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	config.L().Debug("Binding routing group")
	apiRouter := router.Group("/api")
	{
		initLoginRouter(apiRouter)
		initUserRouter(apiRouter)
		initRoleRouter(apiRouter)
	}
	return router
}
