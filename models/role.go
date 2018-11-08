package models

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Role struct {
	Id                   int       `orm:"column(Id);pk"`
	ConcurrencyStamp     string    `orm:"column(ConcurrencyStamp);size(128);null"`
	CreationTime         time.Time `orm:"column(CreationTime);type(datetime)"`
	CreatorUserId        int64     `orm:"column(CreatorUserId);null"`
	DeleterUserId        int64     `orm:"column(DeleterUserId);null"`
	DeletionTime         time.Time `orm:"column(DeletionTime);type(datetime);null"`
	DisplayName          string    `orm:"column(DisplayName);size(45)"`
	IsDefault            bool      `orm:"column(IsDefault);size(1);null"`
	IsDeleted            bool      `orm:"column(IsDeleted);size(1)"`
	IsStatic             bool      `orm:"column(IsStatic);size(1);null"`
	LastModificationTime time.Time `orm:"column(LastModificationTime);type(datetime);null"`
	LastModifierUserId   int64     `orm:"column(LastModifierUserId);null"`
	Name                 string    `orm:"column(Name);size(32)"`
	NormalizedName       string    `orm:"column(NormalizedName);size(32)"`
	TenantId             int       `orm:"column(TenantId);null"`
	Description          string    `orm:"column(Description);size(1024);null"`
}

func (t *Role) TableName() string {
	return "role"
}

func init() {
	orm.RegisterModel(new(Role))
}

// AddRole insert a new Role into database and returns
// last inserted Id on success.
func AddRole(m *Role) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetRoleById retrieves Role by Id. Returns error if
// Id doesn't exist
func GetRoleById(id int) (v *Role, err error) {
	o := orm.NewOrm()
	v = &Role{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllRole retrieves all Role matches certain condition. Returns empty list if
// no records exist
func GetAllRole(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Role))
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

	var l []Role
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

// UpdateRole updates Role by Id and returns error if
// the record to be updated doesn't exist
func UpdateRoleById(m *Role) (err error) {
	o := orm.NewOrm()
	v := Role{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteRole deletes Role by Id and returns error if
// the record to be deleted doesn't exist
func DeleteRole(id int) (err error) {
	o := orm.NewOrm()
	v := Role{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Role{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func GetRoleByName(name string, tenantid int) (r *Role, err error) {
	o := orm.NewOrm()
	r = &Role{}
	if name != "" {
		err := o.QueryTable("role").Filter("TenantId", tenantid).Filter("name", name).One(r)
		if err == orm.ErrMultiRows {
			// 多条的时候报错
			return nil, err
		}
		if err == orm.ErrNoRows {
			// 没有找到记录
			return nil, nil
		}
		return r, nil
	}

	return nil, err
}

func CreateRole(ro *Role, plist *[]Permission) (id int64, err error) {
	o := orm.NewOrm()
	var rid int64
	err = o.Begin()
	rid, err = o.Insert(ro)
	if err != nil {
		err = o.Rollback()
		return rid, err
	}
	if len(*plist) > 0 {
		for i, p := range *plist {
			p.CreationTime = time.Now()
			p.RoleId, _ = strconv.Atoi(strconv.FormatInt(rid, 10))
			p.TenantId = ro.TenantId
			(*plist)[i] = p
		}
		_, err = o.InsertMulti(len(*plist), *plist)
		if err != nil {
			err = o.Rollback()
			return rid, err
		}
	}
	err = o.Commit()
	return rid, err
}

func UpdateRole(ro *Role, plist *[]Permission) (err error) {
	o := orm.NewOrm()
	err = o.Begin()
	_, err = o.Update(ro)
	if err != nil {
		err = o.Rollback()
		return err
	}
	oldPermissions, num := GetPermissionByRole(*ro)
	if num != 0 {
		for _, op := range oldPermissions {
			_, err = o.Delete(&Permission{Id: op.Id})
			if err != nil {
				err = o.Rollback()
				return err
			}
		}
	}

	if len(*plist) > 0 {
		for i, p := range *plist {
			p.CreationTime = time.Now()
			p.RoleId = ro.Id
			p.TenantId = ro.TenantId
			(*plist)[i] = p
		}
		_, err = o.InsertMulti(len(*plist), *plist)
		if err != nil {
			err = o.Rollback()
			return err
		}
	}
	err = o.Commit()
	return err
}

func GetRolesByTenant(tenantid int) (roles []Role, num int64) {
	o := orm.NewOrm()
	num, _ = o.QueryTable("role").Filter("TenantId", tenantid).All(&roles)
	return roles, num
}
