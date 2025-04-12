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

type Transactions struct {
	TransactionId       int64     `orm:"auto"`
	Order               *Orders   `orm:"rel(fk)"`
	Branch              *Branches `orm:"rel(fk)"`
	Amount              float32
	TransactingCurrency int64
	Status              *Status   `orm:"rel(fk);column(status_id)"`
	DateCreated         time.Time `orm:"type(datetime)"`
	DateModified        time.Time `orm:"type(datetime)"`
	CreatedBy           int
	ModifiedBy          int
	Active              int
	Services            *Services   `orm:"rel(fk);column(service_id)"`
	Payments            []*Payments `orm:"reverse(many);null;"`
}

func init() {
	orm.RegisterModel(new(Transactions))
}

// AddTransactions insert a new Transactions into database and returns
// last inserted Id on success.
func AddTransactions(m *Transactions) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetTransactionsById retrieves Transactions by Id. Returns error if
// Id doesn't exist
func GetTransactionsById(id int64) (v *Transactions, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Transactions))

	v = &Transactions{TransactionId: id}
	if err = qs.Filter("TransactionId", id).RelatedSel().One(v); err == nil {
		// q := Transactions{TransactionId: id}
		_, err := o.LoadRelated(v, "Payments")

		if err == nil {
			logs.Info("Loaded payments ", v)
			fmt.Printf("Payments loaded is: %+v\n", v.Payments)

			for _, payment := range v.Payments {
				// payment_ := Payments{PaymentId: payment.PaymentId}
				// err = o.Read(payment)
				// if err == nil {
				_, err := o.LoadRelated(payment, "PaymentMethod")
				if err == nil {
					// payment.PaymentMethod = &paymentMethod
					logs.Info("Payment method is ", payment.PaymentMethod)
					fmt.Printf("Payment method loaded is: %+v\n", payment.PaymentMethod)
				}
				// }
			}
		}

		_, err = o.LoadRelated(v.Order, "OrderDetails")
		if err == nil {
			logs.Info("No error. Continue to loop through orders ", v.Order)
			fmt.Printf("Order loop: %+v\n", v.Order)
			// orderitems := []Order_items{}
			for _, orderD := range v.Order.OrderDetails {
				logs.Info("Each order detail is ", orderD)
				// orderitem := Order_items{OrderItemId: orderD.OrderItemId}
				// err := o.Read(&orderitem)
				// if err == nil {
				_, err := o.LoadRelated(orderD, "Item")
				logs.Info("No error getting item details")
				// _, err := o.LoadRelated(&orderitem, "Item")
				if err == nil {
					// orderitems = append(orderitems, orderitem)
					// item := Items{ItemId: orderD.Item.ItemId}
					// err := o.Read(&item)
					// if err == nil {
					_, err := o.LoadRelated(orderD.Item, "Category")
					if err != nil {
						logs.Error("Failed loading category")
					}
					logs.Info("Item loaded is ", orderD.Item)
					fmt.Printf("Category loaded is: %+v\n", orderD.Item.Category)

					_, err = o.LoadRelated(orderD.Item, "ItemPrice")
					if err != nil {
						logs.Error("Failed loading item price")
					}
					logs.Info("Item price is ", orderD.Item.ItemPrice)

					_, err = o.LoadRelated(orderD.Item, "Branch")
					if err != nil {
						logs.Error("Failed loading branch")
					}
					logs.Info("Item price is ", orderD.Item.Branch)
					// }
				} else {
					logs.Error("Error loading related item ", err.Error())
				}
				// } else {
				// 	logs.Error("An error occurred when reading order item ", err.Error())
				// }
			}
		} else {
			logs.Error("Error loading related..", err.Error())
		}

		logs.Info("About to load related for customer involved")
		_, err = o.LoadRelated(v.Order, "Customer")
		if err == nil {
			logs.Info("No error. Continue to loop through orders ", v.Order)
			fmt.Printf("Order loop: %+v\n", v.Order.Customer)
			// orderitems := []Order_items{}

		} else {
			logs.Error("Error loading related..", err.Error())
		}
		// v.Order = &v.Order

		fmt.Printf("Order details are: %+v\n", v.Order.OrderDetails)

		// for _, orderD := range v.Order.OrderDetails {
		// 	logs.Info("Each order detail is ", orderD)
		// }

		return v, nil
	}
	return nil, err
}

