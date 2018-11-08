package models

import "github.com/astaxie/beego/orm"

// struct table ssouser
type Ssouser struct {
	Id     int    `orm:"column(Id);pk"`
	Phone  string `orm:"column(Phone);size(20)"`
	Passwd string `orm:"column(Passwd);size(45)"`
}

//retuan table name
func (t *Ssouser) tableName() string {
	return "ssouser"
}

func init() {
	orm.RegisterModel(new(Ssouser))
}
