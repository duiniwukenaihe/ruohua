package router

import (
	"github.com/chenhqchn/ruohua/server/api"
	"github.com/chenhqchn/ruohua/server/middleware"
	"github.com/gin-gonic/gin"
)

func initRoleRouter(router *gin.RouterGroup) {
	roleApi := new(api.RoleApi)
	roleRouter := router.Group("/roles")
	roleRouter.Use(middleware.JwtAuth())
	roleRouter.Use(middleware.CheckPrivilege())
	{
		roleRouter.GET("", roleApi.GetRoles)                                               // 获取角色列表
		roleRouter.GET("/:rolename", middleware.CheckRoleIsExist(), roleApi.GetRole)       // 查询单个角色详细信息
		roleRouter.POST("", roleApi.CreateRole)                                            // 创建角色
		roleRouter.DELETE("/:rolename", middleware.CheckRoleIsExist(), roleApi.DeleteRole) // 删除角色
		roleRouter.PUT("/:rolename", middleware.CheckRoleIsExist(), roleApi.UpdateRole)    // 修改角色所有信息(描述信息可以不传)
	}

	policyRouter := router.Group("/policies")
	policyRouter.Use(middleware.JwtAuth())
	policyRouter.Use(middleware.CheckPrivilege())
	{
		policyRouter.POST("/:rolename", middleware.CheckRoleIsExist(), roleApi.AddPolies)         // 为角色添加策略
		policyRouter.DELETE("/:rolename", middleware.CheckRoleIsExist(), roleApi.RemovePolies)    // 为角色移除策略
		policyRouter.GET("/:rolename", middleware.CheckRoleIsExist(), roleApi.GetPoliciesForRole) // 获取角色相关联的策略
	}
}
