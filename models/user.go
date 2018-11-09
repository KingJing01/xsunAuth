package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

type User struct {
	Id                     int64     `orm:"column(Id);pk"`
	AccessFailedCount      int       `orm:"column(AccessFailedCount)"`
	AuthenticationSource   string    `orm:"column(AuthenticationSource);size(64);null"`
	ConcurrencyStamp       string    `orm:"column(ConcurrencyStamp);size(128);null"`
	CreationTime           time.Time `orm:"column(CreationTime);type(datetime)"`
	CreatorUserId          int64     `orm:"column(CreatorUserId);null"`
	DeleterUserId          int64     `orm:"column(DeleterUserId);null"`
	DeletionTime           time.Time `orm:"column(DeletionTime);type(datetime);null"`
	EmailAddress           string    `orm:"column(EmailAddress);size(256)"`
	EmailConfirmationCode  string    `orm:"column(EmailConfirmationCode);size(328);null"`
	IsActive               bool      `orm:"column(IsActive);size(1)"`
	IsDeleted              bool      `orm:"column(IsDeleted);size(1)"`
	IsEmailConfirmed       bool      `orm:"column(IsEmailConfirmed);size(1)"`
	IsLockoutEnabled       bool      `orm:"column(IsLockoutEnabled);size(1)"`
	IsPhoneNumberConfirmed bool      `orm:"column(IsPhoneNumberConfirmed);size(1)"`
	IsTwoFactorEnabled     bool      `orm:"column(IsTwoFactorEnabled);size(1)"`
	LastLoginTime          time.Time `orm:"column(LastLoginTime);type(datetime);null"`
	LastModificationTime   time.Time `orm:"column(LastModificationTime);type(datetime);null"`
	LastModifierUserId     int64     `orm:"column(LastModifierUserId);null"`
	LockoutEndDateUtc      time.Time `orm:"column(LockoutEndDateUtc);type(datetime);null"`
	Name                   string    `orm:"column(Name);size(32)"`
	NormalizedEmailAddress string    `orm:"column(NormalizedEmailAddress);size(256)"`
	NormalizedUserName     string    `orm:"column(NormalizedUserName);size(32)"`
	PasswordResetCode      string    `orm:"column(PasswordResetCode);size(328);null"`
	PhoneNumber            string    `orm:"column(PhoneNumber);size(32);null"`
	SecurityStamp          string    `orm:"column(SecurityStamp);size(128);null"`
	Surname                string    `orm:"column(Surname);size(32)"`
	TenantId               int       `orm:"column(TenantId);null"`
	UserName               string    `orm:"column(UserName);size(32)"`
	SysID                  string    `orm:"column(SysId);size(32)"`
	SsoID                  string    `orm:"column(SsoId);size(32)"`
}

func (t *User) TableName() string {
	return "user"
}

func init() {
	orm.RegisterModel(new(User))
}

// AddUser insert a new User into database and returns
// last inserted Id on success.
func AddUser(m *User) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetUserById retrieves User by Id. Returns error if
// Id doesn't exist
func GetUserById(id int64) (v *User, err error) {
	o := orm.NewOrm()
	v = &User{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllUser retrieves all User matches certain condition. Returns empty list if
// no records exist
func GetAllUser(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(User))
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

	var l []User
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

// UpdateUser updates User by Id and returns error if
// the record to be updated doesn't exist
func UpdateUserById(m *User) (err error) {
	o := orm.NewOrm()
	v := User{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteUser deletes User by Id and returns error if
// the record to be deleted doesn't exist
func DeleteUser(id int64) (err error) {
	o := orm.NewOrm()
	v := User{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&User{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

// 根据用户名查询
func GetUserByName(username string) (result bool, err error) {
	o := orm.NewOrm()
	u := &User{}
	result = true
	err = o.QueryTable("user").Filter("UserName", username).One(u)
	if err == orm.ErrMultiRows {
		// 多条的时候报错
		result = false
		return result, err
	}
	if err == orm.ErrNoRows {
		// 没有找到记录
		result = false
		return result, err
	}
	return true, nil
}

//根据用户名、密码查询
func LoginCheck(username string, password string, sysId string) (result bool, user User, err error) {
	valid := validation.Validation{}
	resultMobile := valid.Mobile(username, "username")
	o := orm.NewOrm()
	u := &User{}
	result = true
	//登录名格式分析  手机号码直接 ssoUser验证 其他的使用user--->sso关联
	if resultMobile.Ok {
		err = o.Raw("select t2.* from ssouser t1 left join user t2 on t1.id = t2.SsoId and t2.SysId=? and t1.Phone=? and t1.Passwd=? ", sysId, username, password).QueryRow(&u)
	} else {
		err = o.Raw("select t2.* from ssouser t1 left join user t2 on t1.id = t2.SsoId and t2.SysId=? and t2.UserName=? and t1.Passwd=? ", sysId, username, password).QueryRow(&u)
	}
	user = *u
	// 判断是否有错误的返回
	if err != nil || int(user.Id) == 0 {
		result = false
		return result, user, err
	}
	return true, user, nil
}
