package functions

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"
	"strings"
	"transaction_service/api"
	"transaction_service/structs/requests"
	"transaction_service/structs/responses"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

// containsIgnoreCase checks if substr is in s, case-insensitive.
func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) &&
		strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func setControllerJSON(c *beego.Controller, value interface{}) {
	if c == nil {
		return
	}
	if c.Data == nil {
		c.Data = map[interface{}]interface{}{}
	}
	c.Data["json"] = value
}

func GetUser(c *beego.Controller, req requests.GetUserRequest) (responses.UserResponseDTO, error) {
	host, _ := beego.AppConfig.String("customerBaseUrl")

	reqText, _ := json.Marshal(req)

	logs.Info("Request to get User: ", string(reqText))

	request := api.NewRequest(
		host,
		"/v1/users/"+req.UserId,
		api.GET)

	// request.Params = {"UserId": strconv.Itoa(int(userid))}
	client := api.Client{
		Request: request,
		Type_:   "params",
	}
	res, err := client.SendRequest()
	if err != nil {
		logs.Error("client.Error: %v", err)
		setControllerJSON(c, err.Error())
		return responses.UserResponseDTO{}, err
	}
	defer res.Body.Close()
	read, err := io.ReadAll(res.Body)
	if err != nil {
		setControllerJSON(c, err.Error())
		return responses.UserResponseDTO{}, err
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, read, "", "  "); err != nil {
		logs.Info("Raw response received is ", string(read))
	} else {
		logs.Info("Raw response received is \n", prettyJSON.String())
	}
	// data := map[string]interface{}{}
	// var dataOri responses.UserOriResponseDTO
	var data responses.UserResponseDTO
	if err := json.Unmarshal(read, &data); err != nil {
		setControllerJSON(c, err.Error())
		return responses.UserResponseDTO{}, err
	}
	setControllerJSON(c, data)

	logs.Info("Resp is ", data)
	// logs.Info("Resp is ", data.User.Branch.Country.DefaultCurrency)

	return data, nil
}

func GetUserWithUsername(c *beego.Controller, req requests.GetUserWithUsernameRequest) (responses.UserResponseDTO, error) {
	host, _ := beego.AppConfig.String("customerBaseUrl")

	reqText, _ := json.Marshal(req)

	logs.Info("Request to get User: ", string(reqText))

	request := api.NewRequest(
		host,
		"/v1/users/get-user-by-username/"+req.Username,
		api.GET)

	// request.Params = {"UserId": strconv.Itoa(int(userid))}
	client := api.Client{
		Request: request,
		Type_:   "params",
	}
	res, err := client.SendRequest()
	if err != nil {
		logs.Error("client.Error: %v", err)
		setControllerJSON(c, err.Error())
		return responses.UserResponseDTO{}, err
	}
	defer res.Body.Close()
	read, err := io.ReadAll(res.Body)
	if err != nil {
		setControllerJSON(c, err.Error())
		return responses.UserResponseDTO{}, err
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, read, "", "  "); err != nil {
		logs.Info("Raw response received is ", string(read))
	} else {
		logs.Info("Raw response received is \n", prettyJSON.String())
	}
	// data := map[string]interface{}{}
	// var dataOri responses.UserOriResponseDTO
	var data responses.UserResponseDTO
	if err := json.Unmarshal(read, &data); err != nil {
		setControllerJSON(c, err.Error())
		return responses.UserResponseDTO{}, err
	}
	setControllerJSON(c, data)

	logs.Info("Resp is ", data)
	// logs.Info("Resp is ", data.User.Branch.Country.DefaultCurrency)

	return data, nil
}