// GetOrdersById retrieves Orders by Id. Returns error if
// Id doesn't exist
func GetTransactionsByUser(id int64) (v *[]Transactions, err error) {
	o := orm.NewOrm()
	v = &[]Transactions{}
	if _, err = o.QueryTable(new(Transactions)).Filter("Order__CreatedBy__UserId", id).RelatedSel().All(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetOrdersById retrieves Orders by Id. Returns error if
// Id doesn't exist
func GetTransactionsByUserWithDate(id int64, fromDate time.Time, toDate time.Time) (v *[]Transactions, err error) {
	o := orm.NewOrm()
	v = &[]Transactions{}
	if _, err = o.QueryTable(new(Transactions)).Filter("Order__CreatedBy__UserId", id).Filter("DateCreated__gte", fromDate).Filter("DateCreated__lte", toDate).RelatedSel().All(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetOrdersById retrieves Orders by Id. Returns error if
// Id doesn't exist
func GetTransactionsByUserWithLimit(id int64, limit int) (v *[]Transactions, err error) {
	o := orm.NewOrm()
	v = &[]Transactions{}
	if _, err = o.QueryTable(new(Transactions)).Filter("Order__CreatedBy__UserId", id).RelatedSel().Limit(limit).All(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllTransactions retrieves all Transactions matches certain condition. Returns empty list if
// no records exist
func GetAllTransactions(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64, search map[string]string) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Transactions))

	logs.Info("Getting all transactions")

	if len(search) > 0 {
		cond := orm.NewCondition()
		for k, v := range search {
			// rewrite dot-notation to Object__Attribute
			k = strings.Replace(k, ".", "__", -1)
			if strings.Contains(k, "isnull") {
				qs = qs.Filter(k, (v == "true" || v == "1"))
			} else {
				logs.Info("Adding or statement")
				cond = cond.Or(k+"__icontains", v)

				// qs = qs.Filter(k+"__icontains", v)

			}
		}
		logs.Info("Condition set ", qs)
		qs = qs.SetCond(cond)
	}
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		logs.Info("Modified query is ", k, " and ", v)
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

	var l []Transactions
	qs = qs.OrderBy(sortFields...).RelatedSel()
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		logs.Info("Data fetched ", qs)
		fmt.Printf("Type of v: %T\n", qs)
		fmt.Printf("Value of v: %+v\n", qs)
		if len(fields) == 0 {
			logs.Info("Data fetched for txns ", l)
			for _, v := range l {
				fmt.Printf("Value of this guy: %+v\n", v)
				// transaction := Transactions{TransactionId: v.TransactionId}

				logs.Info("Order ID is ", v.Order.OrderId)

				logs.Info("Error is nil. Proceeding with ")
				fmt.Printf("Value of order: %+v\n", v.Order)

				// order.Customer = order.Customer
				_, err = o.LoadRelated(v.Order, "OrderDetails")
				if err == nil {
					logs.Info("No error. Continue to loop through orders ", v.Order)
					fmt.Printf("Order loop: %+v\n", v.Order)
					// orderitems := []Order_items{}
					for _, orderD := range v.Order.OrderDetails {
						logs.Info("Each order detail is ", orderD)
						// orderitem := Order_items{OrderItemId: orderD.OrderItemId}
						// err := o.Read(&orderitem)
						// if err == nil {
						_, err := o.LoadRelated(orderD, "Item")
						logs.Info("No error getting item details")
						// _, err := o.LoadRelated(&orderitem, "Item")
						if err == nil {
							// orderitems = append(orderitems, orderitem)
							// item := Items{ItemId: orderD.Item.ItemId}
							// err := o.Read(&item)
							// if err == nil {
							_, err := o.LoadRelated(orderD.Item, "Category")
							if err != nil {
								logs.Error("Failed loading category")
							}
							logs.Info("Item loaded is ", orderD.Item)
							fmt.Printf("Category loaded is: %+v\n", orderD.Item.Category)

							_, err = o.LoadRelated(orderD.Item, "ItemPrice")
							if err != nil {
								logs.Error("Failed loading item price")
							}
							logs.Info("Item price is ", orderD.Item.ItemPrice)

							_, err = o.LoadRelated(orderD.Item, "Branch")
							if err != nil {
								logs.Error("Failed loading branch")
							}
							logs.Info("Item price is ", orderD.Item.Branch)
							// }
						} else {
							logs.Error("Error loading related item ", err.Error())
						}
						// } else {
						// 	logs.Error("An error occurred when reading order item ", err.Error())
						// }
					}
				} else {
					logs.Error("Error loading related..", err.Error())
				}

				_, err = o.LoadRelated(v.Order, "Customer")
				if err == nil {
					fmt.Printf("Order loop: %+v\n", v.Order.Customer)
				}

				_, err = o.LoadRelated(&v, "Payments")
				if err == nil {
					logs.Info("Loaded payments ", v)
					fmt.Printf("Payments loaded is: %+v\n", v.Payments)

					for _, payment := range v.Payments {
						// payment_ := Payments{PaymentId: payment.PaymentId}
						// err = o.Read(payment)
						// if err == nil {
						_, err := o.LoadRelated(payment, "PaymentMethod")
						if err == nil {
							// payment.PaymentMethod = &paymentMethod
							logs.Info("Payment method is ", payment.PaymentMethod)
							fmt.Printf("Payment method loaded is: %+v\n", payment.PaymentMethod)
						}
						// }
					}
				}
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			logs.Info("Data fetched for txns 2 ", l)
			for _, v := range l {
				fmt.Printf("Value of vv: %+v\n", v)
				m := make(map[string]interface{})
				transaction := Transactions{TransactionId: v.TransactionId}
				err := o.Read(&transaction)
				if err == nil {
					logs.Info("Error is not nil. Proceeding")
					_, err := o.LoadRelated(&transaction, "Order__OrderDetails")
					if err == nil {
						for _, order := range transaction.Order.OrderDetails {
							logs.Info("Each order detail is ", order)
						}
					}
				} else {
					logs.Error("Error reading transaction related. ", err.Error())
				}
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

// UpdateTransactions updates Transactions by Id and returns error if
// the record to be updated doesn't exist
func UpdateTransactionsById(m *Transactions) (err error) {
	o := orm.NewOrm()
	v := Transactions{TransactionId: m.TransactionId}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteTransactions deletes Transactions by Id and returns error if
// the record to be deleted doesn't exist
func DeleteTransactions(id int64) (err error) {
	o := orm.NewOrm()
	v := Transactions{TransactionId: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Transactions{TransactionId: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

// GetTransactionCount retrieves Items by Id. Returns error if
// Id doesn't exist
func GetTransactionCount(query map[string]string, search map[string]string) (c int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Transactions))
	if len(search) > 0 {
		cond := orm.NewCondition()
		for k, v := range search {
			// rewrite dot-notation to Object__Attribute
			k = strings.Replace(k, ".", "__", -1)
			if strings.Contains(k, "isnull") {
				qs = qs.Filter(k, (v == "true" || v == "1"))
			} else {
				logs.Info("Adding or statement")
				cond = cond.Or(k+"__icontains", v)

				// qs = qs.Filter(k+"__icontains", v)

			}
		}
		logs.Info("Condition set ", qs)
		qs = qs.SetCond(cond)
	}

	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
		logs.Info("Modified query is ", k, " and ", v)
	}

	if c, err = qs.RelatedSel().Count(); err == nil {
		logs.Info("Count of transactions is ", c)
		return c, nil
	}
	return 0, err
}
