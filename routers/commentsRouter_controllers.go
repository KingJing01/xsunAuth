package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["xsunAuth/controllers:AuthorityManageController"] = append(beego.GlobalControllerRouter["xsunAuth/controllers:AuthorityManageController"],
        beego.ControllerComments{
            Method: "AuthorityError",
            Router: `/AuthorityError`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["xsunAuth/controllers:AuthorityManageController"] = append(beego.GlobalControllerRouter["xsunAuth/controllers:AuthorityManageController"],
        beego.ControllerComments{
            Method: "CreatePermission",
            Router: `/CreatePermission`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["xsunAuth/controllers:AuthorityManageController"] = append(beego.GlobalControllerRouter["xsunAuth/controllers:AuthorityManageController"],
        beego.ControllerComments{
            Method: "CreateRole",
            Router: `/CreateRole`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["xsunAuth/controllers:AuthorityManageController"] = append(beego.GlobalControllerRouter["xsunAuth/controllers:AuthorityManageController"],
        beego.ControllerComments{
            Method: "CreateTenant",
            Router: `/CreateTenant`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["xsunAuth/controllers:AuthorityManageController"] = append(beego.GlobalControllerRouter["xsunAuth/controllers:AuthorityManageController"],
        beego.ControllerComments{
            Method: "DeleteRole",
            Router: `/DeleteRole`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["xsunAuth/controllers:AuthorityManageController"] = append(beego.GlobalControllerRouter["xsunAuth/controllers:AuthorityManageController"],
        beego.ControllerComments{
            Method: "DeleteTenant",
            Router: `/DeleteTenant`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["xsunAuth/controllers:AuthorityManageController"] = append(beego.GlobalControllerRouter["xsunAuth/controllers:AuthorityManageController"],
        beego.ControllerComments{
            Method: "GetPlatformPermissionList",
            Router: `/GetPlatformPermissionList`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["xsunAuth/controllers:AuthorityManageController"] = append(beego.GlobalControllerRouter["xsunAuth/controllers:AuthorityManageController"],
        beego.ControllerComments{
            Method: "GetTenant",
            Router: `/GetTenant`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["xsunAuth/controllers:AuthorityManageController"] = append(beego.GlobalControllerRouter["xsunAuth/controllers:AuthorityManageController"],
        beego.ControllerComments{
            Method: "GetUserInfo",
            Router: `/GetUserInfo`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["xsunAuth/controllers:AuthorityManageController"] = append(beego.GlobalControllerRouter["xsunAuth/controllers:AuthorityManageController"],
        beego.ControllerComments{
            Method: "Login",
            Router: `/Login`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["xsunAuth/controllers:AuthorityManageController"] = append(beego.GlobalControllerRouter["xsunAuth/controllers:AuthorityManageController"],
        beego.ControllerComments{
            Method: "Logout",
            Router: `/Logout`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["xsunAuth/controllers:AuthorityManageController"] = append(beego.GlobalControllerRouter["xsunAuth/controllers:AuthorityManageController"],
        beego.ControllerComments{
            Method: "RoleList",
            Router: `/RoleList`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["xsunAuth/controllers:AuthorityManageController"] = append(beego.GlobalControllerRouter["xsunAuth/controllers:AuthorityManageController"],
        beego.ControllerComments{
            Method: "SetTenantPermission",
            Router: `/SetTenantPermission`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["xsunAuth/controllers:AuthorityManageController"] = append(beego.GlobalControllerRouter["xsunAuth/controllers:AuthorityManageController"],
        beego.ControllerComments{
            Method: "SetUserRoles",
            Router: `/SetUserRoles`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["xsunAuth/controllers:AuthorityManageController"] = append(beego.GlobalControllerRouter["xsunAuth/controllers:AuthorityManageController"],
        beego.ControllerComments{
            Method: "UpdateRole",
            Router: `/UpdateRole`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["xsunAuth/controllers:AuthorityManageController"] = append(beego.GlobalControllerRouter["xsunAuth/controllers:AuthorityManageController"],
        beego.ControllerComments{
            Method: "UpdateTenant",
            Router: `/UpdateTenant`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["xsunAuth/controllers:AuthorityManageController"] = append(beego.GlobalControllerRouter["xsunAuth/controllers:AuthorityManageController"],
        beego.ControllerComments{
            Method: "GetTenantList",
            Router: `/all`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["xsunAuth/controllers:PermissionController"] = append(beego.GlobalControllerRouter["xsunAuth/controllers:PermissionController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["xsunAuth/controllers:PermissionController"] = append(beego.GlobalControllerRouter["xsunAuth/controllers:PermissionController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["xsunAuth/controllers:PermissionController"] = append(beego.GlobalControllerRouter["xsunAuth/controllers:PermissionController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: `/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["xsunAuth/controllers:PermissionController"] = append(beego.GlobalControllerRouter["xsunAuth/controllers:PermissionController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["xsunAuth/controllers:PermissionController"] = append(beego.GlobalControllerRouter["xsunAuth/controllers:PermissionController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:id`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["xsunAuth/controllers:RoleController"] = append(beego.GlobalControllerRouter["xsunAuth/controllers:RoleController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["xsunAuth/controllers:RoleController"] = append(beego.GlobalControllerRouter["xsunAuth/controllers:RoleController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["xsunAuth/controllers:RoleController"] = append(beego.GlobalControllerRouter["xsunAuth/controllers:RoleController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: `/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["xsunAuth/controllers:RoleController"] = append(beego.GlobalControllerRouter["xsunAuth/controllers:RoleController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["xsunAuth/controllers:RoleController"] = append(beego.GlobalControllerRouter["xsunAuth/controllers:RoleController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:id`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["xsunAuth/controllers:TenantController"] = append(beego.GlobalControllerRouter["xsunAuth/controllers:TenantController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["xsunAuth/controllers:TenantController"] = append(beego.GlobalControllerRouter["xsunAuth/controllers:TenantController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["xsunAuth/controllers:TenantController"] = append(beego.GlobalControllerRouter["xsunAuth/controllers:TenantController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: `/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["xsunAuth/controllers:TenantController"] = append(beego.GlobalControllerRouter["xsunAuth/controllers:TenantController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["xsunAuth/controllers:TenantController"] = append(beego.GlobalControllerRouter["xsunAuth/controllers:TenantController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:id`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["xsunAuth/controllers:UserController"] = append(beego.GlobalControllerRouter["xsunAuth/controllers:UserController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["xsunAuth/controllers:UserController"] = append(beego.GlobalControllerRouter["xsunAuth/controllers:UserController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["xsunAuth/controllers:UserController"] = append(beego.GlobalControllerRouter["xsunAuth/controllers:UserController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: `/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["xsunAuth/controllers:UserController"] = append(beego.GlobalControllerRouter["xsunAuth/controllers:UserController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["xsunAuth/controllers:UserController"] = append(beego.GlobalControllerRouter["xsunAuth/controllers:UserController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:id`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["xsunAuth/controllers:UserroleController"] = append(beego.GlobalControllerRouter["xsunAuth/controllers:UserroleController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["xsunAuth/controllers:UserroleController"] = append(beego.GlobalControllerRouter["xsunAuth/controllers:UserroleController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["xsunAuth/controllers:UserroleController"] = append(beego.GlobalControllerRouter["xsunAuth/controllers:UserroleController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: `/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["xsunAuth/controllers:UserroleController"] = append(beego.GlobalControllerRouter["xsunAuth/controllers:UserroleController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["xsunAuth/controllers:UserroleController"] = append(beego.GlobalControllerRouter["xsunAuth/controllers:UserroleController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:id`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