func GetCustomer(c *beego.Controller, req requests.GetCustomerRequest) (responses.CustomerResponseAltDTO, error) {
	host, _ := beego.AppConfig.String("customerBaseUrl")

	reqText, _ := json.Marshal(req)

	logs.Info("Request to get Customer: ", string(reqText))

	request := api.NewRequest(
		host,
		"/v1/customers/"+req.CustomerId,
		api.GET)

	// request.Params = {"UserId": strconv.Itoa(int(userid))}
	client := api.Client{
		Request: request,
		Type_:   "params",
	}
	res, err := client.SendRequest()
	if err != nil {
		logs.Error("client.Error: %v", err)
		setControllerJSON(c, err.Error())
		return responses.CustomerResponseAltDTO{}, err
	}
	defer res.Body.Close()
	read, err := io.ReadAll(res.Body)
	if err != nil {
		setControllerJSON(c, err.Error())
		return responses.CustomerResponseAltDTO{}, err
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, read, "", "  "); err != nil {
		logs.Info("Raw response received is ", string(read))
	} else {
		logs.Info("Raw response received is \n", prettyJSON.String())
	}
	// data := map[string]interface{}{}
	// var dataOri responses.UserOriResponseDTO
	var data responses.CustomerResponseAltDTO
	if err := json.Unmarshal(read, &data); err != nil {
		setControllerJSON(c, err.Error())
		return responses.CustomerResponseAltDTO{}, err
	}
	setControllerJSON(c, data)

	logs.Info("Resp is ", data)
	// logs.Info("Resp is ", data.User.Branch.Country.DefaultCurrency)

	return data, nil
}

func GetBranch(c *beego.Controller, req requests.GetBranchRequest) (responses.BranchResponseDTO, error) {
	host, _ := beego.AppConfig.String("customerBaseUrl")

	reqText, _ := json.Marshal(req)

	logs.Info("Request to get Branch: ", string(reqText))

	request := api.NewRequest(
		host,
		"/v1/branches/"+req.BranchId,
		api.GET)

	// request.Params = {"UserId": strconv.Itoa(int(userid))}
	client := api.Client{
		Request: request,
		Type_:   "params",
	}
	res, err := client.SendRequest()
	if err != nil {
		logs.Error("client.Error: %v", err)
		setControllerJSON(c, err.Error())
		return responses.BranchResponseDTO{}, err
	}
	defer res.Body.Close()
	read, err := io.ReadAll(res.Body)
	if err != nil {
		setControllerJSON(c, err.Error())
		return responses.BranchResponseDTO{}, err
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, read, "", "  "); err != nil {
		logs.Info("Raw response received is ", string(read))
	} else {
		logs.Info("Raw response received is \n", prettyJSON.String())
	}
	// data := map[string]interface{}{}
	// var dataOri responses.UserOriResponseDTO
	var data responses.BranchResponseDTO
	if err := json.Unmarshal(read, &data); err != nil {
		setControllerJSON(c, err.Error())
		return responses.BranchResponseDTO{}, err
	}
	setControllerJSON(c, data)

	logs.Info("Resp is ", data)
	// logs.Info("Resp is ", data.User.Branch.Country.DefaultCurrency)

	return data, nil
}

func GetCustomerWithUsername(c *beego.Controller, req requests.GetCustomerWithUsernameRequest) (responses.CustomerResponseDTO, error) {
	host, _ := beego.AppConfig.String("customerBaseUrl")

	reqText, _ := json.Marshal(req)

	logs.Info("Request to get Customer with username: ", string(reqText))

	request := api.NewRequest(
		host,
		"/v1/customers/username/"+req.Username,
		api.GET)

	// request.Params = {"UserId": strconv.Itoa(int(userid))}
	client := api.Client{
		Request: request,
		Type_:   "params",
	}
	res, err := client.SendRequest()
	if err != nil {
		logs.Error("client.Error: %v", err)
		setControllerJSON(c, err.Error())
		return responses.CustomerResponseDTO{}, err
	}
	defer res.Body.Close()
	read, err := io.ReadAll(res.Body)
	if err != nil {
		setControllerJSON(c, err.Error())
		return responses.CustomerResponseDTO{}, err
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, read, "", "  "); err != nil {
		logs.Info("Raw response received is ", string(read))
	} else {
		logs.Info("Raw response received is \n", prettyJSON.String())
	}
	// data := map[string]interface{}{}
	// var dataOri responses.UserOriResponseDTO
	var data responses.CustomerResponseDTO
	if err := json.Unmarshal(read, &data); err != nil {
		setControllerJSON(c, err.Error())
		return responses.CustomerResponseDTO{}, err
	}
	setControllerJSON(c, data)

	logs.Info("Resp is ", data)
	// logs.Info("Resp is ", data.User.Branch.Country.DefaultCurrency)

	return data, nil
}

