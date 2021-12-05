package service

import (
	"fmt"
	"github.com/chenhqchn/gotools/password"
	"github.com/chenhqchn/gotools/selfreg"
	"github.com/chenhqchn/ruohua/server/database/mysql"
	"github.com/chenhqchn/ruohua/server/models"
	"github.com/chenhqchn/ruohua/server/models/request"
	self_casbin "github.com/chenhqchn/ruohua/server/utils/casbin"
	"github.com/chenhqchn/ruohua/server/utils/resp"
	"github.com/pkg/errors"
)

type UserService struct {
	userDao      mysql.UserDao
	rolePolicies self_casbin.Policies
}

// 创建用户
// 做参数校验，判断传来当值是否满足要求
func (us *UserService) PostUser(user *request.SysUserBody) error {
	// 判断用户名长度
	if len(*user.UserName) < 5 {
		return resp.UserNameLengthError
	}

	// 判断用户名是否已经存在，保证唯一性
	_, err := us.GetUserByUserName(*user.UserName)
	if err != nil && err != resp.UserNotExistError {
		return errors.Wrap(err, "")
	}

	// 判断 email 是否合法
	if !selfreg.VerifyEmail(*user.Email) {
		return resp.UserEmailFormatError
	}

	// 判断手机号格式是否正确
	if !selfreg.VerifyPhone(*user.Phone) {
		return resp.UserPhoneFormatError
	}

	// 判断两次密码输入是否一致，复杂性是否满足数字、大小写字母、特殊字符至少3种以上
	if isok, msg := selfreg.VerifyPassword(*user.FirstPassword, *user.SecondPassword); !isok {
		if msg == resp.UserPasswordNotConsistentError.Error() {
			return resp.UserPasswordNotConsistentError
		}
		if msg == resp.UserPasswordLengthError.Error() {
			return resp.UserPasswordLengthError
		}
		if msg == resp.UserPasswordComplexError.Error() {
			return resp.UserPasswordComplexError
		}
	}

	// 判断性别所属值是否正确, 若传入其它值，将默认设置为男性（1）
	var sex int
	if *user.Sex != 0 && *user.Sex != 1 {
		sex = 1
	} else {
		sex = *user.Sex
	}

	// 判断状态所属值是否正确, 若传入其它值，将默认设置为0（正常）
	var state int
	if *user.State != 0 && *user.State != 1 {
		state = 0
	} else {
		state = *user.State
	}

	// 首次创建用户，last_login，login_ip 均不存在，所以此处不写入数据，由数据库层面生成默认值
	prepareData := models.SysUser{
		DeptID:   *user.DeptID,
		UserName: *user.UserName,
		NickName: *user.NickName,
		Email:    *user.Email,
		Phone:    *user.Phone,
		Sex:      sex,
		Password: password.EncryptPassword(*user.FirstPassword),
		State:    state,
		LoginIP:  "",
		Deleted:  0,
	}
	return us.userDao.PostUser(prepareData)
}

// 通过用户名查询用户信息
func (us *UserService) GetUserByUserName(name string) (*models.SysUser, error) {
	user := models.SysUser{}
	user.UserName = name
	dbRes, err := us.userDao.GetUser(user)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	return dbRes, nil
}

// 查询所有用户
func (us *UserService) GetUsers() (*[]models.SysUser, error) {
	var user = models.SysUser{}
	return us.userDao.GetUsers(user)
}

// 部分更新，目前只用于 程序内部使用，界面、调用API 均使用全量更新方法
// 暂时由内部调用
func (us *UserService) PatchUser(info map[string]interface{}, username string) error {
	return us.userDao.UpdateUser(info, username)
}

