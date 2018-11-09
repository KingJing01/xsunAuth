package controllers

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
	inputModel "xsunAuth/inputmodels"
	"xsunAuth/models"
	"xsunAuth/tools"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	jwt "github.com/dgrijalva/jwt-go"
)

type CreateTenantInput struct {
	Id                   int
	ConnectionString     string
	CreationTime         time.Time
	CreatorUserId        int64
	DeleterUserId        int64
	DeletionTime         time.Time
	EditionId            int64
	IsActive             bool
	IsDeleted            bool
	LastModificationTime time.Time
	LastModifierUserId   int64
	Name                 string
	TenancyName          string
	EmailAddress         string
}

type CreateRoleInput struct {
	Id                   int
	ConcurrencyStamp     string
	CreationTime         time.Time
	CreatorUserId        int64
	DeleterUserId        int64
	DeletionTime         time.Time
	DisplayName          string
	IsDefault            uint64
	IsDeleted            uint64
	IsStatic             uint64
	LastModificationTime time.Time
	LastModifierUserId   int64
	Name                 string
	NormalizedName       string
	TenantId             int
	Description          string
	Permissions          []models.Permission
}

type TenantPermissionsInput struct {
	Tenantid    int
	Permissions []int64
}
type UserRolesInput struct {
	UserId int64
	Roles  []int
}

type AuthorityManageController struct {
	beego.Controller
}

type OperResult struct {
	Result  bool
	Message string
}

type LoginResult struct {
	OperResult
	Token string
}

type UserInfo struct {
	OperResult
	UserId       int64
	Name         string
	UserName     string
	Gender       string
	Age          string
	Avatar       string
	EmailAddress string
	PhoneNumber  string
	Permissions  []string
}
type TenantInfo struct {
	OperResult
	TenantId    int
	IsDelete    bool
	Name        string
	TenancyName string
	CreateTime  string
	Roles       []RoleInfo
}

type RoleInfo struct {
	RoleId      int
	Name        string
	DisplayName string
}

const (
	SecretKey = "sfljdsfjsljdslfdsfsdfjdsf"
)

func (c *AuthorityManageController) URLMapping() {
	c.Mapping("all", c.GetTenantList)
	c.Mapping("GetTenant", c.GetTenant)
	c.Mapping("AuthorityError", c.AuthorityError)
}

func (c *AuthorityManageController) Options() {
	c.Data["json"] = map[string]interface{}{"status": 200, "message": "ok", "moreinfo": ""}
	c.ServeJSON()
}

// @router /all [get]
func (tc *AuthorityManageController) GetTenantList() {
	token := tc.Ctx.Request.Header.Get("Authorization")
	result := OperResult{}
	ok, _, err := tools.CheckAuthority(token, "Platform.Tenant")
	if !ok {
		result.Result = false
		result.Message = err.Error()
		tc.Data["json"] = result
		tc.ServeJSON()
		return
	}
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	// fields: col1,col2,entity.col3
	if v := tc.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := tc.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := tc.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := tc.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := tc.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := tc.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				tc.Data["json"] = errors.New("Error: invalid query key/value pair")
				tc.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllTenant(query, fields, sortby, order, offset, limit)
	if err != nil {
		tc.Data["json"] = err.Error()
	} else {
		tc.Data["json"] = l
	}
	tc.ServeJSON()
}

// @router /GetTenant [get]
func (tc *AuthorityManageController) GetTenant() {
	token := tc.Ctx.Request.Header.Get("Authorization")
	result := OperResult{}
	if ok, _, err := tools.CheckAuthority(token, "Platform.Tenant"); !ok {
		result.Result = false
		result.Message = err.Error()
		tc.Data["json"] = result
		tc.ServeJSON()
		return
	}

	var sid string
	var name string
	var id int

	sid = tc.GetString("id")
	name = tc.GetString("name")
	if sid == "" && name == "" {
		result.Result = false
		result.Message = "检索条件不能为空"
		tc.Data["json"] = result
		tc.ServeJSON()
		return
	}

	if sid != "" {
		tmpid, error := strconv.Atoi(sid)
		if error != nil {
			result.Result = false
			result.Message = "字符串转换成整数失败"
			tc.Data["json"] = result
			tc.ServeJSON()
			return
		} else {
			id = tmpid
		}
	}

	v, err := models.GetTenant(id, name)
	if err != nil {
		result.Result = false
		result.Message = err.Error()
		tc.Data["json"] = result
	} else {
		tc.Data["json"] = v
	}
	tc.ServeJSON()
	//var req []string
	//json.Unmarshal(([]byte)(reqData), &req)
}