func UpdateUserPassword(c *beego.Controller, req requests.UpdateUserPasswordRequest) (responses.UserResponseDTO, error) {
	host, _ := beego.AppConfig.String("customerBaseUrl")

	reqText, _ := json.Marshal(req)

	logs.Info("Request to update User password: ", string(reqText))

	request := api.NewRequest(
		host,
		"/v1/users/password/"+strconv.FormatInt(req.UserId, 10),
		api.PUT)
	if request.InterfaceParams == nil {
		request.InterfaceParams = map[string]interface{}{}
	}

	request.InterfaceParams["UserId"] = req.UserId
	request.InterfaceParams["OldPassword"] = req.OldPassword
	request.InterfaceParams["NewPassword"] = req.NewPassword

	// request.Params = {"UserId": strconv.Itoa(int(userid))}
	client := api.Client{
		Request: request,
		Type_:   "params",
	}
	res, err := client.SendRequest()
	if err != nil {
		logs.Error("client.Error: %v", err)
		setControllerJSON(c, err.Error())
		return responses.UserResponseDTO{}, err
	}
	defer res.Body.Close()
	read, err := io.ReadAll(res.Body)
	if err != nil {
		setControllerJSON(c, err.Error())
		return responses.UserResponseDTO{}, err
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, read, "", "  "); err != nil {
		logs.Info("Raw response received is ", string(read))
	} else {
		logs.Info("Raw response received is \n", prettyJSON.String())
	}
	// data := map[string]interface{}{}
	// var dataOri responses.UserOriResponseDTO
	var data responses.UserResponseDTO
	if err := json.Unmarshal(read, &data); err != nil {
		setControllerJSON(c, err.Error())
		return responses.UserResponseDTO{}, err
	}
	setControllerJSON(c, data)

	logs.Info("Resp is ", data)
	// logs.Info("Resp is ", data.User.Branch.Country.DefaultCurrency)

	return data, nil
}

func GetCurrency(c *beego.Controller, req requests.GetCurrencyRequest) (responses.CurrencyResponseDTO, error) {
	host, _ := beego.AppConfig.String("systemBaseUrl")

	reqText, _ := json.Marshal(req)

	logs.Info("Request to get Customer: ", string(reqText))

	request := api.NewRequest(
		host,
		"/v1/currencies/"+req.CurrencyId,
		api.GET)

	// request.Params = {"UserId": strconv.Itoa(int(userid))}
	client := api.Client{
		Request: request,
		Type_:   "params",
	}
	res, err := client.SendRequest()
	if err != nil {
		logs.Error("client.Error: %v", err)
		setControllerJSON(c, err.Error())
		return responses.CurrencyResponseDTO{}, err
	}
	defer res.Body.Close()
	read, err := io.ReadAll(res.Body)
	if err != nil {
		setControllerJSON(c, err.Error())
		return responses.CurrencyResponseDTO{}, err
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, read, "", "  "); err != nil {
		logs.Info("Raw response received is ", string(read))
	} else {
		logs.Info("Raw response received is \n", prettyJSON.String())
	}
	// data := map[string]interface{}{}
	// var dataOri responses.UserOriResponseDTO
	var data responses.CurrencyResponseDTO
	if err := json.Unmarshal(read, &data); err != nil {
		setControllerJSON(c, err.Error())
		return responses.CurrencyResponseDTO{}, err
	}
	setControllerJSON(c, data)

	logs.Info("Resp is ", data)
	// logs.Info("Resp is ", data.User.Branch.Country.DefaultCurrency)

	return data, nil
}

