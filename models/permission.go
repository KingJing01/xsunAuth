package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Permission struct {
	Id            int64     `orm:"column(Id);pk"`
	CreationTime  time.Time `orm:"column(CreationTime);type(datetime)"`
	CreatorUserId int64     `orm:"column(CreatorUserId);null"`
	Discriminator string    `orm:"column(Discriminator);size(300)"`
	IsGranted     bool      `orm:"column(IsGranted);size(1)"`
	Name          string    `orm:"column(Name);size(128)"`
	TenantId      int       `orm:"column(TenantId);null"`
	RoleId        int       `orm:"column(RoleId);null"`
	UserId        int64     `orm:"column(UserId);null"`
	DisplayName   string    `orm:"column(DisplayName);size(50)"`
}

func (t *Permission) TableName() string {
	return "permission"
}

func init() {
	orm.RegisterModel(new(Permission))
}

// AddPermission insert a new Permission into database and returns
// last inserted Id on success.
func AddPermission(m *Permission) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetPermissionById retrieves Permission by Id. Returns error if
// Id doesn't exist
func GetPermissionById(id int64) (v *Permission, err error) {
	o := orm.NewOrm()
	v = &Permission{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllPermission retrieves all Permission matches certain condition. Returns empty list if
// no records exist
func GetAllPermission(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Permission))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Permission
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdatePermission updates Permission by Id and returns error if
// the record to be updated doesn't exist
func UpdatePermissionById(m *Permission) (err error) {
	o := orm.NewOrm()
	v := Permission{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeletePermission deletes Permission by Id and returns error if
// the record to be deleted doesn't exist
func DeletePermission(id int64) (err error) {
	o := orm.NewOrm()
	v := Permission{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Permission{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func GetPermissionByUser(userid int64) (permissions []Permission, num int64) {
	o := orm.NewOrm()
	num, _ = o.QueryTable("permission").Filter("UserId", userid).All(&permissions)
	return permissions, num
}

func GetPermissionByName(name string) (p *Permission, err error) {
	o := orm.NewOrm()

	p = &Permission{}
	if name != "" {
		err := o.QueryTable("permission").Filter("name", name).One(p)
		if err == orm.ErrMultiRows {
			// 多条的时候报错
			return nil, err
		}
		if err == orm.ErrNoRows {
			// 没有找到记录
			return nil, nil
		}
		return p, nil
	}

	return nil, err
}

//根据单个角色获取权限
func GetPermissionByRole(r Role) (permissions []Permission, num int64) {
	o := orm.NewOrm()
	num, _ = o.QueryTable("permission").Filter("tenantid", r.TenantId).Filter("roleid", r.Id).Filter("userid", 0).All(&permissions)
	return permissions, num
}

// //根据多个角色的权限
// func GetPermissionByRoles(rlist []Role) (permissions []Permission, num int64) {
// 	o := orm.NewOrm()
// 	var tenantid int
// 	var roleids []int
// 	if len(rlist) > 0 {
// 		tenantid = rlist[0].Id
// 		for _, r := range rlist {
// 			roleids = append(roleids, r.Id)
// 		}
// 		num, _ = o.QueryTable("permission").Filter("tenantid", tenantid).Filter("roleid__in", roleids).Filter("userid", 0).All(&permissions)
// 		return permissions, num
// 	} else {
// 		return nil, 0
// 	}
// }

//根据租户获取租户的权限
func GetPermissionByTenant(tenantid int) (permissions []Permission, num int64) {
	o := orm.NewOrm()
	num, _ = o.QueryTable("permission").Filter("tenantid", tenantid).Filter("roleid", 0).Filter("userid", 0).All(&permissions)
	return permissions, num
}

func SetTenantPermissions(tenantid int, permissionsids []int64) (err error) {
	o := orm.NewOrm()
	err = o.Begin()
	permissions := []Permission{}
	tpermissions := &[]Permission{}

	oldPermissions := &[]Permission{}
	oldnum, _ := o.QueryTable("permission").Filter("tenantid", tenantid).Filter("roleid", 0).Filter("userid", 0).All(oldPermissions)
	if oldnum != 0 {
		for _, op := range *oldPermissions {
			_, err = o.Delete(&Permission{Id: op.Id})
			if err != nil {
				err = o.Rollback()
				return err
			}
		}
	}

	num, _ := o.QueryTable("permission").Filter("id__in", permissionsids).All(&permissions)
	if num > 0 {
		for _, p := range permissions {
			tmpp := Permission{}
			tmpp.CreationTime = time.Now()
			tmpp.Name = p.Name
			tmpp.TenantId = tenantid
			tmpp.DisplayName = p.DisplayName
			*tpermissions = append(*tpermissions, tmpp)
		}
		_, err = o.InsertMulti(len(*tpermissions), tpermissions)
		if err != nil {
			err = o.Rollback()
			return err
		}
	}

	err = o.Commit()
	return err

}

func GetPermissionByUserAndPermission(userid int64, permissionName string) (p *Permission, err error) {
	o := orm.NewOrm()
	p = &Permission{}
	err = o.QueryTable("permission").Filter("userid", userid).Filter("name", permissionName).One(p)
	if err == orm.ErrMultiRows {
		// 多条的时候报错
		return nil, err
	}
	if err == orm.ErrNoRows {
		// 没有找到记录
		return nil, err
	}
	return p, nil
}
