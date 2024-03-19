package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Item_quantity struct {
	ItemQuantityId int64  `orm:"auto"`
	ItemId         *Items `orm:"rel(fk)"`
	Quantity       int
	Active         int
	DateCreated    time.Time `orm:"type(datetime)"`
	DateModified   time.Time `orm:"type(datetime)"`
	CreatedBy      int
	ModifiedBy     int
}

func init() {
	orm.RegisterModel(new(Item_quantity))
}

// AddItem_quantity insert a new Item_quantity into database and returns
// last inserted Id on success.
func AddItem_quantity(m *Item_quantity) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetItem_quantityById retrieves Item_quantity by Id. Returns error if
// Id doesn't exist
func GetItem_quantityById(id int64) (v *Item_quantity, err error) {
	o := orm.NewOrm()
	v = &Item_quantity{ItemQuantityId: id}
	if err = o.QueryTable(new(Item_quantity)).Filter("ItemQuantityId", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetItem_quantityById retrieves Item_quantity by Item Id. Returns error if
// Id doesn't exist
func GetItem_quantityByItemId(id int64) (v *Item_quantity, err error) {
	o := orm.NewOrm()
	if err = o.QueryTable(new(Item_quantity)).Filter("ItemId", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllItem_quantity retrieves all Item_quantity matches certain condition. Returns empty list if
// no records exist
func GetAllItem_quantity(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Item_quantity))
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

	var l []Item_quantity
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

// UpdateItem_quantity updates Item_quantity by Id and returns error if
// the record to be updated doesn't exist
func UpdateItem_quantityById(m *Item_quantity) (err error) {
	o := orm.NewOrm()
	v := Item_quantity{ItemQuantityId: m.ItemQuantityId}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteItem_quantity deletes Item_quantity by Id and returns error if
// the record to be deleted doesn't exist
func DeleteItem_quantity(id int64) (err error) {
	o := orm.NewOrm()
	v := Item_quantity{ItemQuantityId: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Item_quantity{ItemQuantityId: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
