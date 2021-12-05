package request

type RolesForUserRequest struct {
	Roles []string `json:"roles" binding:"required"`
}
