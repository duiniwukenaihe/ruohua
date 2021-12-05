package casbin

import (
	"github.com/chenhqchn/ruohua/server/utils/resp"
	"github.com/pkg/errors"
)

type Policies struct {
}

func (p *Policies) AddPolicesForRole(policies [][]string) error {
	isok, err := Enforcer.AddPolicies(policies)

	if err != nil {
		return errors.Wrap(err, "")
	}
	err = Enforcer.SavePolicy()
	if err != nil {
		return errors.Wrap(err, "")
	}
	if !isok {
		return resp.PoliciesExistError
	}
	return nil
}

func (p *Policies) RemovePolicesForRole(policies [][]string) error {
	isok, err := Enforcer.RemovePolicies(policies)
	if err != nil {
		return errors.Wrap(err, "")
	}
	err = Enforcer.SavePolicy()
	if err != nil {
		return errors.Wrap(err, "")
	}
	if !isok {
		return resp.PoliciesNotExistError
	}
	return nil
}

// 用户添加多个角色
func (p *Policies) AddRolesForUser(username string, roles []string) error {
	isok, err := Enforcer.AddRolesForUser(username, roles)
	if err != nil {
		return errors.Wrap(err, "")
	}
	if !isok {
		return resp.RolesForUserExistError
	}
	return nil
}

// 移除用户的某个角色
func (p *Policies) DeleteRoleForUser(username string, rolename string) error {
	isok, err := Enforcer.DeleteRoleForUser(username, rolename)
	if err != nil {
		return errors.Wrap(err, "")
	}
	if !isok {
		return resp.RolesForUserNotExistError
	}
	return nil
}

// 获取用户所属角色
func (p *Policies) GetRolesForUser(username string) ([]string, error) {
	res := []string{}
	if username == "admin" {
		res = []string{"This is an admin user and has all permissions"}
		return res, nil
	}

	res, err := Enforcer.GetRolesForUser(username)
	if err != nil {
		return []string{}, errors.Wrap(err, "")
	}
	if len(res) == 0 {
		return res, resp.RolesForUserAllNotExistError
	}
	return res, nil
}

func (p *Policies) GetPolicies(rolename string) [][]string {
	res := [][]string{}
	polices := Enforcer.GetPolicy()
	for _, policy := range polices {
		if policy[0] == rolename {
			res = append(res, policy)
		}
	}
	return res
}

func (p *Policies) GetUsersForRole(rolename string) ([]string, error) {
	return Enforcer.GetUsersForRole(rolename)
}