// @router /AuthorityError [get]
func (tc *AuthorityManageController) AuthorityError() {
	result := OperResult{}
	result.Result = false
	result.Message = "未登录"
	tc.Data["json"] = result
	tc.ServeJSON()

}

// @Title Login
// @Description 登入接口
// @Param   key     path    string  true        "The email for login"
// @Success 200 {object} controllers.LoginResult
// @Failure 400 Invalid email supplied
// @Failure 404 User not found
// @router /Login [post]
func (tc *AuthorityManageController) Login() {
	l := &inputModel.LoginInfo{}
	lresult := LoginResult{}
	json.Unmarshal(tc.Ctx.Input.RequestBody, l)
	valid := validation.Validation{}
	resultUserName := valid.Required(l.UserName, "username").Message("请输入用户名")
	if resultUserName.Ok == false {
		lresult.Result = false
		lresult.Message = resultUserName.Error.Message
		tc.Data["json"] = lresult
		tc.ServeJSON()
		return
	}
	resultPass := valid.Required(l.Password, "password").Message("请输入密码")
	if resultPass.Ok == false {
		lresult.Result = false
		lresult.Message = resultPass.Error.Message
		tc.Data["json"] = lresult
		tc.ServeJSON()
		return
	}
	resultSysID := valid.Required(l.SysID, "sysId").Message("系统号不能为空")
	if resultSysID.Ok == false {
		lresult.Result = false
		lresult.Message = resultSysID.Error.Message
		tc.Data["json"] = lresult
		tc.ServeJSON()
		return
	}

	result, user, err := models.LoginCheck(l.TenantID, l.UserName, l.Password, l.SysID)
	respmessage := ""
	if result == false {
		if err == nil {
			respmessage = "用户名和密码不匹配，重新登陆"
		} else {
			respmessage = err.Error()
		}
		lresult.Result = false
		lresult.Message = respmessage
		tc.Data["json"] = lresult
		tc.ServeJSON()
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["jti"] = user.Id
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(10)).Unix()
	claims["iat"] = time.Now().Unix()
	token.Claims = claims
	tokenString, err := token.SignedString([]byte(SecretKey))

	lresult.Result = true
	lresult.Token = tokenString
	tc.Data["json"] = lresult

	tc.ServeJSON()
}

//获取用户信息
// @Title GetUserInfo
// @Description 根据TOKEN获取用户信息
// @Param   Authorization     header    string  true        "Token信息"
// @Success 200 {object} controllers.UserInfo
// @Failure 400 Invalid email supplied
// @Failure 404 User not found
// @router /GetUserInfo [get]
func (tc *AuthorityManageController) GetUserInfo() {
	token := tc.Ctx.Request.Header.Get("Authorization")
	oResult := UserInfo{}
	ok, claims, err := tools.CheckLogin(token)
	if !ok {
		oResult.Result = false
		oResult.Message = err.Error()
		tc.Data["json"] = oResult
		tc.ServeJSON()
		return
	}
	tmp := strconv.FormatFloat(claims["jti"].(float64), 'f', -1, 64)
	userid, _ := strconv.ParseInt(tmp, 10, 64)
	u, _ := models.GetUserById(userid)
	if u != nil {
		oResult.UserName = u.UserName
		oResult.Name = u.Name
		oResult.EmailAddress = u.EmailAddress
		oResult.PhoneNumber = u.PhoneNumber
	}
	permissions, _ := models.GetPermissionByUser(userid)
	var arrPermission []string
	for _, v := range permissions {
		arrPermission = append(arrPermission, v.Name)
	}
	oResult.Permissions = arrPermission
	oResult.Result = true

	tc.Data["json"] = oResult
	tc.ServeJSON()
}

