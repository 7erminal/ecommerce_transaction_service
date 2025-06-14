package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Payments struct {
	PaymentId       int64 `orm:"auto"`
	InitiatedBy     int64
	Transaction     *Transactions `orm:"rel(fk)"`
	Sender          *Customers    `orm:"rel(fk);column(sender)"`
	Reciever        *Users        `orm:"rel(fk);column(reciever)"`
	Amount          float64
	PaymentMethod   *Payment_methods `orm:"rel(fk);column(payment_method)"`
	PaymentProof    string           `orm:"null"`
	Status          int64
	PaymentAccount  int
	DateCreated     time.Time `orm:"type(datetime)"`
	DateModified    time.Time `orm:"type(datetime)"`
	CreatedBy       int64
	ModifiedBy      int64
	Active          int
	ReferenceNumber string
}

func init() {
	orm.RegisterModel(new(Payments))
}

// AddPayments insert a new Payments into database and returns
// last inserted Id on success.
func AddPayments(m *Payments) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetOrderCount retrieves Items by Id. Returns error if
// Id doesn't exist
func GetPaymentCount() (c int64, err error) {
	o := orm.NewOrm()
	if c, err = o.QueryTable(new(Payments)).Count(); err == nil {
		return c, nil
	}
	return 0, err
}

// GetPaymentsById retrieves Payments by Id. Returns error if
// Id doesn't exist
func GetPaymentsById(id int64) (v *Payments, err error) {
	o := orm.NewOrm()
	v = &Payments{PaymentId: id}
	if err = o.QueryTable(new(Payments)).Filter("PaymentId", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllPayments retrieves all Payments matches certain condition. Returns empty list if
// no records exist
func GetAllPayments(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Payments))
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

	var l []Payments
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

// UpdatePayments updates Payments by Id and returns error if
// the record to be updated doesn't exist
func UpdatePaymentsById(m *Payments) (err error) {
	o := orm.NewOrm()
	v := Payments{PaymentId: m.PaymentId}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeletePayments deletes Payments by Id and returns error if
// the record to be deleted doesn't exist
func DeletePayments(id int64) (err error) {
	o := orm.NewOrm()
	v := Payments{PaymentId: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Payments{PaymentId: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
