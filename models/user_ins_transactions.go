package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

type UserInsTransactions struct {
	UserInsTransactionId   string `orm:"size(255);unique;column(ins_transaction_id)"`
	UserTransactionId      string `orm:"size(255);column(user_transaction_id)"`
	Amount                 float64
	Data                   string        `orm:"size(255);null"`
	SenderAccountNumber    string        `orm:"size(255)"`
	RecipientAccountNumber string        `orm:"size(255)"`
	Service                *Services     `orm:"rel(fk)"`
	Status                 *Status_codes `orm:"rel(fk);column(status)"`
	Request                string        `orm:"size(255);null"`
	Response               string        `orm:"size(255);null"`
	DateCreated            time.Time     `orm:"type(datetime)"`
	DateModified           time.Time     `orm:"type(datetime)"`
	CreatedBy              int
	ModifiedBy             int
	Active                 int
}

func init() {
	orm.RegisterModel(new(UserInsTransactions))
}

// AddUserInsTransactions insert a new UserInsTransactions into database and returns
// last inserted Id on success.
func AddUserInsTransactions(m *UserInsTransactions) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetUserInsTransactionsById retrieves UserInsTransactions by Id. Returns error if
// Id doesn't exist
func GetUserInsTransactionsById(id string) (v *UserInsTransactions, err error) {
	o := orm.NewOrm()
	v = &UserInsTransactions{UserInsTransactionId: id}
	if err = o.QueryTable(new(UserInsTransactions)).Filter("UserInsTransactionId", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllUserInsTransactions retrieves all UserInsTransactions matches certain condition. Returns empty list if
// no records exist
func GetAllUserInsTransactions(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(UserInsTransactions))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		if strings.HasSuffix(k, "__in") {
			// split by comma into slice
			items := strings.Split(v, ",")
			logs.Info("Items are: ", v)
			// trim spaces
			for i := range items {
				items[i] = strings.TrimSpace(items[i])
			}
			logs.Info("Items after trim are: ", items)
			k = strings.Replace(k, ".", "__", -1)
			qs = qs.Filter(k, items) // []string passed here ✅
		} else {
			k = strings.Replace(k, ".", "__", -1)
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

	var l []UserInsTransactions
	qs = qs.OrderBy(sortFields...).RelatedSel()
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		logs.Info("ORM Fetched records: ", len(l))
		logs.Info("ORM Records: ", l)
		if len(fields) == 0 {
			for _, v := range l {
				// Load related fields as needed
				o.LoadRelated(&v, "UserTransactionId")

				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				// Load related fields as needed
				o.LoadRelated(&v, "UserTransactionId")

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

// UpdateUserInsTransactions updates UserInsTransactions by Id and returns error if
// the record to be updated doesn't exist
func UpdateUserInsTransactionsById(m *UserInsTransactions) (err error) {
	o := orm.NewOrm()
	v := UserInsTransactions{UserInsTransactionId: m.UserInsTransactionId}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteUserInsTransactions deletes UserInsTransactions by Id and returns error if
// the record to be deleted doesn't exist
func DeleteUserInsTransactions(id string) (err error) {
	o := orm.NewOrm()
	v := UserInsTransactions{UserInsTransactionId: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&UserInsTransactions{UserInsTransactionId: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