func GetCountry(c *beego.Controller, req requests.GetCountryRequest) (responses.CountryResponseDTO, error) {
	host, _ := beego.AppConfig.String("systemBaseUrl")

	reqText, _ := json.Marshal(req)

	logs.Info("Request to get Country: ", string(reqText))

	request := api.NewRequest(
		host,
		"/v1/countries/"+req.CountryId,
		api.GET)

	// request.Params = {"UserId": strconv.Itoa(int(userid))}
	client := api.Client{
		Request: request,
		Type_:   "params",
	}
	res, err := client.SendRequest()
	if err != nil {
		logs.Error("client.Error: %v", err)
		setControllerJSON(c, err.Error())
		return responses.CountryResponseDTO{}, err
	}
	defer res.Body.Close()
	read, err := io.ReadAll(res.Body)
	if err != nil {
		setControllerJSON(c, err.Error())
		return responses.CountryResponseDTO{}, err
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, read, "", "  "); err != nil {
		logs.Info("Raw response received is ", string(read))
	} else {
		logs.Info("Raw response received is \n", prettyJSON.String())
	}
	// data := map[string]interface{}{}
	// var dataOri responses.UserOriResponseDTO
	var data responses.CountryResponseDTO
	if err := json.Unmarshal(read, &data); err != nil {
		setControllerJSON(c, err.Error())
		return responses.CountryResponseDTO{}, err
	}
	setControllerJSON(c, data)

	logs.Info("Resp is ", data)
	// logs.Info("Resp is ", data.User.Branch.Country.DefaultCurrency)

	return data, nil
}

func GetStatus(c *beego.Controller, req requests.GetStatusRequest) (responses.StatusResponseDTO, error) {
	host, _ := beego.AppConfig.String("systemBaseUrl")

	reqText, _ := json.Marshal(req)

	logs.Info("Request to get Status: ", string(reqText))

	request := api.NewRequest(
		host,
		"/v1/statuses/"+req.StatusId,
		api.GET)

	// request.Params = {"UserId": strconv.Itoa(int(userid))}
	client := api.Client{
		Request: request,
		Type_:   "params",
	}
	res, err := client.SendRequest()
	if err != nil {
		logs.Error("client.Error: %v", err)
		setControllerJSON(c, err.Error())
		return responses.StatusResponseDTO{}, err
	}
	defer res.Body.Close()
	read, err := io.ReadAll(res.Body)
	if err != nil {
		setControllerJSON(c, err.Error())
		return responses.StatusResponseDTO{}, err
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, read, "", "  "); err != nil {
		logs.Info("Raw response received is ", string(read))
	} else {
		logs.Info("Raw response received is \n", prettyJSON.String())
	}
	// data := map[string]interface{}{}
	// var dataOri responses.UserOriResponseDTO
	var data responses.StatusResponseDTO
	if err := json.Unmarshal(read, &data); err != nil {
		setControllerJSON(c, err.Error())
		return responses.StatusResponseDTO{}, err
	}
	setControllerJSON(c, data)

	logs.Info("Resp is ", data)
	// logs.Info("Resp is ", data.User.Branch.Country.DefaultCurrency)

	return data, nil
}

