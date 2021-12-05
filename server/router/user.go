package router

import (
	"github.com/chenhqchn/ruohua/server/api"
	"github.com/chenhqchn/ruohua/server/middleware"
	"github.com/chenhqchn/ruohua/server/utils/config"
	"github.com/gin-gonic/gin"
)

func initUserRouter(router *gin.RouterGroup) {
	config.L().Debug("Init the user api group")
	userApi := new(api.UserApi)
	userRouter := router.Group("/users")
	userRouter.Use(middleware.JwtAuth())
	userRouter.Use(middleware.CheckPrivilege())
	{
		userRouter.GET("", userApi.GetUsers)                                                         // 查询所有的用户
		userRouter.GET("/:username", middleware.CheckUserNameIsExist(), userApi.GetUserByName)       // 查询单个用户
		userRouter.POST("", userApi.UserCreate)                                                      // 创建用户
		userRouter.DELETE("/:username", middleware.CheckUserNameIsExist(), userApi.DeleteUserByName) // 删除用户
		userRouter.PUT("/:username", middleware.CheckUserNameIsExist(), userApi.UpdateUser)          // 修改用户所有信息
		userRouter.PATCH("/:username", middleware.CheckUserNameIsExist(), userApi.PartialUpdateUser) // 修改用户部分信息

		userRouter.PATCH("/policies/:username", middleware.CheckUserNameIsExist(), userApi.AddRolesForUser)    // 为用户添加角色
		userRouter.DELETE("/policies/:username", middleware.CheckUserNameIsExist(), userApi.DeleteRoleForuser) // 移除用户某个角色
		userRouter.GET("/policies/:username", middleware.CheckUserNameIsExist(), userApi.GetRolesForUser)      // 获取用户所属角色
	}
}