//登出
// @Title Logout
// @Description 登出
// @router /Logout [post]
func (tc *AuthorityManageController) Logout() {
	// token := tc.Ctx.Request.Header.Get("Authorization")
	// result := OperResult{}
	// ok, _, err := tools.CheckAuthority(token, "Logout")
	// if !ok {
	// 	result.Result = false
	// 	result.Message = err.Error()
	// 	tc.Data["json"] = result
	// 	tc.ServeJSON()
	// 	return
	// }

	lresult := LoginResult{}
	lresult.Result = true
	lresult.Token = ""
	lresult.Message = "登出成功"
	tc.Data["json"] = lresult
	tc.ServeJSON()
}

//新增租户
// @Title CreateTenant
// @Description 新增租户
// @Param   Authorization     header    string  true        "Token信息"
// @Param   Name     body     string  true        "租户名"
// @Param   TenancyName     body     string  true        "租户中文名"
// @Param   EmailAddress     body     string  true        "管理员的EMAIL"
// @Success 200 {object} controllers.TenantInfo
// @Failure 400 Invalid email supplied
// @Failure 404 User not found
// @router /CreateTenant [post]
func (tc *AuthorityManageController) CreateTenant() {
	token := tc.Ctx.Request.Header.Get("Authorization")
	oResult := TenantInfo{}
	ok, _, err := tools.CheckAuthority(token, "Platform.Tenant")
	if !ok {
		oResult.Result = false
		oResult.Message = err.Error()
		tc.Data["json"] = oResult
		tc.ServeJSON()
		return
	}

	tinput := &CreateTenantInput{}
	t := &models.Tenant{}
	json.Unmarshal(tc.Ctx.Input.RequestBody, tinput)
	//验证租户是否已存在

	sid := tinput.Id
	name := tinput.Name
	if name == "" {
		oResult.Result = false
		oResult.Message = "租户名不能为空"
		tc.Data["json"] = oResult
		tc.ServeJSON()
		return
	}

	num, _ := models.GetCount(sid, name)
	if num > 0 {
		oResult.Result = false
		oResult.Message = "此租户已存在，不可新增"
		tc.Data["json"] = oResult
		tc.ServeJSON()
		return
	} else {
		t.CreationTime = time.Now()
		t.IsActive = true
		t.IsDeleted = false
		t.Name = tinput.Name
		t.TenancyName = tinput.TenancyName
		result, nt, _ := models.CreateTenant(t, tinput.EmailAddress)
		oResult.TenantId = nt.Id
		oResult.IsDelete = nt.IsDeleted
		oResult.Name = nt.Name
		oResult.TenancyName = nt.TenancyName
		oResult.CreateTime = nt.CreationTime.Format("2006-01-02 15:04:05")
		oResult.Result = result
		oResult.Message = "OK"
	}

	tc.Data["json"] = oResult
	tc.ServeJSON()
}

//更新租户
// @Title UpdateTenant
// @Description 更新租户
// @Param   Authorization     header    string  true        "Token信息"
// @Param   Name     body     string  true        "租户名"
// @Param   TenancyName     body     string  true        "租户中文名"
// @Success 200 {object} controllers.OperResult
// @Failure 400 Invalid email supplied
// @Failure 404 User not found
// @router /UpdateTenant [post]
func (tc *AuthorityManageController) UpdateTenant() {
	token := tc.Ctx.Request.Header.Get("Authorization")
	oResult := OperResult{}
	ok, _, err := tools.CheckAuthority(token, "Platform.Tenant")
	if !ok {
		oResult.Result = false
		oResult.Message = err.Error()
		tc.Data["json"] = oResult
		tc.ServeJSON()
		return
	}

	t := &models.Tenant{}
	json.Unmarshal(tc.Ctx.Input.RequestBody, t)
	//验证租户是否已存在
	sid := t.Id
	name := t.Name

	if sid == 0 {
		oResult.Result = false
		oResult.Message = "租户不存在，不能更新"
		tc.Data["json"] = oResult
		tc.ServeJSON()
		return
	}

	if name == "" {
		oResult.Result = false
		oResult.Message = "租户名不能为空"
		tc.Data["json"] = oResult
		tc.ServeJSON()
		return
	}

	v, err := models.GetTenantByName(name)
	if err != nil {
		oResult.Result = false
		oResult.Message = err.Error()
		tc.Data["json"] = oResult
		tc.ServeJSON()
		return
	} else {
		if v != nil {
			if v.Id != sid {
				oResult.Result = false
				oResult.Message = "存在相同的租户名，不可修改"
				tc.Data["json"] = oResult
				tc.ServeJSON()
				return
			}
		}

	}
	v, err = models.GetTenantById(sid)
	if err != nil {
		oResult.Result = false
		oResult.Message = err.Error()
		tc.Data["json"] = oResult
		tc.ServeJSON()
		return
	}
	v.Name = t.Name
	v.TenancyName = t.TenancyName
	// t.IsActive = true
	// t.IsDeleted = false
	err = models.UpdateTenantById(v)

	if err != nil {
		oResult.Message = err.Error()
		oResult.Result = false
		tc.Data["json"] = oResult
		tc.ServeJSON()
		return
	}
	oResult.Result = true
	tc.Data["json"] = oResult
	tc.ServeJSON()

}

