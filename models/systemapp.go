package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type Systemapp struct {
	Id            int       `orm:"column(Id);pk"`
	SysName       string    `orm:"column(SysName);size(45)"`
	CreationTime  time.Time `orm:"column(CreationTime);type(datetime);null"`
	CreatorUserId int64     `orm:"column(CreatorUserId);null"`
	IsDeleted     int       `orm:"column(IsDeleted);null"`
}

func (t *Systemapp) tableName() string {

	return "systemapp"
}

func init() {
	orm.RegisterModel(new(Systemapp))
}
