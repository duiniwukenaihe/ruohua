package service

import (
	self_password "github.com/chenhqchn/gotools/password"
	self_casbin "github.com/chenhqchn/ruohua/server/utils/casbin"
	self_jwt "github.com/chenhqchn/ruohua/server/utils/jwt"
	"github.com/chenhqchn/ruohua/server/utils/resp"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"time"
)

type LoginService struct {
	userService UserService
}

func (s *LoginService) Login(ctx *gin.Context, username, password string) (string, error) {
	dbRes, err := s.userService.GetUserByUserName(username)
	if err != nil {
		return "", errors.Wrap(err, "")
	}
	if dbRes.UserID == 0 {
		return "", resp.UserNotExistError
	}

	// 用户存在
	// 判断用户是否被禁用
	if dbRes.State == 1 {
		return "", resp.UserStateError
	}

	// 判断密码是否正确
	if !self_password.PasswordCompare(dbRes.Password, password) {
		return "", resp.UserPasswordNotMatchError
	}

	// 用户登录成功，需要更新上次登录时间、登录IP
	updateData := map[string]interface{}{}
	updateData["last_login"] = time.Now().Format("2006-01-02 15:04:05")
	updateData["login_ip"] = ctx.ClientIP()

	//将要更新的数据更新至 DB 中
	err = s.userService.PatchUser(updateData, dbRes.UserName)
	if err != nil {
		return "", resp.UserUpdateError
	}

	// 需要写获取用户对应 role 相关逻辑
	var p self_casbin.Policies
	res, err := p.GetRolesForUser(username)
	// 错误分为两种情况，一、系统原因导致查询失败，二、用户没有关联角色
	if err != nil && err != resp.RolesForUserAllNotExistError {
		// 返回系统原因的错误
		return "", errors.Wrap(err, "")
	}
	if err == resp.RolesForUserAllNotExistError {
		// 将 role 赋值为空，能继续用户登录，只不过登录之后访问接口均提示 403 Forbidden
		res = []string{}
	}
	tokenString, _ := self_jwt.GenToken(username, res)

	return tokenString, nil
}
