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

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

// TransactionsController operations for Transactions
type TransactionsV2Controller struct {
	beego.Controller
}

// URLMapping ...
func (c *TransactionsV2Controller) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("UserPost", c.UserPost)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
	c.Mapping("GetOneTransactionWithTxnRef", c.GetOneTransactionWithTxnRef)
	c.Mapping("GetAllUserTransactions", c.GetAllUserTransactions)
	c.Mapping("PutUserTransaction", c.PutUserTransaction)
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
		logs.Error("Invalid request format")
		c.Data["json"] = "Invalid request format"
		c.ServeJSON()
		return
	}

	logs.Info("Full request: %s", string(reqText))

	userid := req.CreatedBy
	useridInt, err := strconv.ParseInt(userid, 10, 64)
	if err != nil {
		useridInt = 1
	}

	// Get customer by ID
	cust := &models.Customers{}
	if cust, err = models.GetCustomerByPhoneNumber(phoneNumber); err != nil {
		logs.Error("Customer not found: ", err)
		responseMessage = "Customer not found: " + err.Error()
		responseCode = 504
	}
	logs.Info("Customer details: ", cust)

	logs.Info("Fetching status for code: %s", statusCode)
	status, err := models.GetStatus_codesByCode(statusCode)
	if err == nil {
		// Restructure the request to match the model
		serviceCode := req.ServiceCode
		logs.Info("Fetching service for code: %s", serviceCode)
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
					logs.Info("Extra data received are ", req.ExtraData.ExtraData1, req.ExtraData.ExtraData2, req.ExtraData.ExtraData3)

					// If user does not exist, create a system user with the userid 1
					if user, err := models.GetUsersById(useridInt); err == nil {

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
							Package:              req.Package,
							Charge:               0.0,    // Assuming no charge for simplicity
							Status:               status, // Assuming 1 means successful
							CorpId:               req.CorpId,
							ExtraDetails1:        req.ExtraData.ExtraData1,
							ExtraDetails2:        req.ExtraData.ExtraData2,
							ExtraDetails3:        req.ExtraData.ExtraData3,
							DateCreated:          time.Now(),
							DateModified:         time.Now(),
							CreatedBy:            user,
							ModifiedBy:           user,
							Active:               1, // Assuming active status
						}
						if _, err := models.AddBil_transactions(&transaction); err == nil {
							logs.Info("Transaction created successfully: ", transaction)
							responseCode = 200
							responseMessage = "Transaction created successfully"
							bilTxn = transaction
						} else {
							logs.Error("Failed to create transaction: ", err)
							responseMessage = "Failed to create transaction: " + err.Error()
							responseCode = 500
						}
					} else {
						logs.Error("User not found: ", err)
						responseMessage = "Failed to fetch user: " + err.Error()
						responseCode = 500
					}
				} else {
					logs.Error("Biller not found: ", err)
					responseMessage = "Biller not found: " + err.Error()
					responseCode = 502
				}
			} else {
				logs.Error("Failed to create request record: ", err)
				responseMessage = "Failed to create request record: " + err.Error()
				responseCode = 500
			}
		} else {
			logs.Error("Service not found: ", err)
			responseMessage = "Service not found: " + err.Error()
			responseCode = 501
		}
	} else {
		logs.Error("Status not found: ", err)
		responseMessage = "Status not found: " + err.Error()
		responseCode = 503
	}

	response := responses.BilTransactionResponseDTO{
		StatusCode: responseCode,
		StatusDesc: responseMessage,
		Result:     &bilTxn,
	}

	c.Data["json"] = response

	c.ServeJSON()
}

