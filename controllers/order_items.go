package controllers

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"transaction_service/models"
	"transaction_service/structs/requests"
	"transaction_service/structs/responses"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

// Order_itemsController operations for Order_items
type Order_itemsController struct {
	beego.Controller
}

// URLMapping ...
func (c *Order_itemsController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create Order_items
// @Param	body		body 	models.Order_items	true		"body for Order_items content"
// @Success 201 {int} models.Order_items
// @Failure 403 body is empty
// @router / [post]
func (c *Order_itemsController) Post() {
	var v models.Order_items
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)
	if _, err := models.AddOrder_items(&v); err == nil {
		c.Ctx.Output.SetStatus(201)
		c.Data["json"] = v
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get Order_items by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Order_items
// @Failure 403 :id is empty
// @router /:id [get]
func (c *Order_itemsController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	v, err := models.GetOrder_itemsById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get Order_items
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Order_items
// @Failure 403
// @router / [get]
func (c *Order_itemsController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllOrder_items(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Order_items
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Order_items	true		"body for Order_items content"
// @Success 200 {object} responses.OrderItemsCustom
// @Failure 403 :id is not int
// @router /:id [put]
func (c *Order_itemsController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	v := requests.UpdateOrderItemDTO{OrderItemId: id}
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)
	if orderItem, err := models.GetOrder_itemsById(id); err != nil {
		logs.Info("Order item not found: ", err.Error())

	} else {
		logs.Info("Order item found: ", orderItem)
		logs.Info("Order item status: ", v.Status)
		if status, err := models.GetStatusByName(v.Status); err == nil {
			orderItem.Status = status
			orderItem.Comment = v.Comment
			orderItem.ModifiedBy = v.ModifiedBy
			if err := models.UpdateOrder_itemsById(orderItem); err == nil {
				logs.Info("Order item updated successfully: ", orderItem)
				if item, err := models.GetItemsById(orderItem.Item.ItemId); err != nil {
					logs.Error("An Error occurred while getting item: ", err.Error())
				} else {
					item.ItemStatus = status
					if err := models.UpdateItemsById(item); err != nil {
						logs.Error("An Error occurred while updating item: ", err.Error())
					}
				}
				o := responses.OrdersCustom{
					OrderId:      orderItem.Order.OrderId,
					OrderNumber:  orderItem.Order.OrderNumber,
					Quantity:     orderItem.Order.Quantity,
					Cost:         orderItem.Order.Cost,
					CurrencyId:   orderItem.Order.Currency,
					OrderDate:    orderItem.Order.OrderDate,
					DateCreated:  orderItem.Order.DateCreated,
					DateModified: orderItem.Order.DateModified,
				}
				oi := responses.OrderItemsCustom{
					OrderItemId: orderItem.OrderItemId,
					Order:       &o,
					Item:        orderItem.Item,
					Quantity:    orderItem.Quantity,
					Status:      orderItem.Status.Status,
					OrderDate:   orderItem.OrderDate,
					Comment:     orderItem.Comment,
				}
				resp := responses.OrderItemResponseDTO{StatusCode: 200, StatusDesc: "Order item updated successfully", OrderItem: &oi}
				c.Data["json"] = resp
			} else {
				// c.Data["json"] = err.Error()
				logs.Error("An Error occurred while updating order item: ", err.Error())
				resp := responses.OrderItemResponseDTO{StatusCode: 608, StatusDesc: "Order item update failed", OrderItem: nil}
				c.Data["json"] = resp
			}
		} else {
			logs.Error("An Error occurred while updating order item: ", err.Error())
			resp := responses.OrderItemResponseDTO{StatusCode: 608, StatusDesc: "Order item update failed. Status not found", OrderItem: nil}
			c.Data["json"] = resp
		}
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Order_items
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *Order_itemsController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	if err := models.DeleteOrder_items(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
