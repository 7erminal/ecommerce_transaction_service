package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Branches struct {
	BranchId     int64  `orm:"auto"`
	Branch       string `orm:"size(80)"`
	Country      int64  `orm:"column(country_id)"`
	Location     string
	PhoneNumber  string
	Active       int       `orm:"omitempty"`
	DateCreated  time.Time `orm:"type(datetime);omitempty"`
	DateModified time.Time `orm:"type(datetime);omitempty"`
	CreatedBy    int       `orm:"omitempty"`
	ModifiedBy   int       `orm:"omitempty"`
}

func init() {
	orm.RegisterModel(new(Branches))
}

// AddBranches insert a new Branches into database and returns
// last inserted Id on success.
func AddBranches(m *Branches) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetBranchesById retrieves Branches by Id. Returns error if
// Id doesn't exist
func GetBranchesById(id int64) (v *Branches, err error) {
	o := orm.NewOrm()
	v = &Branches{BranchId: id}
	if err = o.QueryTable(new(Branches)).Filter("BranchId", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetBranchesById retrieves Branches by Id. Returns error if
// Id doesn't exist
func GetBranchesByName(name string) (v *Branches, err error) {
	o := orm.NewOrm()
	v = &Branches{Branch: name}
	if err = o.QueryTable(new(Branches)).Filter("Branch", name).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllBranches retrieves all Branches matches certain condition. Returns empty list if
// no records exist
func GetAllBranches(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Branches))
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

	var l []Branches
	qs = qs.OrderBy(sortFields...).RelatedSel()
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

// UpdateBranches updates Branches by Id and returns error if
// the record to be updated doesn't exist
func UpdateBranchesById(m *Branches) (err error) {
	o := orm.NewOrm()
	v := Branches{BranchId: m.BranchId}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteBranches deletes Branches by Id and returns error if
// the record to be deleted doesn't exist
func DeleteBranches(id int64) (err error) {
	o := orm.NewOrm()
	v := Branches{BranchId: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Branches{BranchId: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
