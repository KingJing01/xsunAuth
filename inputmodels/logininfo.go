package inputmodels

//user login  from info
type LoginInfo struct {
	UserName string `form:"username"`
	Password string `form:"password"`
	SysID    string `form:"sysId"`
	TenantID int    `form:"tenantId"`
}