// UserPost ...
// @Title UserPost
// @Description create Transactions for a user
// @Param	body		body 	models.Transactions	true		"body for Transactions content"
// @Success 201 {int} models.Transactions
// @Failure 403 body is empty
// @router /user [post]
func (c *TransactionsV2Controller) UserPost() {
	var req requests.UserTransactionRequestDTO
	json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	// Validate the request

	// authorization := ctx.Input.Header("Authorization")
	// phoneNumber := c.Ctx.Input.Header("PhoneNumber")
	phoneNumber := req.PhoneNumber
	sourceSystem := c.Ctx.Input.Header("SourceSystem")

	responseCode := 400
	responseMessage := "Request not processed"
	txnData := responses.UserTransactions{}

	statusCode := req.Status

	reqText, err := json.Marshal(req)
	if err != nil {
		logs.Error("Invalid request format")
		c.Data["json"] = "Invalid request format"
		c.ServeJSON()
		return
	}

	logs.Info("Full request: %s", string(reqText))

	userid := req.CreatedBy
	useridInt, err := strconv.ParseInt(userid, 10, 64)
	if err != nil {
		useridInt = 1
	}

	// Get customer by ID
	user := &models.Users{}
	if user, err = models.GetUsersByUsername(phoneNumber); err != nil {
		logs.Error("User not found: ", err)
		responseMessage = "User not found: " + err.Error()
		responseCode = 504
	}
	logs.Info("User details: ", user)

	logs.Info("Fetching status for code: %s", statusCode)
	status, err := models.GetStatus_codesByCode(statusCode)
	if err == nil {
		// Restructure the request to match the model
		serviceCode := req.ServiceCode
		logs.Info("Fetching service for code: %s", serviceCode)
		if service, err := models.GetServicesByCode(serviceCode); err == nil {
			requestIdStr := req.RequestId
			requestId, _ := strconv.ParseInt(requestIdStr, 10, 64)
			// Create a request record
			v := models.UserRequest{
				ApiRequestId:    requestId,
				UserId:          user,
				Request:         string(reqText),
				RequestType:     service.ServiceName,
				RequestStatus:   status.StatusDescription,
				RequestAmount:   req.Amount,
				RequestResponse: "",
				RequestDate:     time.Now(),
				DateCreated:     time.Now(),
				DateModified:    time.Now(),
			}
			if _, err := models.AddUserRequest(&v); err == nil {
				logs.Info("Extra data received are ", req.ExtraData.ExtraData1, req.ExtraData.ExtraData2, req.ExtraData.ExtraData3)

				// Get customer by ID
				cust := &models.Customers{}
				if cust, err = models.GetCustomerByPhoneNumber(phoneNumber); err != nil {
					logs.Error("Customer not found: ", err)
					responseMessage = "Customer not found: " + err.Error()
					responseCode = 504
				}
				logs.Info("Customer details: ", cust)
				// If user does not exist, create a system user with the userid 1
				if user, err := models.GetUsersById(useridInt); err == nil {

					// Create a transaction record
					transaction := models.UserTransactions{
						TransactionId:                "TRX-" + strconv.FormatInt(time.Now().Unix(), 10) + strconv.FormatInt(v.RequestId, 10),
						Service:                      service, // Assuming service ID is 1 for airtime
						Request:                      &v,
						TransactionCustomerReference: cust,
						Amount:                       req.Amount,
						TransactingCurrency:          "GHC",
						Reference:                    req.Reference,
						SourceChannel:                sourceSystem,
						Source:                       req.SourceAccountNumber,
						Destination:                  req.DestinationAccountNumber,
						Package:                      req.Package,
						Charge:                       0.0,    // Assuming no charge for simplicity
						Status:                       status, // Assuming 1 means successful
						ExtraDetails1:                req.ExtraData.ExtraData1,
						ExtraDetails2:                req.ExtraData.ExtraData2,
						ExtraDetails3:                req.ExtraData.ExtraData3,
						DateCreated:                  time.Now(),
						DateModified:                 time.Now(),
						CreatedBy:                    user,
						ModifiedBy:                   user,
						Active:                       1, // Assuming active status
					}
					if _, err := models.AddUserTransactions(&transaction); err == nil {
						logs.Info("Transaction created successfully: ", transaction)
						// Save in user_ins_transactions table to add details in the steps
						userInsTransaction := models.UserInsTransactions{
							UserInsTransactionId:   "TRX-INS-" + strconv.FormatInt(time.Now().Unix(), 10) + strconv.FormatInt(v.RequestId, 10),
							UserTransactionId:      transaction.TransactionId,
							Amount:                 req.Amount,
							Data:                   "",
							SenderAccountNumber:    req.SourceAccountNumber,
							RecipientAccountNumber: req.DestinationAccountNumber,
							Service:                &models.Services{ServiceId: service.ServiceId},
							Status:                 &models.Status_codes{StatusId: status.StatusId},
							Request:                string(reqText),
							Response:               "",
							DateCreated:            time.Now(),
							DateModified:           time.Now(),
							CreatedBy:              int(user.UserId),
							ModifiedBy:             int(user.UserId),
							Active:                 1,
						}

						txnData = responses.UserTransactions{
							TransactionId:                transaction.TransactionId,
							Service:                      service.ServiceName,
							TransactionCustomerReference: cust.FullName,
							Amount:                       transaction.Amount,
							TransactingCurrency:          transaction.TransactingCurrency,
							SourceChannel:                transaction.SourceChannel,
							Source:                       transaction.Source,
							Destination:                  transaction.Destination,
							Package:                      transaction.Package,
							Charge:                       transaction.Charge,
							Commission:                   transaction.Commission,
							ExternalReferenceNumber:      transaction.ExternalReferenceNumber,
							Status:                       status.StatusDescription,
							ExtraDetails1:                transaction.ExtraDetails1,
							ExtraDetails2:                transaction.ExtraDetails2,
							ExtraDetails3:                transaction.ExtraDetails3,
							ClientResponseCode:           transaction.ClientResponseCode,
							DateCreated:                  transaction.DateCreated,
							DateModified:                 transaction.DateModified,
							CreatedBy:                    user.Username,
							ModifiedBy:                   user.Username,
							Active:                       transaction.Active,
						}
						if userInsTransactionJSON, marshalErr := json.Marshal(userInsTransaction); marshalErr != nil {
							logs.Error("failed to marshal userInsTransaction to json: %v", marshalErr)
						} else {
							logs.Info("userInsTransaction: %s", string(userInsTransactionJSON))
						}
						if _, err := models.AddUserInsTransactions(&userInsTransaction); err == nil {
							logs.Info("UserInsTransaction created successfully: ", userInsTransaction)

							responseCode = 200
							responseMessage = "Transaction created successfully"
						} else {
							logs.Error("Failed to create transaction: ", err)
							responseMessage = "Failed to create transaction: " + err.Error()
							responseCode = 500
						}
					} else {
						logs.Error("Failed to add user transaction: ", err)
						responseMessage = "Failed to add user transaction: " + err.Error()
						responseCode = 500
					}
				} else {
					logs.Error("User not found: ", err)
					responseMessage = "Failed to fetch user: " + err.Error()
					responseCode = 500
				}
			} else {
				logs.Error("Failed to create request record: ", err)
				responseMessage = "Failed to create request record: " + err.Error()
				responseCode = 500
			}
		} else {
			logs.Error("Service not found: ", err)
			responseMessage = "Service not found: " + err.Error()
			responseCode = 501
		}
	} else {
		logs.Error("Status not found: ", err)
		responseMessage = "Status not found: " + err.Error()
		responseCode = 503
	}

	response := responses.UserTransactionResponseDTO{
		StatusCode: responseCode,
		StatusDesc: responseMessage,
		Result:     &txnData,
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
		if vJSON, marshalErr := json.Marshal(v); marshalErr != nil {
			logs.Error("failed to marshal transaction v to json: %v", marshalErr)
		} else {
			logs.Info("transaction v: %s", string(vJSON))
		}

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

// GetUserTransactionWithTxnRef ...
// @Title Get User Transaction With Txn Ref
// @Description get Transactions by Txn Ref
// @Param	ref		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Transactions
// @Failure 403 :ref is empty
// @router /user/ref/:ref [get]
func (c *TransactionsV2Controller) GetUserTransactionWithTxnRef() {
	idStr := c.Ctx.Input.Param(":ref")
	// id, _ := strconv.ParseInt(idStr, 0, 64)

	statusCode := 400
	responseMessage := "Request not processed"
	bilTxn := models.UserTransactions{}
	txnData := responses.UserTransactions{}
	v, err := models.GetUserTransactionsById(idStr)
	if err != nil {
		responseMessage = "Transaction not found: " + err.Error()
		c.Ctx.Output.SetStatus(200)
	} else {
		bilTxn = *v
		if vJSON, marshalErr := json.Marshal(v); marshalErr != nil {
			logs.Error("failed to marshal transaction v to json: %v", marshalErr)
		} else {
			logs.Info("transaction v: %s", string(vJSON))
		}

		statusCode = 200
		responseMessage = "Transaction retrieved successfully"

		txnData = responses.UserTransactions{
			TransactionId:                bilTxn.TransactionId,
			Service:                      bilTxn.Service.ServiceName,
			TransactionCustomerReference: bilTxn.TransactionCustomerReference.FullName,
			Amount:                       bilTxn.Amount,
			TransactingCurrency:          bilTxn.TransactingCurrency,
			SourceChannel:                bilTxn.SourceChannel,
			Source:                       bilTxn.Source,
			Destination:                  bilTxn.Destination,
			Package:                      bilTxn.Package,
			Charge:                       bilTxn.Charge,
			Commission:                   bilTxn.Commission,
			ExternalReferenceNumber:      bilTxn.ExternalReferenceNumber,
			Status:                       bilTxn.Status.StatusCode,
			ExtraDetails1:                bilTxn.ExtraDetails1,
			ExtraDetails2:                bilTxn.ExtraDetails2,
			ExtraDetails3:                bilTxn.ExtraDetails3,
			ClientResponseCode:           bilTxn.ClientResponseCode,
			DateCreated:                  bilTxn.DateCreated,
			DateModified:                 bilTxn.DateModified,
			CreatedBy:                    bilTxn.CreatedBy.FullName,
			ModifiedBy:                   bilTxn.ModifiedBy.FullName,
			Active:                       bilTxn.Active,
		}
	}

	response := responses.UserTransactionResponseDTO{
		StatusCode: statusCode,
		StatusDesc: responseMessage,
		Result:     &txnData,
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

// GetAllUserTransactions ...
// @Title Get All User Transactions
// @Description get User Transactions
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Transactions
// @Failure 403
// @router /user-transactions [get]
func (c *TransactionsV2Controller) GetAllUserTransactions() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	responseCode := 400
	responseMessage := "Request not processed"
	txnData := []responses.UserTransactions{}

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

	l, err := models.GetAllUserTransactions(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = err.Error()
		logs.Error("Failed to retrieve transactions: ", err)
		responseMessage = "Failed to retrieve transactions: " + err.Error()
		responseCode = 500
	} else {
		// c.Data["json"] = l
		for _, v := range l {
			m := v.(models.UserTransactions)

			txnData = append(txnData, responses.UserTransactions{
				TransactionId:                m.TransactionId,
				Service:                      m.Service.ServiceName,
				TransactionCustomerReference: m.TransactionCustomerReference.FullName,
				Amount:                       m.Amount,
				TransactingCurrency:          m.TransactingCurrency,
				SourceChannel:                m.SourceChannel,
				Source:                       m.Source,
				Destination:                  m.Destination,
				Package:                      m.Package,
				Charge:                       m.Charge,
				Commission:                   m.Commission,
				ExternalReferenceNumber:      m.ExternalReferenceNumber,
				Status:                       m.Status.StatusDescription,
				ExtraDetails1:                m.ExtraDetails1,
				ExtraDetails2:                m.ExtraDetails2,
				ExtraDetails3:                m.ExtraDetails3,
				ClientResponseCode:           m.ClientResponseCode,
				DateCreated:                  m.DateCreated,
				DateModified:                 m.DateModified,
				CreatedBy:                    m.CreatedBy.FullName,
				ModifiedBy:                   m.ModifiedBy.FullName,
				Active:                       m.Active,
			})
		}
		responseCode = 200
		responseMessage = "Transactions retrieved successfully"
	}

	response := responses.UserTransactionsResponseDTO{
		StatusCode: responseCode,
		StatusDesc: responseMessage,
		Result:     &txnData,
	}

	c.Data["json"] = response

	c.ServeJSON()
}

// Put User Transactions ...
// @Title Put User Transactions
// @Description update the Transactions
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Transactions	true		"body for Transactions content"
// @Success 200 {object} models.Transactions
// @Failure 403 :id is not int
// @router /user-transaction/:id [put]
func (c *TransactionsV2Controller) PutUserTransaction() {
	idStr := c.Ctx.Input.Param(":id")
	req := requests.UpdateUserTransactionRequest{}
	json.Unmarshal(c.Ctx.Input.RequestBody, &req)

	statusCode := 400
	statusMessage := "Something went wrong"
	txnData := responses.UserTransactions{}

	status := req.Status

	if transaction, err := models.GetUserTransactionsById(idStr); err == nil {
		status, err := models.GetStatus_codesByCode(status)
		if err == nil {
			transaction.Status = status
			transaction.ExternalReferenceNumber = req.ClientReference
			if err := models.UpdateUserTransactionsById(transaction); err == nil {
				statusCode = 200
				statusMessage = "Transaction updated successfully"

				customerRef := ""
				if transaction.TransactionCustomerReference != nil {
					customerRef = transaction.TransactionCustomerReference.FullName
				}

				txnData = responses.UserTransactions{
					TransactionId:                transaction.TransactionId,
					TransactionCustomerReference: customerRef,
					Amount:                       transaction.Amount,
					TransactingCurrency:          transaction.TransactingCurrency,
					SourceChannel:                transaction.SourceChannel,
					Source:                       transaction.Source,
					Destination:                  transaction.Destination,
					Package:                      transaction.Package,
					Charge:                       transaction.Charge,
					Commission:                   transaction.Commission,
					ExternalReferenceNumber:      transaction.ExternalReferenceNumber,
					Status:                       transaction.Status.StatusDescription,
					ExtraDetails1:                transaction.ExtraDetails1,
					ExtraDetails2:                transaction.ExtraDetails2,
					ExtraDetails3:                transaction.ExtraDetails3,
					ClientResponseCode:           transaction.ClientResponseCode,
					DateCreated:                  transaction.DateCreated,
					DateModified:                 transaction.DateModified,
					CreatedBy:                    transaction.CreatedBy.FullName,
					ModifiedBy:                   transaction.ModifiedBy.FullName,
					Active:                       transaction.Active,
				}
			} else {
				c.Data["json"] = err.Error()
				logs.Error("Error updating user transaction:: ", err.Error())
				statusMessage = "Failed to update transaction: " + err.Error()
				statusCode = 500
			}
		} else {
			statusMessage = "Invalid status provided"
		}
	}

	response := responses.UserTransactionResponseDTO{
		StatusCode: statusCode,
		StatusDesc: statusMessage,
		Result:     &txnData,
	}

	c.Data["json"] = response

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
