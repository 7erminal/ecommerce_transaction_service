package api

import (
	"net/http"
	"os"

	"github.com/astaxie/beego/httplib"
	"github.com/beego/beego/v2/core/logs"
)

type Client struct {
	Request *Request
	Type_   string
}

func (c *Client) SendRequest() (*http.Response, error) {
	// TODO: BaseURLとPathのvalidation
	fullURL := c.Request.BaseURL + c.Request.Path
	logs.Info("Full third party URL: ", fullURL)
	beegoRequest := httplib.NewBeegoRequest(
		fullURL,
		c.Request.Method.String())
	for key, value := range c.Request.HeaderField {
		beegoRequest.Header(key, value)
	}
	for key, value := range c.Request.FileField {
		logs.Info("File field seen...")
		logs.Info("Key is ", key, " and value is ", value)
		fileInfo, err := os.Stat(value)
		if err != nil {
			logs.Error("Error accessing file:", err)
		} else {
			logs.Info("File size:", fileInfo.Size())
		}
		beegoRequest.PostFile(key, value)
	}
	if c.Type_ == "params" {
		for key, value := range c.Request.Params {
			beegoRequest.Param(key, value)
		}
	} else if c.Type_ == "body" {
		if c.Request.InterfaceParams == nil {
			logs.Info("Interface params are nil")
			beegoRequest.JSONBody(c.Request.Params)
		} else {
			logs.Info("Interface params are not nil")
			beegoRequest.JSONBody(c.Request.InterfaceParams)
		}
	}

	beegoRequest.Debug(true)

	res, err := beegoRequest.Response()
	if err != nil {
		return nil, err
	}
	return res, nil
}
