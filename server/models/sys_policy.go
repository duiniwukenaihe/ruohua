package models

//type SysPolicy struct {
//	Rolename     string `json:"rolename"`
//	ApiInterface string `json:"api_interface"`
//	Method       string `json:"method"`
//}

type SysPolices struct {
	Roles [][]string `json:"roles" binding:"required"`
}
