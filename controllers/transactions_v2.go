package controllers

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
	"transaction_service/models"
	"transaction_service/structs/requests"
	"transaction_service/structs/responses"

	beego "github.com/beego/beego/v2/server/web"
)

// TransactionsController operations for Transactions
type TransactionsV2Controller struct {
	beego.Controller
}

// URLMapping ...
func (c *TransactionsV2Controller) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
	c.Mapping("GetOneTransactionWithTxnRef", c.GetOneTransactionWithTxnRef)
}

// Post ...
// @Title Post
// @Description create Transactions
// @Param	body		body 	models.Transactions	true		"body for Transactions content"
// @Success 201 {int} models.Transactions
// @Failure 403 body is empty
// @router / [post]
func (c *TransactionsV2Controller) Post() {
	var req requests.BilTransactionRequestDTO
	json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	// Validate the request

	// authorization := ctx.Input.Header("Authorization")
	// phoneNumber := c.Ctx.Input.Header("PhoneNumber")
	phoneNumber := req.PhoneNumber
	sourceSystem := c.Ctx.Input.Header("SourceSystem")

	responseCode := 400
	responseMessage := "Request not processed"
	bilTxn := models.Bil_transactions{}

	statusCode := "PENDING" // Assuming 5002 is the status code for "Request Pending"

	reqText, err := json.Marshal(req)
	if err != nil {
		c.Data["json"] = "Invalid request format"
		c.ServeJSON()
		return
	}

	// Get customer by ID
	if cust, err := models.GetCustomerByPhoneNumber(phoneNumber); err == nil {
		status, err := models.GetStatus_codesByCode(statusCode)
		if err == nil {
			// Restructure the request to match the model
			serviceCode := req.ServiceCode
			if service, err := models.GetServicesByCode(serviceCode); err == nil {
				requestIdStr := req.RequestId
				requestId, _ := strconv.ParseInt(requestIdStr, 10, 64)
				// Create a request record
				v := models.Request{
					ApiRequestId:    requestId,
					CustId:          cust,
					Request:         string(reqText),
					RequestType:     service.ServiceName,
					RequestStatus:   status.StatusDescription,
					RequestAmount:   req.Amount,
					RequestResponse: "",
					RequestDate:     time.Now(),
					DateCreated:     time.Now(),
					DateModified:    time.Now(),
				}
				if _, err := models.AddRequest(&v); err == nil {
					if biller, err := models.GetBillerByCode(req.BillerCode); err == nil {
						// Create a transaction record
						transaction := models.Bil_transactions{
							TransactionRefNumber: "TRX-" + strconv.FormatInt(time.Now().Unix(), 10) + strconv.FormatInt(v.RequestId, 10),
							Service:              service, // Assuming service ID is 1 for airtime
							BillerCode:           biller.BillerCode,
							Request:              &v,
							TransactionBy:        cust,
							Amount:               req.Amount,
							TransactingCurrency:  "GHC", // Assuming USD for simplicity
							SourceChannel:        sourceSystem,
							Source:               req.Source,
							Destination:          req.Destination,
							Charge:               0.0,    // Assuming no charge for simplicity
							Status:               status, // Assuming 1 means successful
							ExtraDetails1:        req.ExtraData.ExtraData1,
							ExtraDetails2:        req.ExtraData.ExtraData2,
							ExtraDetails3:        req.ExtraData.ExtraData3,
							DateCreated:          time.Now(),
							DateModified:         time.Now(),
							CreatedBy:            1,
							ModifiedBy:           1,
							Active:               1, // Assuming active status
						}
						if _, err := models.AddBil_transactions(&transaction); err == nil {
							responseCode = 201
							responseMessage = "Transaction created successfully"
							bilTxn = transaction
						} else {
							responseMessage = "Failed to create transaction: " + err.Error()
							responseCode = 500
						}
					} else {
						responseMessage = "Biller not found: " + err.Error()
						responseCode = 502
					}
				}
			} else {
				responseMessage = "Service not found: " + err.Error()
				responseCode = 501
			}
		} else {
			responseMessage = "Status not found: " + err.Error()
			responseCode = 503
		}
	} else {
		responseMessage = "Customer not found: " + err.Error()
		responseCode = 504
	}

	response := responses.BilTransactionResponseDTO{
		StatusCode: responseCode,
		StatusDesc: responseMessage,
		Result:     &bilTxn,
	}

	c.Data["json"] = response

	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get Transactions by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Transactions
// @Failure 403 :id is empty
// @router /:id [get]
func (c *TransactionsV2Controller) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)

	statusCode := 400
	responseMessage := "Request not processed"
	bilTxn := models.Bil_transactions{}
	v, err := models.GetBil_transactionsById(id)
	if err != nil {
		responseMessage = "Transaction not found: " + err.Error()
		c.Ctx.Output.SetStatus(200)
	} else {
		bilTxn = *v
		statusCode = 200
		responseMessage = "Transaction retrieved successfully"
	}

	response := responses.BilTransactionResponseDTO{
		StatusCode: statusCode,
		StatusDesc: responseMessage,
		Result:     &bilTxn,
	}
	c.Data["json"] = response
	c.ServeJSON()
}

// GetOneTransactionWithTxnRef ...
// @Title Get One Transaction With Txn Ref
// @Description get Transactions by Txn Ref
// @Param	ref		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Transactions
// @Failure 403 :ref is empty
// @router /ref/:ref [get]
func (c *TransactionsV2Controller) GetOneTransactionWithTxnRef() {
	idStr := c.Ctx.Input.Param(":ref")
	// id, _ := strconv.ParseInt(idStr, 0, 64)

	statusCode := 400
	responseMessage := "Request not processed"
	bilTxn := models.Bil_transactions{}
	v, err := models.GetBil_transactionsByTransactionRefNum(idStr)
	if err != nil {
		responseMessage = "Transaction not found: " + err.Error()
		c.Ctx.Output.SetStatus(200)
	} else {
		bilTxn = *v
		statusCode = 200
		responseMessage = "Transaction retrieved successfully"
	}

	response := responses.BilTransactionResponseDTO{
		StatusCode: statusCode,
		StatusDesc: responseMessage,
		Result:     &bilTxn,
	}
	c.Data["json"] = response
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
func (c *TransactionsV2Controller) GetAll() {
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

	l, err := models.GetAllBil_transactions(query, fields, sortby, order, offset, limit)
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
// @Param	body		body 	models.Transactions	true		"body for Transactions content"
// @Success 200 {object} models.Transactions
// @Failure 403 :id is not int
// @router /:id [put]
func (c *TransactionsV2Controller) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	v := models.Bil_transactions{TransactionId: id}
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)
	if err := models.UpdateBil_transactionsById(&v); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
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
func (c *TransactionsV2Controller) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	if err := models.DeleteBil_transactions(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
