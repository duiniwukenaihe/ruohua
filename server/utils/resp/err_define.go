package resp

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var (
	UserNameLengthError            = errors.New("[user] Username length is too short(length >= 5)")
	UserNotExistError              = errors.New("[user] User does not exist")
	UserEmailFormatError           = errors.New("[user] Email format error, please check again")
	UserPhoneFormatError           = errors.New("[user] Phone format error, please check again")
	UserPasswordNotConsistentError = errors.New("[user] The two passwords entered are inconsistent")
	UserPasswordNotMatchError      = errors.New("[user] Username and password do not match")
	UserPasswordLengthError        = errors.New("[user] Password length is too short(legnth > 6)")
	UserPasswordComplexError       = errors.New("[user] The password is too simple(At least 3 types of uppercase and " +
		"lowercase letters, numbers, and special characters)")
	UserPasswordEnterError  = errors.New("[user] Please enter your password twice")
	UserDeleteError         = errors.New("[user] Default user: admin cannot be deleted")
	UserStateError          = errors.New("[user] Login failed, user is disabled")
	UserUpdateError         = errors.New("[user] Failed to update user info")
	UserLoginParaError      = errors.New("[user] Login failed, Please enter username and password")
	UserNotUpdateError      = errors.New("[user]: There are no info can be updated")
	AdminNotDeleteRoleError = errors.New("[user] Admin are not allowed to delete roles")
	AdminNotAddRoleError    = errors.New("[user] Admin are not allowed to add roles")

	RoleExistError               = errors.New("[role]: This role already exists")
	RoleStateError               = errors.New("[role]: The state must be 0 or 1")
	RoleNotExistError            = errors.New("[role]: This role does not exist")
	RoleDescriError              = errors.New("[role]: The parameter description must be specified")
	TooManyRolesDeleteError      = errors.New("[role]: this api interface can be provided only one rolename")
	RoleNotUpdateError           = errors.New("[role]: There are no info can be updated")
	RolesForUserExistError       = errors.New("[role] User in this role already exists")
	RolesForUserNotExistError    = errors.New("[role] User in this role does not exist")
	RolesForUserAllNotExistError = errors.New("[role] User does not have any role")
	RoleEmptyError               = errors.New("[role] Rolename can't be empty")
	RoleNotDeleteError = errors.New("[role] This role is already associated with multiple users, please unbind it first")

	PoliciesExistError    = errors.New("[policies] Policies already exists")
	PoliciesNotExistError = errors.New("[policies] Policies does not exist")
	PoliciesAddParaError = errors.New("[policies] Please enter the parameters in the following format," +
		" [[\"api_interface\",\"method\"],...]")

	ApiInternalError = gin.H{"staus": 500, "msg": "Internal error!"}
	ParameterError   = errors.New("[all]: This parameter cannot be passed a null value")
)
