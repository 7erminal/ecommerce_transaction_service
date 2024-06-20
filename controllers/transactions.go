package controllers

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"transaction_service/models"
	"transaction_service/structs/requests"
	"transaction_service/structs/responses"

	beego "github.com/beego/beego/v2/server/web"
)

// TransactionsController operations for Transactions
type TransactionsController struct {
	beego.Controller
}

// URLMapping ...
func (c *TransactionsController) URLMapping() {
	// c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create Transactions
// @Param	body		body 	models.Transactions	true		"body for Transactions content"
// @Success 201 {int} models.Transactions
// @Failure 403 body is empty
// @router / [post]
// func (c *TransactionsController) Post() {
// 	var v models.Transactions
// 	json.Unmarshal(c.Ctx.Input.RequestBody, &v)
// 	if _, err := models.AddTransactions(&v); err == nil {
// 		c.Ctx.Output.SetStatus(201)
// 		c.Data["json"] = v
// 	} else {
// 		c.Data["json"] = err.Error()
// 	}
// 	c.ServeJSON()
// }

// GetOne ...
// @Title Get One
// @Description get Transactions by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Transactions
// @Failure 403 :id is empty
// @router /:id [get]
func (c *TransactionsController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	v, err := models.GetTransactionsById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get Transactions
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Transactions
// @Failure 403
// @router / [get]
func (c *TransactionsController) GetAll() {
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

	l, err := models.GetAllTransactions(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Transactions
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	requests.UpdateTransactionRequestDTO	true		"body for Transactions content"
// @Success 200 {object} models.Transactions
// @Failure 403 :id is not int
// @router /:id [put]
func (c *TransactionsController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	var v requests.UpdateTransactionRequestDTO
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	if z, err := models.GetTransactionsById(id); err == nil {
		if td, err := models.GetTransaction_detailsByTransaction(z); err == nil {
			td.SenderAccountNumber = v.SenderAccountNumber
			td.RecipientAccountNumber = v.RecipientAccountNumber
			if err := models.UpdateTransaction_detailsById(td); err == nil {
				var resp = responses.TransactionResponseDTO{StatusCode: 200, Transaction: z, StatusDesc: "Transaction Updated"}
				c.Ctx.Output.SetStatus(201)
				c.Data["json"] = resp
			} else {
				var resp = responses.TransactionResponseDTO{StatusCode: 301, Transaction: z, StatusDesc: "Transaction Update failed. " + err.Error()}
				c.Ctx.Output.SetStatus(201)
				c.Data["json"] = resp
			}
		} else {
			var resp = responses.TransactionResponseDTO{StatusCode: 301, Transaction: z, StatusDesc: "Transaction Update failed. " + err.Error()}
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = resp
		}
	} else {
		var resp = responses.TransactionResponseDTO{StatusCode: 301, Transaction: nil, StatusDesc: "Transaction Update failed. " + err.Error()}
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = resp
	}

	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Transactions
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *TransactionsController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	if err := models.DeleteTransactions(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
