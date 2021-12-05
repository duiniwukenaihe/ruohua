package mysql

import (
	"fmt"
	"github.com/chenhqchn/ruohua/server/models"
	"github.com/chenhqchn/ruohua/server/utils/config"
	db_util "github.com/chenhqchn/ruohua/server/utils/db"
	"github.com/chenhqchn/ruohua/server/utils/resp"
	"github.com/pkg/errors"
	"time"
)

const (
	getRolesSQL   = `SELECT %s FROM sys_role WHERE deleted = 0;`
	getRoleSQL    = `SELECT %s FROM sys_role WHERE role_name = ? AND deleted = 0;`
	createRoleSQL = `INSERT INTO sys_role (role_name,state,description) VALUES (?,?,?);`
	deleteRoleSQL = `UPDATE sys_role SET deleted = 1 WHERE role_name = ?;`
	updateRoleSQL = `UPDATE sys_role SET %s WHERE role_name = ?;`
)

type RoleDao struct {
}

// 获取角色列表
func (rd *RoleDao) GetRoles(role models.SysRole) (*[]models.SysRole, error) {
	sqlFields, _ := role.FieldToDB()
	db := config.D()

	stmt, err := db.Prepare(fmt.Sprintf(getRolesSQL, sqlFields))
	if err != nil {
		return nil, errors.Wrap(err, "database sql preparation failed")
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, errors.Wrap(err, "data query error")
	}
	defer rows.Close()

	roles := []models.SysRole{}
	for rows.Next() {
		var role models.SysRole
		_, scanFields := role.FieldToDB()
		err := rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "data mapping to field failed")
		}
		roles = append(roles, role)
	}
	return &roles, nil
}

// 获取角色详细信息
func (rd *RoleDao) GetRole(roleInfo models.SysRole) (*models.SysRole, error) {
	sqlFields, scanFields := roleInfo.FieldToDB()
	db := config.D()

	stmt, err := db.Prepare(fmt.Sprintf(getRoleSQL, sqlFields))
	if err != nil {
		return nil, errors.Wrap(err, "Database sql preparation failed")
	}
	defer stmt.Close()

	rows, err := stmt.Query(*roleInfo.RoleName)
	if err != nil {
		return nil, errors.Wrap(err, "Data query failed")
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to map data to field")
		}
	}
	return &roleInfo, nil
}

// 角色创建
func (rd *RoleDao) CreateRole(roleInfo models.SysRole) (*models.SysRole, error) {
	db := config.D()
	tx, err := db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "Transaction failed to start")
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(createRoleSQL)
	if err != nil {
		return nil, errors.Wrap(err, "Database sql preparation failed")
	}
	defer stmt.Close()

	execRes, err := stmt.Exec(*roleInfo.RoleName, *roleInfo.State, *roleInfo.Description)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to excuting create sql")
	}
	tx.Commit()

	realRoleID, _ := execRes.LastInsertId()
	roleInfo.RoleID = &realRoleID
	realTime, _ := time.ParseInLocation("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"), time.Local)
	roleInfo.CreateTime = &realTime
	roleInfo.UpdateTime = &realTime

	deleted := 0
	roleInfo.Deleted = &deleted
	return &roleInfo, nil
}

// 角色删除
func (rd *RoleDao) DeleteRole(roleName string) error {
	db := config.D()
	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "Transaction failed to start")
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(deleteRoleSQL)
	if err != nil {
		return errors.Wrap(err, "Database sql preparation failed")
	}
	defer stmt.Close()

	_, err = stmt.Exec(roleName)
	if err != nil {
		return errors.Wrap(err, "Failed to excuting delete sql ")
	}
	tx.Commit()
	return nil
}

func (rd *RoleDao) UpdateRole(updateData map[string]interface{}, rolename string) error {
	db := config.D()
	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "Transaction failed to start")
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(db_util.BuildCommon(updateRoleSQL, updateData))
	if err != nil {
		return errors.Wrap(err, "Database sql preparation failed")
	}
	defer stmt.Close()

	execRes, err := stmt.Exec(rolename)
	if err != nil {
		fmt.Println(err.Error())
		return errors.Wrap(err, "Failed to excuting update sql")
	}
	tx.Commit()

	affectedRows, _ := execRes.RowsAffected()
	if affectedRows == 0 {
		return resp.RoleNotUpdateError
	}

	return nil
}
