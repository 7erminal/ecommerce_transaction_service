package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Items struct {
	ItemId          int64        `orm:"auto" orm:"omitempty"`
	ItemName        string       `orm:"size(80)"`
	Description     string       `orm:"size(250)" orm:"omitempty"`
	Category        *Categories  `orm:"rel(fk)"`
	ItemPrice       *Item_prices `orm:"rel(fk)" orm:"omitempty"`
	AvailableSizes  string       `orm:"size(250)" orm:"omitempty"`
	AvailableColors string       `orm:"size(250)" orm:"omitempty"`
	ImagePath       string       `orm:"size(250)" orm:"omitempty"`
	Quantity        int          `orm:"omitempty"`
	Active          int          `orm:"omitempty"`
	DateCreated     time.Time    `orm:"type(datetime)" orm:"omitempty"`
	DateModified    time.Time    `orm:"type(datetime)" orm:"omitempty"`
	CreatedBy       int          `orm:"omitempty"`
	ModifiedBy      int          `orm:"omitempty"`
}

func init() {
	orm.RegisterModel(new(Items))
}

// AddItems insert a new Items into database and returns
// last inserted Id on success.
func AddItems(m *Items) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetItemsById retrieves Items by Id. Returns error if
// Id doesn't exist
func GetItemsById(id int64) (v *Items, err error) {
	o := orm.NewOrm()
	v = &Items{ItemId: id}
	if err = o.QueryTable(new(Items)).Filter("ItemId", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllItems retrieves all Items matches certain condition. Returns empty list if
// no records exist
func GetAllItems(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Items))
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

	var l []Items
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

// UpdateItems updates Items by Id and returns error if
// the record to be updated doesn't exist
func UpdateItemsById(m *Items) (err error) {
	o := orm.NewOrm()
	v := Items{ItemId: m.ItemId}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteItems deletes Items by Id and returns error if
// the record to be deleted doesn't exist
func DeleteItems(id int64) (err error) {
	o := orm.NewOrm()
	v := Items{ItemId: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Items{ItemId: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
