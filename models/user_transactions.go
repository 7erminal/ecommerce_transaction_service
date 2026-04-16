package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type UserTransactions struct {
	TransactionId                string       `orm:"size(255);unique"`
	Service                      *Services    `orm:"rel(fk)"`
	Request                      *UserRequest `orm:"rel(fk)"`
	TransactionCustomerReference *Customers   `orm:"rel(fk);column(transaction_customer_reference)"`
	Amount                       float64
	TransactingCurrency          string `orm:"size(255)"`
	SourceChannel                string `orm:"size(255)"`
	Source                       string `orm:"size(255)"`
	Destination                  string `orm:"size(255)"`
	Package                      string `orm:"size(255)"`
	Charge                       float64
	Commission                   float64
	ExternalReferenceNumber      string `orm:"size(255)"`
	IsAsync                      bool
	CorpId                       string        `orm:"size(255)"`
	Status                       *Status_codes `orm:"rel(fk)"`
	ExtraDetails1                string        `orm:"size(255); column(extra_details_1)"`
	ExtraDetails2                string        `orm:"size(255); column(extra_details_2)"`
	ExtraDetails3                string        `orm:"size(255); column(extra_details_3)"`
	ClientResponseCode           string        `orm:"size(255); column(client_response_code)"`
	DateCreated                  time.Time     `orm:"type(datetime)"`
	DateModified                 time.Time     `orm:"type(datetime)"`
	CreatedBy                    *Users        `orm:"rel(fk);column(created_by)"`
	ModifiedBy                   *Users        `orm:"rel(fk);column(modified_by)"`
	Active                       int
}

func (t *UserTransactions) TableName() string {
	return "user_transactions"
}

func init() {
	orm.RegisterModel(new(UserTransactions))
}

// AddUserTransactions insert a new UserTransactions into database and returns
// last inserted Id on success.
func AddUserTransactions(m *UserTransactions) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetUserTransactionsById retrieves UserTransactions by Id. Returns error if
// Id doesn't exist
func GetUserTransactionsById(id string) (v *UserTransactions, err error) {
	o := orm.NewOrm()
	v = &UserTransactions{TransactionId: id}
	if err = o.QueryTable(new(UserTransactions)).Filter("TransactionId", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllUserTransactions retrieves all UserTransactions matches certain condition. Returns empty list if
// no records exist
func GetAllUserTransactions(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(UserTransactions))
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

	var l []UserTransactions
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

// UpdateUserTransactionsById updates UserTransactions by Id and returns error if
// the record to be updated doesn't exist
func UpdateUserTransactionsById(m *UserTransactions) (err error) {
	o := orm.NewOrm()
	v := UserTransactions{TransactionId: m.TransactionId}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteUserTransactions deletes UserTransactions by Id and returns error if
// the record to be deleted doesn't exist
func DeleteUserTransactions(id string) (err error) {
	o := orm.NewOrm()
	v := UserTransactions{TransactionId: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&UserTransactions{TransactionId: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
