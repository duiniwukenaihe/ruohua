package api

import (
	"github.com/chenhqchn/ruohua/server/models"
	"github.com/chenhqchn/ruohua/server/service"
	"github.com/chenhqchn/ruohua/server/utils/config"
	"github.com/chenhqchn/ruohua/server/utils/resp"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

type RoleApi struct {
	roleService service.RoleService
}

// 获取角色列表
func (ra *RoleApi) GetRoles(ctx *gin.Context) {
	res, err := ra.roleService.GetRoles()
	if err != nil {
		config.L().Error(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": 500, "msg": "Internal error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": 200, "data": res, "msg": "Get roles successful"})
}

// 根据角色名获取角色详细信息
func (ra *RoleApi) GetRole(ctx *gin.Context) {
	roleName := ctx.Param("rolename")
	roleRes, roleUsers, err := ra.roleService.GetRole(roleName)
	if err != nil {
		switch err {
		default:
			ctx.JSON(http.StatusInternalServerError, resp.ApiInternalError)
			return
		}
	}

	res := map[string]interface{}{}
	res["role"] = roleRes
	res["users"] = roleUsers
	ctx.JSON(http.StatusOK, gin.H{"status": 200, "data": res, "msg": "Successfully get role details"})
}

//创建角色
func (ra *RoleApi) CreateRole(ctx *gin.Context) {
	role := models.SysRole{}
	if err := ctx.ShouldBindBodyWith(&role, binding.JSON); err == nil {
		res, err := ra.roleService.CreateRole(role)
		if err != nil {
			switch err {
			case resp.RoleExistError:
				ctx.JSON(http.StatusConflict, gin.H{"status": 409, "msg": err.Error()})
				return
			case resp.RoleEmptyError:
				ctx.JSON(http.StatusBadRequest, gin.H{"status": 400, "msg": err.Error()})
				return
			default:
				ctx.JSON(http.StatusInternalServerError, resp.ApiInternalError)
				return
			}
		}
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "data": res, "msg": "Role created successfully"})
	} else {
		config.L().Warnf("Parameter binding failed: \n%s", err.Error())
		resp.GetBindingError(ctx, role, err)
	}
}

// 删除角色
func (ra *RoleApi) DeleteRole(ctx *gin.Context) {
	rolename := ctx.Param("rolename")
	err := ra.roleService.DeleteRole(rolename)
	if err != nil {
		switch err {
		case resp.RoleNotDeleteError, resp.PoliciesNotExistError:
			ctx.JSON(http.StatusBadRequest, gin.H{"status": 400, "msg": err.Error()})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, resp.ApiInternalError)
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"status": 200, "msg": "Role deleted successfully"})
}

// 更新角色
func (ra *RoleApi) UpdateRole(ctx *gin.Context) {
	roleName := ctx.Param("rolename")
	role := models.SysRole{}
	if err := ctx.ShouldBindBodyWith(&role, binding.JSON); err == nil {
		err = ra.roleService.UpdateRole(role, roleName)
		if err != nil {
			config.L().Error(err.Error())
			switch err {
			case resp.RoleStateError, resp.RoleDescriError:
				ctx.JSON(http.StatusBadRequest, gin.H{"status": 400, "msg": err.Error()})
				return
			case resp.RoleNotUpdateError:
				ctx.JSON(http.StatusAccepted, gin.H{"status": 202, "msg": err.Error()})
				return
			default:
				ctx.JSON(http.StatusInternalServerError, resp.ApiInternalError)
				return
			}
		}
		config.L().Infof("Role: %s updated successfully", *role.RoleName)
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "msg": "Role updated successfully"})
		return
	} else {
		config.L().Warnf("Parameter binding failed: \n%s", err.Error())
		resp.GetBindingError(ctx, role, err)
	}
}

// 为 角色添加策略
func (ra *RoleApi) AddPolies(ctx *gin.Context) {
	rolename := ctx.Param("rolename")
	policies := models.SysPolices{}
	if err := ctx.ShouldBindBodyWith(&policies, binding.JSON); err == nil {
		err := ra.roleService.AddPolicies(rolename, policies)
		if err != nil {
			config.L().Error(err)
			switch err {
			case resp.PoliciesExistError:
				ctx.JSON(http.StatusConflict, gin.H{"status": 409, "msg": err.Error()})
				return
			case resp.ParameterError, resp.PoliciesAddParaError:
				ctx.JSON(http.StatusBadRequest, gin.H{"status": 400, "msg": err.Error()})
				return
			default:
				ctx.JSON(http.StatusInternalServerError, resp.ApiInternalError)
				return
			}
		}
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "msg": "Successfully add policies"})
	} else {
		resp.GetBindingError(ctx, policies, err)
		return
	}
}

// 为角色移除策略
func (ra *RoleApi) RemovePolies(ctx *gin.Context) {
	rolename := ctx.Param("rolename")
	policies := models.SysPolices{}
	if err := ctx.ShouldBindBodyWith(&policies, binding.JSON); err == nil {
		err := ra.roleService.RemovePolicies(rolename, policies)
		if err != nil {
			config.L().Error(err)
			switch err {
			case resp.PoliciesNotExistError:
				ctx.JSON(http.StatusBadRequest, gin.H{"status": 400, "msg": err.Error()})
				return
			default:
				ctx.JSON(http.StatusInternalServerError, resp.ApiInternalError)
				return
			}
		}
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "msg": "Successfully remove policies"})
	} else {
		resp.GetBindingError(ctx, policies, err)
		return
	}
}

func (ra *RoleApi) GetPoliciesForRole(ctx *gin.Context) {
	res := ra.roleService.GetPolicies(ctx.Param("rolename"))
	ctx.JSON(http.StatusOK, gin.H{"status": 200, "msg": "Successfully get policies for role", "data": res})
}
