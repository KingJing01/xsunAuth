package inputmodels

//user login  from info
type LoginInfo struct {
	UserName  string `form:"username"`
	Password  string `form:"password"`
	SysID     string `form:"sysId"`
	tokenTime int64  `form:"tokenTime"`
}