//删除租户--软删除
// @Title DeleteTenant
// @Description 删除租户
// @Param   Authorization     header    string  true        "Token信息"
// @Param   Id     body     int  true        "租户ID"
// @Success 200 {object} controllers.OperResult
// @Failure 400 Invalid email supplied
// @Failure 404 User not found
// @router /DeleteTenant [post]
func (tc *AuthorityManageController) DeleteTenant() {
	token := tc.Ctx.Request.Header.Get("Authorization")
	oResult := OperResult{}
	ok, _, err := tools.CheckAuthority(token, "Platform.Tenant")
	if !ok {
		oResult.Result = false
		oResult.Message = err.Error()
		tc.Data["json"] = oResult
		tc.ServeJSON()
		return
	}

	t := &models.Tenant{}
	json.Unmarshal(tc.Ctx.Input.RequestBody, t)
	sid := t.Id
	if sid == 0 {
		oResult.Message = "租户名不能为空,不可删除"
		oResult.Result = false
		tc.Data["json"] = oResult
		tc.ServeJSON()
		return
	}
	v, err := models.GetTenantById(sid)
	if err != nil {
		oResult.Message = err.Error()
		oResult.Result = false
		tc.Data["json"] = oResult
		tc.ServeJSON()
		return
	}
	v.IsDeleted = true
	err = models.UpdateTenantById(v)

	if err != nil {
		oResult.Message = err.Error()
		oResult.Result = false
		tc.Data["json"] = oResult
		tc.ServeJSON()
		return
	}

	oResult.Result = true
	tc.Data["json"] = oResult
	tc.ServeJSON()
}

//设置租户的基础权限
// @Title SetTenantPermission
// @Description 设置租户的基础权限
// @Param   Authorization     header    string  true        "Token信息"
// @Param   Tenantid     body     int  true        "租户ID"
// @Param   Permissions     body     []int  true        "权限LIST"
// @Success 200 {object} controllers.OperResult
// @Failure 400 Invalid email supplied
// @Failure 404 User not found
// @router /SetTenantPermission [post]
func (tc *AuthorityManageController) SetTenantPermission() {
	token := tc.Ctx.Request.Header.Get("Authorization")
	oResult := OperResult{}
	ok, _, err := tools.CheckAuthority(token, "Platform.Tenant")
	if !ok {
		oResult.Result = false
		oResult.Message = err.Error()
		tc.Data["json"] = oResult
		tc.ServeJSON()
		return
	}

	tp := &TenantPermissionsInput{}
	json.Unmarshal(tc.Ctx.Input.RequestBody, tp)
	if tp.Tenantid == 0 {
		oResult.Message = "租户不存在，不能设置权限"
		oResult.Result = false
		tc.Data["json"] = oResult
		tc.ServeJSON()
		return
	}
	err = models.SetTenantPermissions(tp.Tenantid, tp.Permissions)
	if err != nil {
		oResult.Result = false
		oResult.Message = err.Error()
		tc.Data["json"] = oResult
		tc.ServeJSON()
		return
	}
	oResult.Result = true
	tc.Data["json"] = oResult
	tc.ServeJSON()
}