// 暴露给 api 层使用
func (us *UserService) PartialUpdateUser(username string, body *request.UpdateUserBody) error {
	updateFields := map[string]interface{}{}
	if body.DeptID != nil {
		updateFields["dept_id"] = *body.DeptID
	}
	if body.NickName != nil {
		updateFields["nickname"] = *body.NickName
	}
	if body.Email != nil {
		if !selfreg.VerifyEmail(*body.Email) {
			return resp.UserEmailFormatError
		}
		updateFields["email"] = *body.Email
	}
	if body.Phone != nil {
		if !selfreg.VerifyPhone(*body.Phone) {
			return resp.UserPhoneFormatError
		}
		updateFields["phone"] = *body.Phone
	}
	if body.Sex != nil {
		if *body.Sex != 0 && *body.Sex != 1 {
			updateFields["sex"] = 1
		} else {
			updateFields["sex"] = *body.Sex
		}
	}

	// 更新密码，判断密码是否满足复杂性要求或者是否一致
	if body.FirstPassword != nil && body.SecondPassword != nil {
		if isok, msg := selfreg.VerifyPassword(*body.FirstPassword, *body.SecondPassword); !isok {
			if msg == resp.UserPasswordNotConsistentError.Error() {
				return resp.UserPasswordNotConsistentError
			}
			if msg == resp.UserPasswordLengthError.Error() {
				return resp.UserPasswordLengthError
			}
			if msg == resp.UserPasswordComplexError.Error() {
				return resp.UserPasswordComplexError
			}
		}
		updateFields["password"] = password.EncryptPassword(*body.FirstPassword)
	}

	// 用若需要更新密码，但传入其中一个参数，first_password second_password
	if body.FirstPassword != nil || body.SecondPassword != nil {
		return resp.UserPasswordEnterError
	}

	if body.State != nil {
		if *body.State != 0 && *body.State != 1 {
			updateFields["state"] = 0
		} else {
			updateFields["state"] = *body.State
		}
	}

	return us.userDao.UpdateUser(updateFields, username)
}

func (us *UserService) PutUser(username string, body *request.SysUserBody) error {
	updateFields := map[string]interface{}{}
	updateFields["dept_id"] = *body.DeptID
	if username != "admin" {
		updateFields["nickname"] = *body.NickName
	}
	if !selfreg.VerifyEmail(*body.Email) {
		return resp.UserEmailFormatError
	}
	updateFields["email"] = *body.Email
	if !selfreg.VerifyPhone(*body.Phone) {
		return resp.UserPhoneFormatError
	}
	updateFields["phone"] = *body.Phone

	if *body.Sex != 0 && *body.Sex != 1 {
		updateFields["sex"] = 1
	} else {
		updateFields["sex"] = *body.Sex
	}
	if body.FirstPassword != nil && body.SecondPassword != nil {
		if isok, msg := selfreg.VerifyPassword(*body.FirstPassword, *body.SecondPassword); !isok {
			if msg == resp.UserPasswordNotConsistentError.Error() {
				return resp.UserPasswordNotConsistentError
			}
			if msg == resp.UserPasswordLengthError.Error() {
				return resp.UserPasswordNotConsistentError
			}
			if msg == resp.UserPasswordComplexError.Error() {
				return resp.UserPasswordNotConsistentError
			}
		}
		updateFields["password"] = password.EncryptPassword(*body.FirstPassword)
	}
	// 被修改用户为 admin，不允许被禁用，即使传入任何值，也保持 state 为0
	if username == "admin" {
		updateFields["state"] = 0
	} else if *body.State != 0 && *body.State != 1 {
		updateFields["state"] = 0
	} else {
		updateFields["state"] = *body.State
	}
	return us.userDao.UpdateUser(updateFields, username)
}

// 删除用户
func (us *UserService) DeleteUserByName(username string) error {
	if username == "admin" {
		return resp.UserDeleteError
	}

	return us.userDao.DeleteUserByName(username)
}

// 为用户添加角色
func (us *UserService) AddRolesForUser(username string, roles []string) error {
	if username == "admin" {
		return resp.AdminNotAddRoleError
	}

	addData := []string{}
	var ra RoleService
	for _, role := range roles {
		dbData, _, err := ra.GetRole(role)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to get the rolename: %s\n", role))
		}

		if dbData.RoleID == nil {
			return resp.RoleNotExistError
		}
		addData = append(addData, role)
	}
	return us.rolePolicies.AddRolesForUser(username, addData)
}

func (us *UserService) DeleteRoleForuser(username string, roles []string) error {
	if username == "admin" {
		return resp.AdminNotDeleteRoleError
	}
	// 空参数请求，"roles": []
	if len(roles) == 0 {
		return resp.ParameterError
	}
	// 这个接口只能传一个角色名，
	if len(roles) > 1 {
		return resp.TooManyRolesDeleteError
	}

	return us.rolePolicies.DeleteRoleForUser(username, roles[0])
}

func (us *UserService) GetRolesForUser(username string) ([]string, error) {
	return us.rolePolicies.GetRolesForUser(username)
}
