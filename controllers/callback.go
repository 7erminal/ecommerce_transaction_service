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

// CallbackController operations for Callback
type CallbackController struct {
	beego.Controller
}

// URLMapping ...
func (c *CallbackController) URLMapping() {
	c.Mapping("Post", c.Post)
}

// Post ...
// @Title Create
// @Description create Callback
// @Param	body		body 	requests.CallbackRequest	true		"body for Callback content"
// @Success 201 {object} responses.CallbackResponse
// @Failure 403 body is empty
// @router /process [post]
func (c *CallbackController) Post() {
	var v requests.CallbackRequest
	logs.Info("Received callback request: ", string(c.Ctx.Input.RequestBody))
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.Data["json"] = map[string]string{"error": "Invalid request body"}
		c.Ctx.Output.SetStatus(400)
		c.ServeJSON()
		return
	}

	responseCode := false
	responseMessage := "Invalid request"
	transaction := responses.Bil_transactionCustom{}

	vBytes, _ := json.Marshal(v)
	logs.Info("Transaction callback Request data: ", string(vBytes))

	// Handle successful callback
	transactionRef := ""
	if v.ClientReference != nil {
		logs.Info("Transaction ID found in request: ", *v.ClientReference)
		transactionRef = *v.ClientReference
	}
	logs.Info("About to get transaction by ID: ", transactionRef)
	// id, err := strconv.ParseInt(transactionId, 10, 64)
	// if err != nil {
	// 	logs.Error("Invalid transaction ID: %v", err)
	// 	responseCode = false
	// 	responseMessage = "Invalid transaction ID"
	// 	resp := responses.CallbackResponse{
	// 		StatusCode:    responseCode,
	// 		StatusMessage: responseMessage,
	// 		Result:        nil,
	// 	}
	// 	c.Data["json"] = resp
	// 	c.Ctx.Output.SetStatus(400)
	// 	c.ServeJSON()
	// 	return
	// }
	if resp, err := models.GetBil_transactionsByTransactionRefNum(transactionRef); err == nil {
		logs.Info("Request ID: ", resp.TransactionId, " found for callback processing")
		if resp != nil {
			// Update the transaction status
			statusCode := v.Status
			logs.Info("Updating transaction status to: ", statusCode)

			status, err := models.GetStatus_codesByCode(statusCode)
			if err == nil {
				resp.Status = status
				resp.DateModified = time.Now()
				if v.ExternalTransactionId != nil {
					resp.ExternalReferenceNumber = *v.ExternalTransactionId
				}
				resp.DateModified = time.Now()
				resp.Charge = v.Charges
			} else {
				c.Data["json"] = map[string]string{"error": "Status code not found"}
				c.Ctx.Output.SetStatus(404)
			}

			if err := models.UpdateBil_transactionsById(resp); err != nil {
				logs.Info("Failed to update transaction status: %v", err)
				responseCode = false
				responseMessage = "Failed to update transaction status"
				resp := responses.CallbackResponse{
					StatusCode:    responseCode,
					StatusMessage: responseMessage,
					Result:        nil,
				}
				c.Data["json"] = resp
				c.Ctx.Output.SetStatus(200)
			} else {
				// c.Data["json"] = map[string]string{"message": "Transaction updated successfully"}

				// Update payment history
				var fields []string
				var sortby []string
				var order []string
				var query = make(map[string]string)
				var limit int64 = 10
				var offset int64

				querySearch := "BilTransactionId__TransactionId:" + strconv.FormatInt(resp.TransactionId, 10)
				// query: k:v,k:v
				if v := querySearch; v != "" {
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

				var insResponses []*responses.Bil_ins_transactionCustom

				insTxns, err := models.GetAllBil_ins_transactions(query, fields, sortby, order, offset, limit)
				if err == nil {
					for _, v := range insTxns {
						logs.Info("Payment history found: ", v)
						insTxnObj := v.(models.Bil_ins_transactions)

						insTxnObj.Status = status
						insTxnObj.DateModified = time.Now()

						if err := models.UpdateBil_ins_transactionsById(&insTxnObj); err != nil {
							logs.Error("Failed to update payment history: %v", err)
						} else {
							logs.Info("Payment history updated successfully for Payment ID: ", resp.TransactionId)
						}

						bilTxnResp := responses.Bil_ins_transactionCustom{
							BilInsTransactionId:    insTxnObj.BilInsTransactionId,
							Amount:                 strconv.FormatFloat(insTxnObj.Amount, 'f', 2, 64),
							Biller:                 insTxnObj.Biller.BillerName,
							SenderAccountNumber:    insTxnObj.SenderAccountNumber,
							RecipientAccountNumber: insTxnObj.RecipientAccountNumber,
							Network:                insTxnObj.Network,
							Request:                insTxnObj.Request,
							Response:               insTxnObj.Response,
							Active:                 insTxnObj.Active,
						}
						insResponses = append(insResponses, &bilTxnResp)
					}

				} else {
					logs.Error("Failed to retrieve payment history: %v", err)
				}

				// Update request with callback data
				resText, err := json.Marshal(v)
				if err != nil {
					logs.Error("Failed to marshal callback request: %v", err)
					// c.Data["json"] = "Invalid request format"
					// c.ServeJSON()
					// return
				}

				logs.Info("Callback response text: %s", string(resText))
				logs.Info("Updating request", resp.Request.RequestId, " with callback response")
				if request, err := models.GetRequestById(resp.Request.RequestId); err == nil {
					logs.Info("Found request: ", request.RequestId)
					request.CallbackResponse = string(resText)

					request.DateModified = time.Now()
					if err := models.UpdateRequestById(request); err != nil {
						logs.Error("Failed to update request: %v", err)
						// c.Data["json"] = "Failed to update request"
						// c.ServeJSON()
						// return
					} else {
						logs.Info("Request updated successfully with callback response")
					}
				} else {
					logs.Error("Failed to retrieve request by ID: %v", err)
				}

				responseCode = true
				responseMessage = "Transaction updated successfully"

				transaction = responses.Bil_transactionCustom{
					TransactionId:           strconv.FormatInt(resp.TransactionId, 10),
					TransactionRefNumber:    resp.TransactionRefNumber,
					Service:                 resp.Service.ServiceName,
					BillerCode:              resp.BillerCode,
					Amount:                  strconv.FormatFloat(resp.Amount, 'f', 2, 64),
					TransactingCurrency:     resp.TransactingCurrency,
					SourceChannel:           resp.SourceChannel,
					Source:                  resp.Source,
					Destination:             resp.Destination,
					Charge:                  strconv.FormatFloat(resp.Charge, 'f', 2, 64),
					Commission:              strconv.FormatFloat(resp.Commission, 'f', 2, 64),
					ExternalReferenceNumber: resp.ExternalReferenceNumber,
					Status:                  resp.Status.StatusDescription,
					ExtraDetails1:           resp.ExtraDetails1,
					ExtraDetails2:           resp.ExtraDetails2,
					ExtraDetails3:           resp.ExtraDetails3,
					DateProcessed:           resp.DateModified,
					Active:                  resp.Active,
					InsTxns:                 insResponses,
				}

				c.Ctx.Output.SetStatus(200)
			}
		} else {
			logs.Info("Transaction not found for ID: %s", transactionRef)
			responseCode = false
			responseMessage = "Transaction not found"
			// c.Data["json"] = map[string]string{"error": "Transaction not found"}
			c.Ctx.Output.SetStatus(200)
		}
	} else {
		c.Data["json"] = map[string]string{"error": "Failed to retrieve transaction"}
		logs.Info("Failed to retrieve transaction: %s", err.Error())
		responseCode = false
		responseMessage = "Failed to retrieve transaction"
		// c.Data["json"] = map[string]string{"error": "Transaction not found"}
		c.Ctx.Output.SetStatus(200)
	}

	resp := responses.CallbackResponse{
		StatusCode:    responseCode,
		StatusMessage: responseMessage,
		Result:        &transaction,
	}
	c.Data["json"] = resp

	c.ServeJSON()
}