func GetItem(c *beego.Controller, req requests.GetItemRequest) (responses.ItemResponseDTO, error) {
	host, _ := beego.AppConfig.String("itemBaseUrl")

	reqText, _ := json.Marshal(req)

	logs.Info("Request to get Item: ", string(reqText))

	request := api.NewRequest(
		host,
		"/v1/items/"+req.ItemId,
		api.GET)

	// request.Params = {"UserId": strconv.Itoa(int(userid))}
	client := api.Client{
		Request: request,
		Type_:   "params",
	}
	res, err := client.SendRequest()
	if err != nil {
		logs.Error("client.Error: %v", err)
		setControllerJSON(c, err.Error())
		return responses.ItemResponseDTO{}, err
	}
	defer res.Body.Close()
	read, err := io.ReadAll(res.Body)
	if err != nil {
		setControllerJSON(c, err.Error())
		return responses.ItemResponseDTO{}, err
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, read, "", "  "); err != nil {
		logs.Info("Raw response received is ", string(read))
	} else {
		logs.Info("Raw response received is \n", prettyJSON.String())
	}
	// data := map[string]interface{}{}
	// var dataOri responses.UserOriResponseDTO
	var data responses.ItemResponseDTO
	if err := json.Unmarshal(read, &data); err != nil {
		setControllerJSON(c, err.Error())
		return responses.ItemResponseDTO{}, err
	}
	setControllerJSON(c, data)

	logs.Info("Resp is ", data)
	// logs.Info("Resp is ", data.User.Branch.Country.DefaultCurrency)

	return data, nil
}

func GetItemQuantity(c *beego.Controller, req requests.GetItemRequest) (responses.ItemQuantityResponseDTO, error) {
	host, _ := beego.AppConfig.String("itemBaseUrl")

	reqText, _ := json.Marshal(req)

	logs.Info("Request to get Item: ", string(reqText))

	request := api.NewRequest(
		host,
		"/v1/items/quantity/"+req.ItemId,
		api.GET)

	// request.Params = {"UserId": strconv.Itoa(int(userid))}
	client := api.Client{
		Request: request,
		Type_:   "params",
	}
	res, err := client.SendRequest()
	if err != nil {
		logs.Error("client.Error: %v", err)
		setControllerJSON(c, err.Error())
		return responses.ItemQuantityResponseDTO{}, err
	}
	defer res.Body.Close()
	read, err := io.ReadAll(res.Body)
	if err != nil {
		setControllerJSON(c, err.Error())
		return responses.ItemQuantityResponseDTO{}, err
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, read, "", "  "); err != nil {
		logs.Info("Raw response received is ", string(read))
	} else {
		logs.Info("Raw response received is \n", prettyJSON.String())
	}
	// data := map[string]interface{}{}
	// var dataOri responses.UserOriResponseDTO
	var data responses.ItemQuantityResponseDTO
	if err := json.Unmarshal(read, &data); err != nil {
		setControllerJSON(c, err.Error())
		return responses.ItemQuantityResponseDTO{}, err
	}
	setControllerJSON(c, data)

	logs.Info("Resp is ", data)
	// logs.Info("Resp is ", data.User.Branch.Country.DefaultCurrency)

	return data, nil
}

func UpdateItemQuantity(c *beego.Controller, req requests.UpdateItemQuantityRequest) (responses.ItemResponseDTO, error) {
	host, _ := beego.AppConfig.String("itemBaseUrl")

	reqText, _ := json.Marshal(req)

	logs.Info("Request to update Item quantity: ", string(reqText))

	request := api.NewRequest(
		host,
		"/v1/items/quantity/"+req.ItemId,
		api.PUT)

	// request.Params = {"UserId": strconv.Itoa(int(userid))}
	request.InterfaceParams["ItemId"] = req.ItemId
	request.InterfaceParams["Quantity"] = req.Quantity
	client := api.Client{
		Request: request,
		Type_:   "body",
	}
	res, err := client.SendRequest()
	if err != nil {
		logs.Error("client.Error: %v", err)
		setControllerJSON(c, err.Error())
		return responses.ItemResponseDTO{}, err
	}
	defer res.Body.Close()
	read, err := io.ReadAll(res.Body)
	if err != nil {
		setControllerJSON(c, err.Error())
		return responses.ItemResponseDTO{}, err
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, read, "", "  "); err != nil {
		logs.Info("Raw response received is ", string(read))
	} else {
		logs.Info("Raw response received is \n", prettyJSON.String())
	}
	// data := map[string]interface{}{}
	// var dataOri responses.UserOriResponseDTO
	var data responses.ItemResponseDTO
	if err := json.Unmarshal(read, &data); err != nil {
		setControllerJSON(c, err.Error())
		return responses.ItemResponseDTO{}, err
	}
	setControllerJSON(c, data)

	logs.Info("Resp is ", data)
	// logs.Info("Resp is ", data.User.Branch.Country.DefaultCurrency)

	return data, nil
}