//获取平台的基础权限
// @Title GetPlatformPermissionList
// @Description 获取平台的基础权限
// @Param   Authorization     header    string  true        "Token信息"
// @Success 200 {object} controllers.OperResult
// @Failure 400 Invalid email supplied
// @Failure 404 User not found
// @router /GetPlatformPermissionList [get]
func (tc *AuthorityManageController) GetPlatformPermissionList() {
	token := tc.Ctx.Request.Header.Get("Authorization")
	result := OperResult{}
	ok, _, err := tools.CheckAuthority(token, "Platform.Permission")
	if !ok {
		result.Result = false
		result.Message = err.Error()
		tc.Data["json"] = result
		tc.ServeJSON()
		return
	}
	permission, _ := models.GetPermissionByTenant(0) //获取平台的基础权限清单
	tc.Data["json"] = permission
	tc.ServeJSON()
}

//新增基础权限
// @Title CreatePermission
// @Description 新增基础权限
// @Param   Authorization     header    string  true        "Token信息"
// @Param   Id     body     int64  true        "权限ID"
// @Param   Name     body     string  true        "权限Name"
// @Param   DisplayName     body     string  true        "权限中文名"
// @Success 200 {object} controllers.OperResult
// @Failure 400 Invalid email supplied
// @Failure 404 User not found
// @router /CreatePermission [post]
func (tc *AuthorityManageController) CreatePermission() {
	token := tc.Ctx.Request.Header.Get("Authorization")
	oResult := OperResult{}
	ok, _, err := tools.CheckAuthority(token, "Platform.Permission")
	if !ok {
		oResult.Result = false
		oResult.Message = err.Error()
		tc.Data["json"] = oResult
		tc.ServeJSON()
		return
	}

	var id int64
	p := &models.Permission{}
	json.Unmarshal(tc.Ctx.Input.RequestBody, p)
	if p.Name == "" {
		oResult.Result = false
		oResult.Message = "权限名不能为空"
		tc.Data["json"] = oResult
		tc.ServeJSON()
		return
	}
	//判断是否有名称相同的权限
	tp, err := models.GetPermissionByName(p.Name)
	if err != nil {
		oResult.Result = false
		oResult.Message = err.Error()
		tc.Data["json"] = oResult
		tc.ServeJSON()
		return
	} else {
		if tp != nil {
			oResult.Result = false
			oResult.Message = "存在名称相同的权限，不可新增"
			tc.Data["json"] = oResult
			tc.ServeJSON()
			return
		}
	}
	p.CreationTime = time.Now()
	id, err = models.AddPermission(p)
	if err != nil {
		oResult.Result = false
		oResult.Message = err.Error()
		tc.Data["json"] = oResult
		tc.ServeJSON()
		return
	}
	oResult.Result = true
	oResult.Message = strconv.FormatInt(id, 10)
	tc.Data["json"] = oResult
	tc.ServeJSON()
}

//获取角色清单
// @Title RoleList
// @Description 获取角色清单
// @Param   Authorization     header    string  true        "Token信息"
// @Success 200 {object} controllers.TenantInfo
// @Failure 400 Invalid email supplied
// @Failure 404 User not found
// @router /RoleList [get]
func (tc *AuthorityManageController) RoleList() {
	token := tc.Ctx.Request.Header.Get("Authorization")
	result := TenantInfo{}
	ok, claims, err := tools.CheckAuthority(token, "Tenant.Role")
	if !ok {
		result.Result = false
		result.Message = err.Error()
		tc.Data["json"] = result
		tc.ServeJSON()
		return
	}

	tmp := strconv.FormatFloat(claims["jti"].(float64), 'f', -1, 64)
	userid, _ := strconv.ParseInt(tmp, 10, 64)

	user, _ := models.GetUserById(userid)
	tenantid := user.TenantId
	roles, _ := models.GetRolesByTenant(tenantid)
	rrs := &[]RoleInfo{}
	for _, r := range roles {
		rr := RoleInfo{}
		rr.RoleId = r.Id
		rr.Name = r.Name
		rr.DisplayName = r.DisplayName
		*rrs = append(*rrs, rr)
	}
	result.Result = true
	result.Message = "OK"
	result.Roles = *rrs
	tc.Data["json"] = result
	tc.ServeJSON()
}

