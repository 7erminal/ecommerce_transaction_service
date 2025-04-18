package functions

import (
	"encoding/json"
	"io"
	"transaction_service/api"
	"transaction_service/structs/responses"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

func UpdateCustomer(c *beego.Controller, customerid string, transactionDate string) (resp responses.CustomerResponseDTO) {
	host, _ := beego.AppConfig.String("customerBaseUrl")

	logs.Info("Sending user name ", customerid)

	request := api.NewRequest(
		host,
		"/v1/customers/last-txn/"+customerid,
		api.PUT)
	request.Params["TransactionDate"] = transactionDate
	// request.Params["Dob"] = req.Dob
	// request.Params["Gender"] = req.Gender
	// request.Params["PhoneNumber"] = req.PhoneNumber
	// request.Params["Username"] = req.Username
	// request.Params["MaritalStatus"] = ""
	// request.Params = {"UserId": strconv.Itoa(int(userid))}
	client := api.Client{
		Request: request,
		Type_:   "params",
	}
	res, err := client.SendRequest()
	if err != nil {
		logs.Error("client.Error: %v", err)
		c.Data["json"] = err.Error()
	}
	defer res.Body.Close()
	read, err := io.ReadAll(res.Body)
	if err != nil {
		c.Data["json"] = err.Error()
	}

	logs.Info("Raw response received is ", res)
	// data := map[string]interface{}{}
	var data responses.CustomerResponseDTO
	json.Unmarshal(read, &data)
	c.Data["json"] = data

	logs.Info("Resp is ", data)

	return data
}

func CheckItemAfterOrder(itemId string) (resp responses.CustomerResponseDTO) {
	host, _ := beego.AppConfig.String("itemsBaseUrl")

	logs.Info("Sending item ID ", itemId)

	var data responses.CustomerResponseDTO

	data.StatusCode = 406
	data.Customer = nil
	data.StatusDesc = "Something went wrong"

	request := api.NewRequest(
		host,
		"/v1/items/check-item-quantity/"+itemId,
		api.GET)
	// request.Params["Dob"] = req.Dob
	// request.Params["Gender"] = req.Gender
	// request.Params["PhoneNumber"] = req.PhoneNumber
	// request.Params["Username"] = req.Username
	// request.Params["MaritalStatus"] = ""
	// request.Params = {"UserId": strconv.Itoa(int(userid))}
	client := api.Client{
		Request: request,
		Type_:   "params",
	}
	res, err := client.SendRequest()
	if err != nil {
		logs.Error("client.Error: %v", err)
		// c.Data["json"] = err.Error()
	}
	defer res.Body.Close()
	read, err := io.ReadAll(res.Body)
	if err != nil {
		logs.Error("An error occurred when reading body ", err.Error())
		// c.Data["json"] = err.Error()
	}

	logs.Info("Raw response received is ", res)
	// data := map[string]interface{}{}
	json.Unmarshal(read, &data)
	// c.Data["json"] = data

	logs.Info("Resp is ", data)

	return data
}
