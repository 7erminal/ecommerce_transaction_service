package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Orders struct {
	OrderId       int64      `orm:"auto"`
	Customer      *Customers `orm:"rel(fk);null"`
	OrderNumber   string
	Quantity      int
	Cost          float32
	OrderDesc     string `orm:"column(order_desc)"`
	OrderLocation string `orm:"column(order_location)"`
	// Currency     *Currencies `orm:"rel(fk)"`
	Currency     int64
	OrderDate    time.Time `orm:"type(datetime)"`
	OrderEndDate time.Time `orm:"type(datetime)"`
	ReturnedDate time.Time `orm:"type(datetime)"`
	DateCreated  time.Time `orm:"type(datetime)"`
	DateModified time.Time `orm:"type(datetime)"`
	CreatedBy    *Users    `orm:"column(created_by);rel(fk);"`
	ModifiedBy   int64
	OrderDetails []*Order_items `orm:"reverse(many);null;"`
}

func init() {
	orm.RegisterModel(new(Orders))
}

// AddOrders insert a new Orders into database and returns
// last inserted Id on success.
func AddOrders(m *Orders) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetOrdersById retrieves Orders by Id. Returns error if
// Id doesn't exist
func GetOrdersById(id int64) (v *Orders, err error) {
	o := orm.NewOrm()
	v = &Orders{OrderId: id}
	if err = o.QueryTable(new(Orders)).Filter("OrderId", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetOrdersById retrieves Orders by Id. Returns error if
// Id doesn't exist
func GetOrdersByUser(id int64) (v *[]Orders, err error) {
	o := orm.NewOrm()
	v = &[]Orders{}
	if _, err = o.QueryTable(new(Orders)).Filter("CreatedBy__UserId", id).RelatedSel().All(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetOrderCount retrieves Items by Id. Returns error if
// Id doesn't exist
func GetOrderCount(query map[string]string) (c int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Orders))
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
	}

	if c, err = qs.Count(); err == nil {
		return c, nil
	}
	return 0, err
}

// GetOrderCount retrieves Items by Id. Returns error if
// Id doesn't exist
func GetOrderCountByType() (c int64, err error) {
	o := orm.NewOrm()
	if c, err = o.QueryTable(new(Orders)).Count(); err == nil {
		return c, nil
	}
	return 0, err
}

// GetAllOrders retrieves all Orders matches certain condition. Returns empty list if
// no records exist
func GetAllOrders(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Orders))
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

	var l []Orders
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

// UpdateOrders updates Orders by Id and returns error if
// the record to be updated doesn't exist
func UpdateOrdersById(m *Orders) (err error) {
	o := orm.NewOrm()
	v := Orders{OrderId: m.OrderId}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteOrders deletes Orders by Id and returns error if
// the record to be deleted doesn't exist
func DeleteOrders(id int64) (err error) {
	o := orm.NewOrm()
	v := Orders{OrderId: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Orders{OrderId: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
