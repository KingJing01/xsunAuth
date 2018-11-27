package main

import (
	"regexp"
	_ "xsunAuth/routers"
	"xsunAuth/tools"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {

	orm.RegisterDataBase("default", "mysql", "root:root@tcp(127.0.0.1:3308)/godatabase?charset=utf8&loc=Local")

	orm.Debug = true

	var FilterUser = func(ctx *context.Context) {
		//正则匹配访问路径
		b, _ := regexp.MatchString("/v1/authoritymanage/xsunLogin", ctx.Request.RequestURI)
		if b == false {
			ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")                           //允许访问源
			ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS")    //允许post访问
			ctx.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization") //header的类型
			ctx.ResponseWriter.Header().Set("Access-Control-Max-Age", "1728000")
			ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
			ctx.ResponseWriter.Header().Set("content-type", "application/json") //返回数据格式是json
			token := ctx.Request.Header.Get("Authorization")
			if token == "" && ctx.Request.RequestURI != "/v1/authoritymanage/AuthorityError" && ctx.Request.RequestURI != "/v1/authoritymanage/Login" {
				ctx.Redirect(302, "/v1/authoritymanage/AuthorityError")

			} else if token != "" && ctx.Request.RequestURI != "/v1/authoritymanage/AuthorityError" && ctx.Request.RequestURI != "/v1/authoritymanage/Login" {
				result, _, _ := tools.CheckLogin(token)
				if result == false {
					ctx.Redirect(302, "/v1/authoritymanage/AuthorityError")
				}
			}
		}
	}

	beego.InsertFilter("/v1/*", beego.BeforeRouter, FilterUser)

}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
