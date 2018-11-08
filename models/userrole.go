package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Userrole struct {
	Id            int64     `orm:"column(Id);pk"`
	CreationTime  time.Time `orm:"column(CreationTime);type(datetime)"`
	CreatorUserId int64     `orm:"column(CreatorUserId);null"`
	RoleId        int       `orm:"column(RoleId)"`
	TenantId      int       `orm:"column(TenantId);null"`
	UserId        int64     `orm:"column(UserId)"`
	SysId         int       `orm:"column(SysId);null"`
}

func (t *Userrole) TableName() string {
	return "userrole"
}

func init() {
	orm.RegisterModel(new(Userrole))
}

// AddUserrole insert a new Userrole into database and returns
// last inserted Id on success.
func AddUserrole(m *Userrole) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetUserroleById retrieves Userrole by Id. Returns error if
// Id doesn't exist
func GetUserroleById(id int64) (v *Userrole, err error) {
	o := orm.NewOrm()
	v = &Userrole{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllUserrole retrieves all Userrole matches certain condition. Returns empty list if
// no records exist
func GetAllUserrole(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Userrole))
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

	var l []Userrole
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

// UpdateUserrole updates Userrole by Id and returns error if
// the record to be updated doesn't exist
func UpdateUserroleById(m *Userrole) (err error) {
	o := orm.NewOrm()
	v := Userrole{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteUserrole deletes Userrole by Id and returns error if
// the record to be deleted doesn't exist
func DeleteUserrole(id int64) (err error) {
	o := orm.NewOrm()
	v := Userrole{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Userrole{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func GetUserRolesByUser(userid int64) (userroles []Userrole, num int64) {
	o := orm.NewOrm()
	num, _ = o.QueryTable("userrole").Filter("userid", userid).All(&userroles)
	return userroles, num
}

//设置用户角色
func SetUserRoles(userid int64, tenantid int, roleids []int) (err error) {
	o := orm.NewOrm()
	err = o.Begin()
	//删除用户对应的角色和对应的权限
	oldUserroles, oldnum := GetUserRolesByUser(userid)
	if oldnum > 0 {
		for _, or := range oldUserroles {
			_, err = o.Delete(&Userrole{Id: or.Id})
			if err != nil {
				err = o.Rollback()
				return err
			}
		}
		oldUserPermissions := &[]Permission{}
		opNum, _ := o.QueryTable("permission").Filter("TenantId", tenantid).Filter("roleid", 0).Filter("userid", userid).All(oldUserPermissions)
		if opNum > 0 {
			for _, op := range *oldUserPermissions {
				_, err = o.Delete(&Permission{Id: op.Id})
				if err != nil {
					err = o.Rollback()
					return err
				}
			}
		}
	}
	//新配用户的角色和对应的权限
	troles := &[]Role{}
	rolePermissions := &[]Permission{}
	newUserPermissions := &[]Permission{}
	userroles := &[]Userrole{}
	tnum, _ := o.QueryTable("role").Filter("TenantId", tenantid).Filter("id__in", roleids).Filter("userid", 0).All(troles)
	rpnum, _ := o.QueryTable("permission").Filter("TenantId", tenantid).Filter("roleid__in", roleids).Filter("userid", 0).All(rolePermissions)
	if tnum > 0 {
		for _, r := range *troles {
			ur := Userrole{}
			ur.CreationTime = time.Now()
			ur.RoleId = r.Id
			ur.UserId = userid
			ur.TenantId = tenantid
			*userroles = append(*userroles, ur)
		}
		_, err = o.InsertMulti(len(*userroles), userroles)
		if err != nil {
			err = o.Rollback()
			return err
		}
		if rpnum > 0 {
			for _, rp := range *rolePermissions {
				up := Permission{}
				up.CreationTime = time.Now()
				up.Name = rp.Name
				up.TenantId = tenantid
				up.RoleId = 0
				up.UserId = userid
				up.DisplayName = rp.DisplayName
				*newUserPermissions = append(*newUserPermissions, up)
			}
			_, err = o.InsertMulti(len(*newUserPermissions), newUserPermissions)
			if err != nil {
				err = o.Rollback()
				return err
			}
		}

	}
	err = o.Commit()
	return err
}