//新增角色，同时为角色配权限
// @Title CreateRole
// @Description 新增角色
// @Param   Authorization     header    string  true        "Token信息"
// @Param   TenantId     body     int  true        "租户ID"
// @Param   Name     body     string  true        "角色Name"
// @Param   DisplayName     body     string  true        "角色中文名"
// @Success 200 {object} controllers.OperResult
// @Failure 400 Invalid email supplied
// @Failure 404 User not found
// @router /CreateRole [post]
func (tc *AuthorityManageController) CreateRole() {
	token := tc.Ctx.Request.Header.Get("Authorization")
	oResult := OperResult{}
	ok, _, err := tools.CheckAuthority(token, "Tenant.Role")
	if !ok {
		oResult.Result = false
		oResult.Message = err.Error()
		tc.Data["json"] = oResult
		tc.ServeJSON()
		return
	}

	var id int64
	rinput := &CreateRoleInput{}
	r := &models.Role{}
	plist := &[]models.Permission{}
	json.Unmarshal(tc.Ctx.Input.RequestBody, rinput)
	if rinput.Name == "" {
		oResult.Result = false
		oResult.Message = "角色名不能为空"
		tc.Data["json"] = oResult
		tc.ServeJSON()
		return
	}
	tr, err := models.GetRoleByName(rinput.Name, rinput.TenantId)
	if err != nil {
		oResult.Result = false
		oResult.Message = err.Error()
		tc.Data["json"] = oResult
		tc.ServeJSON()
		return
	} else {
		if tr != nil {
			oResult.Result = false
			oResult.Message = "存在名称相同的角色，不可新增"
			tc.Data["json"] = oResult
			tc.ServeJSON()
			return
		}
	}
	r.CreationTime = time.Now()
	r.Name = rinput.Name
	r.DisplayName = rinput.DisplayName
	r.NormalizedName = strings.ToUpper(rinput.Name)
	r.TenantId = rinput.TenantId
	if len(rinput.Permissions) > 0 {
		for _, p := range rinput.Permissions {
			*plist = append(*plist, p)
		}
	}
	id, err = models.CreateRole(r, plist)
	if err != nil {
		oResult.Result = false
		oResult.Message = err.Error()
		tc.Data["json"] = oResult
		tc.ServeJSON()
		return
	}
	oResult.Result = true
	oResult.Message = strconv.FormatInt(id, 10)
	tc.Data["json"] = oResult
	tc.ServeJSON()
}

//更新角色
// @Title UpdateRole
// @Description 更新角色
// @Param   Authorization     header    string  true        "Token信息"
// @Param   Id     body     int  true        "角色Id"
// @Param   Name     body     string  true        "角色Name"
// @Param   DisplayName     body     string  true        "角色中文名"
// @Success 200 {object} controllers.OperResult
// @Failure 400 Invalid email supplied
// @Failure 404 User not found
// @router /UpdateRole [post]
func (tc *AuthorityManageController) UpdateRole() {
	token := tc.Ctx.Request.Header.Get("Authorization")
	oResult := OperResult{}
	ok, _, err := tools.CheckAuthority(token, "Tenant.Role")
	if !ok {
		oResult.Result = false
		oResult.Message = err.Error()
		tc.Data["json"] = oResult
		tc.ServeJSON()
		return
	}

	rinput := &CreateRoleInput{}
	plist := &[]models.Permission{}
	json.Unmarshal(tc.Ctx.Input.RequestBody, rinput)
	if rinput.Id == 0 {
		oResult.Result = false
		oResult.Message = "角色不存在，不能更新"
		tc.Data["json"] = oResult
		tc.ServeJSON()
		return
	}
	if rinput.Name == "" {
		oResult.Result = false
		oResult.Message = "角色名不能为空"
		tc.Data["json"] = oResult
		tc.ServeJSON()
		return
	}
	tr, err := models.GetRoleByName(rinput.Name, rinput.TenantId)
	if err != nil {
		oResult.Result = false
		oResult.Message = err.Error()
		tc.Data["json"] = oResult
		tc.ServeJSON()
		return
	} else {
		if tr != nil {
			if tr.Id != rinput.Id {
				oResult.Result = false
				oResult.Message = "存在相同的角色名，不可修改"
				tc.Data["json"] = oResult
				tc.ServeJSON()
				return
			}
		}
	}
	tr, err = models.GetRoleById(rinput.Id)
	if err != nil {
		oResult.Result = false
		oResult.Message = err.Error()
		tc.Data["json"] = oResult
		tc.ServeJSON()
		return
	}
	tr.CreationTime = time.Now()
	tr.Name = rinput.Name
	tr.DisplayName = rinput.DisplayName
	tr.NormalizedName = strings.ToUpper(rinput.Name)

	if len(rinput.Permissions) > 0 {
		for _, p := range rinput.Permissions {
			*plist = append(*plist, p)
		}
	}
	err = models.UpdateRole(tr, plist)
	if err != nil {
		oResult.Result = false
		oResult.Message = err.Error()
		tc.Data["json"] = oResult
		tc.ServeJSON()
		return
	}
	oResult.Result = true
	tc.Data["json"] = oResult
	tc.ServeJSON()
}

