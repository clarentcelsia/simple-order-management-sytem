package util

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	MResponse "restaurant/model/response"

	"github.com/go-resty/resty/v2"
)

type client struct {
	*resty.Client
}

var instantiated *client
var once sync.Once

func NewRestyClient() *client {
	once.Do(func() {
		instantiated = &client{resty.New()}
		instantiated.SetTimeout(5 * time.Minute)
	})
	return instantiated
}

func HandleError(c *gin.Context, db interface{}, message string, detail string, dbstatus int) {
	response := MResponse.PageResponse{
		Status:  dbstatus,
		Message: message,
		Items:   db,
	}

	//serialize the given struct into json response body
	c.JSON(dbstatus, response)
	log.WithFields(log.Fields{
		"detail": detail,
	}).Error(db)
}

func HandleSuccess(c *gin.Context, db interface{}, message string, dbstatus int) {
	response := MResponse.PageResponse{
		Status:  dbstatus,
		Message: message,
		Items:   db,
	}
	c.JSON(dbstatus, response)
	log.WithFields(log.Fields{
		"detail": message,
	}).Info(message)
}

func HttpPostReq(c *gin.Context, url string, headers map[string]string, queryParams map[string]string, formData map[string]string, body map[string]interface{}) (map[string]interface{}, []byte, int, error) {

	var result interface{}
	res, errHttpReq := NewRestyClient().R().
		SetHeaders(headers).
		SetQueryParams(queryParams).
		SetFormData(formData).
		SetBody(body).
		Post(url)

	if errHttpReq != nil || res.StatusCode() != 200 {
		return nil, res.Body(), res.StatusCode(), errors.New("error Http Request to " + url)
	}

	err := json.Unmarshal(res.Body(), &result)
	if err != nil {
		return nil, res.Body(), res.StatusCode(), errors.New("error unmarshal to interface")
	}

	mapResponse := result.(map[string]interface{})
	return mapResponse, nil, res.StatusCode(), nil
}

func HttpGetReq(c *gin.Context, url string, headers map[string]string, queryParams map[string]string, formData map[string]string, body map[string]interface{}) (map[string]interface{}, []byte, int, error) {

	var result interface{}
	res, errHttpReq := NewRestyClient().R().
		SetHeaders(headers).
		SetQueryParams(queryParams).
		SetFormData(formData).
		SetBody(body).
		Get(url)

	if errHttpReq != nil || res.StatusCode() != 200 {
		return nil, res.Body(), res.StatusCode(), errors.New("error Http Request to " + url)
	}

	err := json.Unmarshal(res.Body(), &result)
	if err != nil {
		return nil, res.Body(), res.StatusCode(), errors.New("error unmarshal to interface")
	}

	mapResponse := result.(map[string]interface{})
	return mapResponse, nil, res.StatusCode(), nil
}
