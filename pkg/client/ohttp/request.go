package ohttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	onion_log "github.com/kms9/publicyc/pkg/onion-log"
	"github.com/kms9/publicyc/pkg/onion-log/logger"
)

type Requests struct {
	Config  *Config
	clients sync.Pool
	_logger *onion_log.Log
}

// Curl 请求使用baseContext
func (r *Requests) Curl(trace *logger.BaseContentInfo, project, api string, params map[string]interface{}) ([]byte, int, error) {
	req, exists := r.Config.Projects[strings.ToLower(project)]
	if !exists {
		return nil, 0, ErrProjectIsNull
	}
	config, exists := req.API[strings.ToLower(api)]
	if !exists {
		return nil, 0, ErrAPIIsNull
	}

	client := r.clients.Get().(*resty.Client)
	defer r.clients.Put(client)
	response, err := request(trace, client, req.URL, config, params)

	if err != nil || response == nil ||
		response.StatusCode() >= http.StatusBadRequest ||
		response.Time() > r.Config.SlowRequest || r.Config.Debug {

		r.print(response, req, config, err, project, api, params, trace)
	}

	if err != nil {
		return nil, 0, err
	}
	return response.Body(), response.StatusCode(), nil
}

// CurlWithGinContext 返回基于ginContext的curl方法
func (r *Requests) CurlWithGinContext(ctx *gin.Context) func(project, api string, params map[string]interface{}) ([]byte, int, error) {
	trace := onion_log.GetBaseByContext(ctx)
	return func(project, requestFunc string, params map[string]interface{}) ([]byte, int, error) {
		return r.Curl(trace, project, requestFunc, params)
	}
}

// print 打印错误信息及数据
func (r *Requests) print(response *resty.Response, req *ProjectConfig, config *RequestConfig, err error, project, api string, params map[string]interface{}, trace *logger.BaseContentInfo) {
	body := []byte{}
	code := 0
	latency := time.Duration(0)
	if response != nil {
		body = response.Body()
		code = response.StatusCode()
		latency = response.Time()
	}
	if len(body) > 2000 {
		body = body[0:2000]
	}

	loggerFields := r._logger.With(trace).WithFields(map[string]interface{}{
		"mod":         "request",
		"url":         req.URL + config.Path,
		"method":      config.Method,
		"code":        code,
		"params":      params,
		"body":        string(body),
		"err":         err,
		"project":     project,
		"requestFunc": api,
		"resLatency":  latency.Milliseconds(),
		"slow":        latency > r.Config.SlowRequest,
	})

	if r.Config.Debug {
		loggerFields.Info()
	} else {
		loggerFields.Warn()
	}
}

// Unmarshal 解析对应的结果
func (r *Requests) Unmarshal(b []byte, err error, result interface{}) error {
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, result)
	return err
}

// getParams 获取参数信息
func getParams(params map[string]interface{}, keys []string) (map[string]string, error) {
	result := map[string]string{}
	for _, v := range keys {
		if val, exists := params[v]; exists {
			result[v] = fmt.Sprintf("%v", val)
			continue
		}
		return nil, fmt.Errorf("params 缺少参数 %s", v)
	}
	return result, nil
}

// request 请求信息
func request(baseContent *logger.BaseContentInfo, client *resty.Client, url string, config *RequestConfig, params map[string]interface{}) (*resty.Response, error) {
	var err error

	if url == "" {
		return nil, ErrNoURL
	}

	// 超时设置: 秒
	if config.TimeOut != 0 {
		client.SetTimeout(config.TimeOut)
	} else {
		client.SetTimeout(defaultSetTimeout)
	}

	request := client.R()

	// url params 设置
	if len(config.PathParams) > 0 {
		pathParams, err := getParams(params, config.PathParams)
		if err != nil {
			return nil, err
		}
		request = request.SetPathParams(pathParams)
	}

	// query 设置
	if len(config.Query) > 0 {
		queryParams, err := getParams(params, config.Query)
		if err != nil {
			return nil, err
		}
		request = request.SetQueryParams(queryParams)
	}

	// header
	if len(config.Header) > 0 {
		headerParams, err := getParams(params, config.Header)
		if err != nil {
			return nil, err
		}
		request = request.SetHeaders(headerParams)
	}

	// 传递trace
	if baseContent != nil {
		request = request.SetHeaders(map[string]string{
			"trace-id": baseContent.TraceID,
			"app-from": baseContent.From,
			"span-id":  baseContent.SpanID,
		})
		if request.Header.Get("uid") == "" {
			request = request.SetHeader("uid", baseContent.UID)
		}
	}

	// body
	if body, exists := params["body"]; exists && config.BodyRequire {
		request = request.SetBody(body)
	}

	var response *resty.Response
	switch strings.ToUpper(config.Method) {
	case "GET":
		response, err = request.Get(url + config.Path)
	case "POST":
		response, err = request.Post(url + config.Path)
	case "PUT":
		response, err = request.Put(url + config.Path)
	case "DELETE":
		response, err = request.Delete(url + config.Path)
	case "PATCH":
		response, err = request.Patch(url + config.Path)
	case "HEAD":
		response, err = request.Head(url + config.Path)
	default:
		return nil, ErrNoMethod
	}

	return response, err
}

// Curl 请求
func Curl(trace *logger.BaseContentInfo, address, path, method string, header map[string]string, timeOut time.Duration, body interface{}) ([]byte, int, error) {
	config := &RequestConfig{
		TimeOut:     timeOut,
		Path:        url.QueryEscape(path),
		Method:      method,
		Header:      []string{},
		BodyRequire: body == nil,
	}
	params := map[string]interface{}{}
	for k, v := range header {
		config.Header = append(config.Header, k)
		params[k] = v
	}
	if config.BodyRequire {
		params["body"] = body
	}

	response, err := request(trace, resty.New(), address, config, params)
	if err != nil {
		return nil, 0, err
	}
	return response.Body(), response.StatusCode(), nil
}