func GetBiller(c *beego.Controller, req requests.GetBillerRequest) (responses.BillerResponseDTO, error) {
	host, _ := beego.AppConfig.String("systemBaseUrl")

	reqText, _ := json.Marshal(req)

	logs.Info("Request to get Item: ", string(reqText))

	request := api.NewRequest(
		host,
		"/v1/billers/"+req.BillerId,
		api.GET)

	// request.Params = {"UserId": strconv.Itoa(int(userid))}
	client := api.Client{
		Request: request,
		Type_:   "params",
	}
	res, err := client.SendRequest()
	if err != nil {
		logs.Error("client.Error: %v", err)
		setControllerJSON(c, err.Error())
		return responses.BillerResponseDTO{}, err
	}
	defer res.Body.Close()
	read, err := io.ReadAll(res.Body)
	if err != nil {
		setControllerJSON(c, err.Error())
		return responses.BillerResponseDTO{}, err
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, read, "", "  "); err != nil {
		logs.Info("Raw response received is ", string(read))
	} else {
		logs.Info("Raw response received is \n", prettyJSON.String())
	}
	// data := map[string]interface{}{}
	// var dataOri responses.UserOriResponseDTO
	var data responses.BillerResponseDTO
	if err := json.Unmarshal(read, &data); err != nil {
		setControllerJSON(c, err.Error())
		return responses.BillerResponseDTO{}, err
	}
	setControllerJSON(c, data)

	logs.Info("Resp is ", data)
	// logs.Info("Resp is ", data.User.Branch.Country.DefaultCurrency)

	return data, nil
}

func GetService(c *beego.Controller, req requests.GetServiceRequest) (responses.ServiceResponseDTO, error) {
	host, _ := beego.AppConfig.String("systemBaseUrl")

	reqText, _ := json.Marshal(req)

	logs.Info("Request to get Service: ", string(reqText))

	request := api.NewRequest(
		host,
		"/v1/services/"+req.ServiceId,
		api.GET)

	// request.Params = {"UserId": strconv.Itoa(int(userid))}
	client := api.Client{
		Request: request,
		Type_:   "params",
	}
	res, err := client.SendRequest()
	if err != nil {
		logs.Error("client.Error: %v", err)
		setControllerJSON(c, err.Error())
		return responses.ServiceResponseDTO{}, err
	}
	defer res.Body.Close()
	read, err := io.ReadAll(res.Body)
	if err != nil {
		setControllerJSON(c, err.Error())
		return responses.ServiceResponseDTO{}, err
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, read, "", "  "); err != nil {
		logs.Info("Raw response received is ", string(read))
	} else {
		logs.Info("Raw response received is \n", prettyJSON.String())
	}
	// data := map[string]interface{}{}
	// var dataOri responses.UserOriResponseDTO
	var data responses.ServiceResponseDTO
	if err := json.Unmarshal(read, &data); err != nil {
		setControllerJSON(c, err.Error())
		return responses.ServiceResponseDTO{}, err
	}
	setControllerJSON(c, data)

	logs.Info("Resp is ", data)
	// logs.Info("Resp is ", data.User.Branch.Country.DefaultCurrency)

	return data, nil
}

