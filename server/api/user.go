package api

import (
	"github.com/chenhqchn/ruohua/server/models/request"
	"github.com/chenhqchn/ruohua/server/service"
	"github.com/chenhqchn/ruohua/server/utils/resp"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

type UserApi struct {
	userService service.UserService
	roleService service.RoleService
}

func (ua *UserApi) UserCreate(ctx *gin.Context) {
	sysuserBody := request.SysUserBody{}
	// 参数校验，required 类型是否传参
	if err := ctx.ShouldBindBodyWith(&sysuserBody, binding.JSON); err == nil {
		err := ua.userService.PostUser(&sysuserBody)
		if err != nil {
			switch err {
			case resp.UserNameLengthError, resp.UserEmailFormatError, resp.UserPhoneFormatError,
				resp.UserPasswordNotConsistentError, resp.UserPasswordLengthError, resp.UserPasswordComplexError:
				ctx.JSON(http.StatusBadRequest, gin.H{"status": 400, "msg": err.Error()})
				return
			default:
				ctx.JSON(http.StatusInternalServerError, resp.ApiInternalError)
				return
			}
		}
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "msg": "Successfuly created user"})
	} else {
		resp.GetBindingError(ctx, sysuserBody, err)
	}
}

// 查询所有的用户
func (ua *UserApi) GetUsers(ctx *gin.Context) {
	res, err := ua.userService.GetUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, resp.ApiInternalError)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": 200, "msg": "Successfully get user list", "data": res})
}

// 根据用户名查询用户信息
func (ua *UserApi) GetUserByName(ctx *gin.Context) {
	res, err := ua.userService.GetUserByUserName(ctx.Param("username"))
	if err != nil {
		switch err {
		default:
			ctx.JSON(http.StatusInternalServerError, resp.ApiInternalError)
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"status": 200, "msg": "Successfully get user details", "data": res})
}

// 删除用户
func (ua *UserApi) DeleteUserByName(ctx *gin.Context) {
	err := ua.userService.DeleteUserByName(ctx.Param("username"))
	if err != nil {
		switch err {
		case resp.UserDeleteError:
			ctx.JSON(http.StatusBadRequest, gin.H{"status": 400, "msg": err.Error()})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, resp.ApiInternalError)
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"status": 200, "msg": "User deleted successfully"})
}

// 部分更新（更新哪个字段，就传入哪个字段参数）
func (ua *UserApi) PartialUpdateUser(ctx *gin.Context) {
	username := ctx.Param("username")

	updateBody := request.UpdateUserBody{}
	if err := ctx.ShouldBindBodyWith(&updateBody, binding.JSON); err == nil {
		err := ua.userService.PartialUpdateUser(username, &updateBody)
		if err != nil {
			switch err {
			case resp.UserEmailFormatError, resp.UserPhoneFormatError, resp.UserPasswordNotConsistentError,
				resp.UserPasswordLengthError, resp.UserPasswordComplexError, resp.UserPasswordEnterError, resp.UserNotUpdateError:
				ctx.JSON(http.StatusBadRequest, gin.H{"status": 400, "msg": err.Error()})
				return
			default:
				ctx.JSON(http.StatusInternalServerError, resp.ApiInternalError)
				return
			}
		}
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "msg": "User data updated successfully"})
	} else {
		resp.GetBindingError(ctx, updateBody, err)
	}
}

// 全量更新，必须传入所有参数
func (ua *UserApi) UpdateUser(ctx *gin.Context) {
	username := ctx.Param("username")

	updateBody := request.SysUserBody{}
	if err := ctx.ShouldBindBodyWith(&updateBody, binding.JSON); err == nil {
		err := ua.userService.PutUser(username, &updateBody)
		if err != nil {
			switch err {
			case resp.UserEmailFormatError, resp.UserPhoneFormatError, resp.UserPasswordNotConsistentError,
				resp.UserPasswordLengthError, resp.UserPasswordComplexError, resp.UserPasswordEnterError:
				ctx.JSON(http.StatusBadRequest, gin.H{"status": 400, "msg": err.Error()})
				return
			default:
				ctx.JSON(http.StatusInternalServerError, resp.ApiInternalError)
				return
			}
		}
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "msg": "User data updated successfully"})
	} else {
		resp.GetBindingError(ctx, updateBody, err)
	}
}

// 用户添加角色
func (ua *UserApi) AddRolesForUser(ctx *gin.Context) {
	username := ctx.Param("username")
	roles := request.RolesForUserRequest{}
	if err := ctx.ShouldBindBodyWith(&roles, binding.JSON); err == nil {
		err := ua.userService.AddRolesForUser(username, roles.Roles)
		if err != nil {
			switch err {
			case resp.RolesForUserExistError, resp.RoleNotExistError,resp.AdminNotAddRoleError:
				ctx.JSON(http.StatusBadRequest, gin.H{"status": 400, "msg": err.Error()})
				return
			default:
				ctx.JSON(http.StatusInternalServerError, resp.ApiInternalError)
				return
			}
		}
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "msg": "User added roles successfully"})

	} else {
		resp.GetBindingError(ctx, roles, err)
		return
	}
}

// 移除用户单个角色
func (ua *UserApi) DeleteRoleForuser(ctx *gin.Context) {
	username := ctx.Param("username")
	roles := request.RolesForUserRequest{}
	if err := ctx.ShouldBindBodyWith(&roles, binding.JSON); err == nil {
		err := ua.userService.DeleteRoleForuser(username, roles.Roles)
		if err != nil {
			switch err {
			case resp.TooManyRolesDeleteError, resp.RolesForUserNotExistError,
			resp.AdminNotDeleteRoleError,resp.ParameterError:
				ctx.JSON(http.StatusBadRequest, gin.H{"status": 400, "msg": err.Error()})
				return
			default:
				ctx.JSON(http.StatusInternalServerError, resp.ApiInternalError)
				return
			}
		}
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "msg": "User deleted roles successfully"})
	} else {
		resp.GetBindingError(ctx, roles, err)
		return
	}
}

// 获取用户所属角色
func (ua *UserApi) GetRolesForUser(ctx *gin.Context) {
	res, err := ua.userService.GetRolesForUser(ctx.Param("username"))
	if err != nil {
		switch err {
		case resp.RolesForUserAllNotExistError:
			ctx.JSON(http.StatusBadRequest, gin.H{"status": 400, "msg": err.Error()})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, resp.ApiInternalError)
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"status": 200, "msg": "Successfully get roles for user", "data": res})
}
