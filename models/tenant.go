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

const (
	AdminUserName   = "admin"
	DefaultPassWord = "123456"
)

type Tenant struct {
	Id                   int       `orm:"column(Id);pk"`
	ConnectionString     string    `orm:"column(ConnectionString);size(1024);null"`
	CreationTime         time.Time `orm:"column(CreationTime);type(datetime)"`
	CreatorUserId        int64     `orm:"column(CreatorUserId);null"`
	DeleterUserId        int64     `orm:"column(DeleterUserId);null"`
	DeletionTime         time.Time `orm:"column(DeletionTime);type(datetime);null"`
	EditionId            int64     `orm:"column(EditionId);null"`
	IsActive             bool      `orm:"column(IsActive);"`
	IsDeleted            bool      `orm:"column(IsDeleted);"`
	LastModificationTime time.Time `orm:"column(LastModificationTime);type(datetime);null"`
	LastModifierUserId   int64     `orm:"column(LastModifierUserId);null"`
	Name                 string    `orm:"column(Name);size(128)"`
	TenancyName          string    `orm:"column(TenancyName);size(64)"`
	SysId                int       `orm:"column(SysId);size(45)"`
}

func (t *Tenant) TableName() string {
	return "tenant"
}

func init() {
	orm.RegisterModel(new(Tenant))
}

// AddTenant insert a new Tenant into database and returns
// last inserted Id on success.
func AddTenant(m *Tenant) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetTenantById retrieves Tenant by Id. Returns error if
// Id doesn't exist
func GetTenantById(id int) (v *Tenant, err error) {
	o := orm.NewOrm()
	v = &Tenant{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllTenant retrieves all Tenant matches certain condition. Returns empty list if
// no records exist
func GetAllTenant(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Tenant))
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

	var l []Tenant
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

// UpdateTenant updates Tenant by Id and returns error if
// the record to be updated doesn't exist
func UpdateTenantById(m *Tenant) (err error) {
	o := orm.NewOrm()
	v := Tenant{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteTenant deletes Tenant by Id and returns error if
// the record to be deleted doesn't exist2
func DeleteTenant(id int) (err error) {
	o := orm.NewOrm()
	v := Tenant{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Tenant{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func GetTenant(id int, name string) (v *Tenant, err error) {
	o := orm.NewOrm()

	v = &Tenant{}
	if id != 0 {
		v.Id = id
		if err = o.Read(v); err == nil {
			return v, nil
		}
	}
	if name != "" {
		err := o.QueryTable("tenant").Filter("name", name).One(v)
		if err == orm.ErrMultiRows {
			// 多条的时候报错
			return nil, err
		}
		if err == orm.ErrNoRows {
			// 没有找到记录
			return nil, err
		}
		return v, nil
	}

	return nil, err
}

func GetTenantByName(name string) (v *Tenant, err error) {
	o := orm.NewOrm()

	v = &Tenant{}
	if name != "" {
		err := o.QueryTable("tenant").Filter("name", name).One(v)
		if err == orm.ErrMultiRows {
			// 多条的时候报错
			return nil, err
		}
		if err == orm.ErrNoRows {
			// 没有找到记录
			return nil, nil
		}
		return v, nil
	}

	return nil, err
}

func GetCount(id int, name string) (num int64, err error) {
	o := orm.NewOrm()
	if id != 0 {
		num, err = o.QueryTable("tenant").Filter("id", id).Count()
		return num, err
	}
	if name != "" {
		num, err = o.QueryTable("tenant").Filter("name", name).Count()
		return num, err
	}
	err = errors.New("查询参数不能为空")
	return 1, err
}

func CreateTenant(t *Tenant, emailAddress string) (result bool, tenant *Tenant, err error) {
	o := orm.NewOrm()
	var tid int64
	err = o.Begin()
	tid, err = o.Insert(t)
	if err != nil {
		err = o.Rollback()
		return false, nil, err
	}
	//新增租户时，默认新增租户管理员
	var uid int64
	u := &User{}
	u.TenantId, _ = strconv.Atoi(strconv.FormatInt(tid, 10))
	u.AccessFailedCount = 0
	u.CreationTime = time.Now()
	u.EmailAddress = emailAddress
	u.Name = AdminUserName
	u.UserName = AdminUserName
	u.NormalizedEmailAddress = strings.ToUpper(u.EmailAddress)
	u.NormalizedUserName = strings.ToUpper(AdminUserName)
	u.Password = DefaultPassWord
	u.Surname = AdminUserName
	uid, err = o.Insert(u)
	if err != nil {
		err = o.Rollback()
		return false, nil, err
	}
	//默认新增租户管理员角色
	r := &Role{}
	var rid int64
	r.TenantId, _ = strconv.Atoi(strconv.FormatInt(tid, 10))
	r.CreationTime = time.Now()
	r.DisplayName = "管理员"
	r.Name = "admin"
	r.NormalizedName = "ADMIN"
	rid, err = o.Insert(r)
	if err != nil {
		err = o.Rollback()
		return false, nil, err
	}
	//租户管理员分配管理员角色
	ur := &Userrole{}
	ur.CreationTime = time.Now()
	ur.TenantId, _ = strconv.Atoi(strconv.FormatInt(tid, 10))
	ur.RoleId, _ = strconv.Atoi(strconv.FormatInt(rid, 10))
	ur.UserId = uid
	_, err = o.Insert(ur)
	if err != nil {
		err = o.Rollback()
		return false, nil, err
	}
	err = o.Commit()
	t.Id, _ = strconv.Atoi(strconv.FormatInt(rid, 10))
	return true, t, err
}