//删除角色--软删除
// @Title DeleteRole
// @Description 删除角色
// @Param   Authorization     header    string  true        "Token信息"
// @Param   Id     body     int  true        "角色Id"
// @Success 200 {object} controllers.OperResult
// @Failure 400 Invalid email supplied
// @Failure 404 User not found
// @router /DeleteRole [post]
func (tc *AuthorityManageController) DeleteRole() {
	token := tc.Ctx.Request.Header.Get("Authorization")
	oResult := OperResult{}
	ok, _, err := tools.CheckAuthority(token, "Tenant.Role")
	if !ok {
		oResult.Result = false
		oResult.Message = err.Error()
		tc.Data["json"] = oResult
		tc.ServeJSON()
		return
	}

	r := &models.Role{}
	json.Unmarshal(tc.Ctx.Input.RequestBody, r)
	sid := r.Id
	if sid == 0 {
		oResult.Message = "角色名不能为空,不可删除"
		oResult.Result = false
		tc.Data["json"] = oResult
		tc.ServeJSON()
		return
	}
	v, err := models.GetRoleById(sid)
	if err != nil {
		oResult.Message = err.Error()
		oResult.Result = false
		tc.Data["json"] = oResult
		tc.ServeJSON()
		return
	}
	v.IsDeleted = true
	err = models.UpdateRoleById(v)

	if err != nil {
		oResult.Message = err.Error()
		oResult.Result = false
		tc.Data["json"] = oResult
		tc.ServeJSON()
		return
	}

	oResult.Result = true
	tc.Data["json"] = oResult
	tc.ServeJSON()
}

//给用户配角色
// @Title SetUserRoles
// @Description 给用户配角色
// @Param   Authorization     header    string  true        "Token信息"
// @Param   UserId     body     int64  true        "用户Id"
// @Param   Roles     body     []int  true        "角色List"
// @Success 200 {object} controllers.OperResult
// @Failure 400 Invalid email supplied
// @Failure 404 User not found
// @router /SetUserRoles [post]
func (tc *AuthorityManageController) SetUserRoles() {
	token := tc.Ctx.Request.Header.Get("Authorization")
	oResult := OperResult{}
	ok, _, err := tools.CheckAuthority(token, "Tenant.User")
	if !ok {
		oResult.Result = false
		oResult.Message = err.Error()
		tc.Data["json"] = oResult
		tc.ServeJSON()
		return
	}

	ur := &UserRolesInput{}
	json.Unmarshal(tc.Ctx.Input.RequestBody, ur)
	if ur.UserId == 0 {
		oResult.Result = false
		oResult.Message = "用户不能为空"
		tc.Data["json"] = oResult
		tc.ServeJSON()
		return
	}
	v, err := models.GetUserById(ur.UserId)
	if err != nil {
		oResult.Result = false
		oResult.Message = "用户不存在"
		tc.Data["json"] = oResult
		tc.ServeJSON()
		return
	}
	err = models.SetUserRoles(v.Id, v.TenantId, ur.Roles)
	if err != nil {
		oResult.Result = false
		oResult.Message = err.Error()
		tc.Data["json"] = oResult
		tc.ServeJSON()
		return
	}
	oResult.Result = true
	tc.Data["json"] = oResult
	tc.ServeJSON()
}
