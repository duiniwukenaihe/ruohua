package service

import (
	"github.com/chenhqchn/ruohua/server/database/mysql"
	"github.com/chenhqchn/ruohua/server/models"
	self_casbin "github.com/chenhqchn/ruohua/server/utils/casbin"
	"github.com/chenhqchn/ruohua/server/utils/resp"
	"github.com/pkg/errors"
)

type RoleService struct {
	roleDao      mysql.RoleDao
	rolePolicies self_casbin.Policies
}

// 获取角色列表
func (rs *RoleService) GetRoles() (*[]models.SysRole, error) {
	var role models.SysRole
	return rs.roleDao.GetRoles(role)
}

// 获取单个角色详细信息
func (rs *RoleService) GetRole(roleName string) (*models.SysRole, []string, error) {
	roleInfo := models.SysRole{}
	roleInfo.RoleName = &roleName
	dbData, err := rs.roleDao.GetRole(roleInfo)
	if err != nil {
		return nil, nil, errors.Wrap(err, "")
	}

	// 获取角色绑定的用户
	res, err := rs.rolePolicies.GetUsersForRole(roleName)
	if err != nil {
		return nil, nil, errors.Wrap(err, "")
	}
	return dbData, res, nil
}

func (rs *RoleService) CreateRole(role models.SysRole) (*models.SysRole, error) {
	if *role.RoleName == "" {
		return nil, resp.RoleEmptyError
	}
	dbRes, _, err := rs.GetRole(*role.RoleName)
	if err != nil {
		return nil, err
	}

	if dbRes != nil && dbRes.RoleID != nil {
		return nil, resp.RoleExistError
	}

	var state int
	if *role.State != 0 && *role.State != 1 {
		state = 0
	} else {
		state = *role.State
	}

	var description string
	if role.Description == nil {
		description = ""
	} else {
		description = *role.Description
	}

	roleInfoToDB := models.SysRole{
		RoleName:    role.RoleName,
		State:       &state,
		Description: &description,
	}

	return rs.roleDao.CreateRole(roleInfoToDB)
}

func (rs *RoleService) DeleteRole(roleName string) error {
	// 判断此角色是否存在用户绑定
	res, err := rs.rolePolicies.GetUsersForRole(roleName)
	if err != nil {
		return errors.Wrap(err, "")
	}
	if len(res) > 0 {
		return resp.RoleNotDeleteError
	}
	// 先查询与此角色相关的策略
	var policesObj models.SysPolices
	policesObj.Roles = rs.GetPolicies(roleName)

	if len(policesObj.Roles) != 0 {
		// 再删除与此角色相关的所有策略
		err = rs.RemovePolicies(roleName, policesObj)
		if err != nil {
			return errors.Wrap(err, "")
		}
	}

	return rs.roleDao.DeleteRole(roleName)
}

func (rs *RoleService) UpdateRole(role models.SysRole, roleName string) error {
	if *role.State != 0 && *role.State != 1 {
		return resp.RoleStateError
	}

	if role.Description == nil {
		return resp.RoleDescriError
	}

	updateFields := map[string]interface{}{}
	updateFields["state"] = *role.State
	updateFields["description"] = *role.Description

	return rs.roleDao.UpdateRole(updateFields, roleName)
}

// 为角色添加策略
func (rs *RoleService) AddPolicies(rolename string, polices models.SysPolices) error {
	if len(polices.Roles) == 0 {
		return resp.ParameterError
	}
	rolesSlice := [][]string{}
	for _, value := range polices.Roles {
		if len(value) != 2 {
			return resp.PoliciesAddParaError
		}
		value = append(value, "")
		copy(value[1:], value[0:])
		value[0] = rolename
		rolesSlice = append(rolesSlice, value)
	}
	return rs.rolePolicies.AddPolicesForRole(rolesSlice)
}

// 从角色移除策略
func (rs *RoleService) RemovePolicies(rolename string, polices models.SysPolices) error {
	rolesSlice := [][]string{}
	for _, value := range polices.Roles {
		value = append(value, "")
		copy(value[1:], value[0:])
		value[0] = rolename
		rolesSlice = append(rolesSlice, value)
	}
	return rs.rolePolicies.RemovePolicesForRole(rolesSlice)
}

func (rs *RoleService) GetPolicies(rolename string) [][]string {
	return rs.rolePolicies.GetPolicies(rolename)
}
