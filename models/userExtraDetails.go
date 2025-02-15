package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type UserExtraDetails struct {
	UserDetailsId int64     `orm:"auto"`
	User          int64     `orm:"column(user_id)"`
	Branch        *Branches `orm:"rel(fk);column(branch);omitempty;null"`
	Shop          int64
	Nickname      string    `orm:"size(100);omitempty;null"`
	DateCreated   time.Time `orm:"type(datetime)"`
	DateModified  time.Time `orm:"type(datetime)"`
	CreatedBy     int
	ModifiedBy    int
	Active        int
}

func init() {
	orm.RegisterModel(new(UserExtraDetails))
}

// AddCustomers insert a new Customers into database and returns
// last inserted Id on success.
func AddUserExtraDetails(m *UserExtraDetails) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetCustomersById retrieves Customers by Id. Returns error if
// Id doesn't exist
func GetUserExtraDetailsById(id int64) (v *UserExtraDetails, err error) {
	o := orm.NewOrm()
	v = &UserExtraDetails{UserDetailsId: id}
	if err = o.QueryTable(new(UserExtraDetails)).Filter("UserDetailsId", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetCustomersByUserId retrieves Customers by User Id. Returns error if
// Id doesn't exist
func GetUserExtraDetailsByUser(user int64) (v *UserExtraDetails, err error) {
	o := orm.NewOrm()
	v = &UserExtraDetails{User: user}
	if err = o.QueryTable(new(UserExtraDetails)).Filter("User", user).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetCustomersByUserId retrieves Customers by User Id. Returns error if
// Id doesn't exist
func GetUserExtraDetailsByBranch(branch *Branches) (v *UserExtraDetails, err error) {
	o := orm.NewOrm()
	v = &UserExtraDetails{Branch: branch}
	if err = o.QueryTable(new(UserExtraDetails)).Filter("Branch", branch).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllCustomers retrieves all Customers matches certain condition. Returns empty list if
// no records exist
func GetAllUserExtraDetails(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(UserExtraDetails))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
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

	var l []UserExtraDetails
	qs = qs.OrderBy(sortFields...).RelatedSel()

	// Without Limit
	if _, err = qs.All(&l, fields...); err == nil {
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

	// With limit
	// if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
	// 	if len(fields) == 0 {
	// 		for _, v := range l {
	// 			ml = append(ml, v)
	// 		}
	// 	} else {
	// 		// trim unused fields
	// 		for _, v := range l {
	// 			m := make(map[string]interface{})
	// 			val := reflect.ValueOf(v)
	// 			for _, fname := range fields {
	// 				m[fname] = val.FieldByName(fname).Interface()
	// 			}
	// 			ml = append(ml, m)
	// 		}
	// 	}
	// 	return ml, nil
	// }
	return nil, err
}

// UpdateCustomers updates Customers by Id and returns error if
// the record to be updated doesn't exist
func UpdateUserExtraDetailsById(m *UserExtraDetails) (err error) {
	o := orm.NewOrm()
	v := UserExtraDetails{UserDetailsId: m.UserDetailsId}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteCustomers deletes Customers by Id and returns error if
// the record to be deleted doesn't exist
func DeleteUserExtraDetails(id int64) (err error) {
	o := orm.NewOrm()
	v := UserExtraDetails{UserDetailsId: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&UserExtraDetails{UserDetailsId: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