func GetOperator(c *beego.Controller, req requests.GetOperatorRequest) (responses.OperatorResponseDTO, error) {
	host, _ := beego.AppConfig.String("systemBaseUrl")

	reqText, _ := json.Marshal(req)

	logs.Info("Request to get Operator: ", string(reqText))

	request := api.NewRequest(
		host,
		"/v1/operators/"+req.OperatorId,
		api.GET)

	// request.Params = {"UserId": strconv.Itoa(int(userid))}
	client := api.Client{
		Request: request,
		Type_:   "params",
	}
	res, err := client.SendRequest()
	if err != nil {
		logs.Error("client.Error: %v", err)
		setControllerJSON(c, err.Error())
		return responses.OperatorResponseDTO{}, err
	}
	defer res.Body.Close()
	read, err := io.ReadAll(res.Body)
	if err != nil {
		setControllerJSON(c, err.Error())
		return responses.OperatorResponseDTO{}, err
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, read, "", "  "); err != nil {
		logs.Info("Raw response received is ", string(read))
	} else {
		logs.Info("Raw response received is \n", prettyJSON.String())
	}
	// data := map[string]interface{}{}
	// var dataOri responses.UserOriResponseDTO
	var data responses.OperatorResponseDTO
	if err := json.Unmarshal(read, &data); err != nil {
		setControllerJSON(c, err.Error())
		return responses.OperatorResponseDTO{}, err
	}
	setControllerJSON(c, data)

	logs.Info("Resp is ", data)
	// logs.Info("Resp is ", data.User.Branch.Country.DefaultCurrency)

	return data, nil
}

func GetServices(c *beego.Controller) (responses.ServicesResponseDTO, error) {
	host, _ := beego.AppConfig.String("itemBaseUrl")

	logs.Info("Request to get Service")

	request := api.NewRequest(
		host,
		"/v1/services/",
		api.GET)

	// request.Params = {"UserId": strconv.Itoa(int(userid))}
	client := api.Client{
		Request: request,
		Type_:   "params",
	}
	res, err := client.SendRequest()
	if err != nil {
		logs.Error("client.Error: %v", err)
		setControllerJSON(c, err.Error())
		return responses.ServicesResponseDTO{}, err
	}
	defer res.Body.Close()
	read, err := io.ReadAll(res.Body)
	if err != nil {
		setControllerJSON(c, err.Error())
		return responses.ServicesResponseDTO{}, err
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, read, "", "  "); err != nil {
		logs.Info("Raw response received is ", string(read))
	} else {
		logs.Info("Raw response received is \n", prettyJSON.String())
	}
	// data := map[string]interface{}{}
	// var dataOri responses.UserOriResponseDTO
	var data responses.ServicesResponseDTO
	if err := json.Unmarshal(read, &data); err != nil {
		setControllerJSON(c, err.Error())
		return responses.ServicesResponseDTO{}, err
	}
	setControllerJSON(c, data)

	logs.Info("Resp is ", data)
	// logs.Info("Resp is ", data.User.Branch.Country.DefaultCurrency)

	return data, nil
}

func GetBillers(c *beego.Controller) (responses.BillersResponseDTO, error) {
	host, _ := beego.AppConfig.String("systemBaseUrl")

	logs.Info("Request to get Billers")

	request := api.NewRequest(
		host,
		"/v1/billers/",
		api.GET)

	// request.Params = {"UserId": strconv.Itoa(int(userid))}
	client := api.Client{
		Request: request,
		Type_:   "params",
	}
	res, err := client.SendRequest()
	if err != nil {
		logs.Error("client.Error: %v", err)
		setControllerJSON(c, err.Error())
		return responses.BillersResponseDTO{}, err
	}
	defer res.Body.Close()
	read, err := io.ReadAll(res.Body)
	if err != nil {
		setControllerJSON(c, err.Error())
		return responses.BillersResponseDTO{}, err
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, read, "", "  "); err != nil {
		logs.Info("Raw response received is ", string(read))
	} else {
		logs.Info("Raw response received is \n", prettyJSON.String())
	}
	// data := map[string]interface{}{}
	// var dataOri responses.UserOriResponseDTO
	var data responses.BillersResponseDTO
	if err := json.Unmarshal(read, &data); err != nil {
		setControllerJSON(c, err.Error())
		return responses.BillersResponseDTO{}, err
	}
	setControllerJSON(c, data)

	logs.Info("Resp is ", data)
	// logs.Info("Resp is ", data.User.Branch.Country.DefaultCurrency)

	return data, nil
}
