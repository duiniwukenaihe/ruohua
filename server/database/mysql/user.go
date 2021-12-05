package mysql

import (
	"fmt"
	"github.com/chenhqchn/ruohua/server/models"
	"github.com/chenhqchn/ruohua/server/utils/config"
	db_util "github.com/chenhqchn/ruohua/server/utils/db"
	"github.com/chenhqchn/ruohua/server/utils/resp"
	"github.com/pkg/errors"
)

const (
	queryUserSQL  = `SELECT %s FROM sys_user WHERE username = ? AND deleted = 0;`
	updateUserSQL = `UPDATE sys_user set %s WHERE username = ?;`
	postUserSQL   = `INSERT INTO sys_user (dept_id,username,nickname,email,phone,sex,password,state,deleted,login_ip,last_login
					) VALUES (?,?,?,?,?,?,?,?,?,?,?);`
	getUsersSQL   = `SELECT %s FROM sys_user WHERE deleted = 0;`
	deleteUserSQL = `UPDATE sys_user set deleted = 1 WHERE username = ?;`
)

type UserDao struct {
}

// 根据用户名查询用户信息
func (ud *UserDao) GetUser(user models.SysUser) (*models.SysUser, error) {
	sqlFields, scanFields := user.FieldToDB()
	db := config.D()

	stmt, err := db.Prepare(fmt.Sprintf(queryUserSQL, sqlFields))
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	defer stmt.Close()

	rows, err := stmt.Query(user.UserName)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "")
		}
	}
	return &user, nil
}

// 查询所有用户
func (ud *UserDao) GetUsers(user models.SysUser) (*[]models.SysUser, error) {
	sqlFields, _ := user.FieldToDB()
	db := config.D()

	stmt, err := db.Prepare(fmt.Sprintf(getUsersSQL, sqlFields))
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	defer rows.Close()

	users := []models.SysUser{}
	for rows.Next() {
		var user models.SysUser
		_, scanFields := user.FieldToDB()
		err := rows.Scan(scanFields...)
		if err != nil {
			return nil, errors.Wrap(err, "")
		}
		users = append(users, user)
	}
	return &users, nil
}

func (ud *UserDao) UpdateUser(info map[string]interface{}, username string) error {
	// 用户不传任何参数的情况
	if len(info) == 0 {
		return resp.UserNotUpdateError
	}

	db := config.D()
	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "")
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(db_util.BuildCommon(updateUserSQL, info))
	if err != nil {
		return errors.Wrap(err, "")
	}
	defer stmt.Close()

	_, err = stmt.Exec(username)
	if err != nil {
		return errors.Wrap(err, "")
	}
	tx.Commit()
	return nil
}

// 创建用户
func (ud *UserDao) PostUser(user models.SysUser) error {
	db := config.D()
	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "")
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(postUserSQL)
	if err != nil {
		return errors.Wrap(err, "")
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.DeptID, user.UserName, user.NickName, user.Email, user.Phone, user.Sex, user.Password,
		user.State, user.Deleted, user.LoginIP, user.LastLogin)
	if err != nil {
		return errors.Wrap(err, "")
	}
	tx.Commit()

	return nil
}

// 删除用户
func (ud *UserDao) DeleteUserByName(username string) error {
	db := config.D()
	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "")
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(deleteUserSQL)
	if err != nil {
		return errors.Wrap(err, "")
	}
	defer stmt.Close()

	_, err = stmt.Exec(username)
	if err != nil {
		return errors.Wrap(err, "")
	}
	tx.Commit()

	//// 提交事务后，判断修改的行数，如果为0，代表删除的用户不存在
	//rows, _ := result.RowsAffected()
	//if rows == 0 {
	//	return resp.UserNotExistError
	//}
	return nil
}
