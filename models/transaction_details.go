package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Transaction_details struct {
	TransactionDetailId    int64         `orm:"auto"`
	TransactionId          *Transactions `orm:"rel(fk); column(transaction_id)"`
	Amount                 float32
	Comment                string `orm:"size(255); omitempty; null"`
	SenderAccountNumber    string `orm:"size(255)"`
	RecipientAccountNumber string `orm:"size(255)"`
	StatusCode             string `orm:"size(20)"`
	StatusMessage          string `orm:"size(255)"`
	SenderId               int
	TransactionType        string    `orm:"size(255)"`
	DateCreated            time.Time `orm:"type(datetime)"`
	DateModified           time.Time `orm:"type(datetime)"`
	CreatedBy              int
	ModifiedBy             int
	Active                 int
}

func init() {
	orm.RegisterModel(new(Transaction_details))
}

// AddTransaction_details insert a new Transaction_details into database and returns
// last inserted Id on success.
func AddTransaction_details(m *Transaction_details) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetTransaction_detailsById retrieves Transaction_details by Id. Returns error if
// Id doesn't exist
func GetTransaction_detailsById(id int64) (v *Transaction_details, err error) {
	o := orm.NewOrm()
	v = &Transaction_details{TransactionDetailId: id}
	if err = o.QueryTable(new(Transaction_details)).Filter("TransactionId", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetTransaction_detailsByTransactionId retrieves Transaction_details by Id. Returns error if
// Id doesn't exist
func GetTransaction_detailsByTransaction(id *Transactions) (v *Transaction_details, err error) {
	o := orm.NewOrm()
	v = &Transaction_details{TransactionId: id}
	if err = o.QueryTable(new(Transaction_details)).Filter("TransactionId", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllTransaction_details retrieves all Transaction_details matches certain condition. Returns empty list if
// no records exist
func GetAllTransaction_details(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Transaction_details))
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

	var l []Transaction_details
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

// UpdateTransaction_details updates Transaction_details by Id and returns error if
// the record to be updated doesn't exist
func UpdateTransaction_detailsById(m *Transaction_details) (err error) {
	o := orm.NewOrm()
	v := Transaction_details{TransactionDetailId: m.TransactionDetailId}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteTransaction_details deletes Transaction_details by Id and returns error if
// the record to be deleted doesn't exist
func DeleteTransaction_details(id int64) (err error) {
	o := orm.NewOrm()
	v := Transaction_details{TransactionDetailId: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Transaction_details{TransactionDetailId: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
