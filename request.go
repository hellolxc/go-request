package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"syscall"
	"time"

	"github.com/hellolxc/go-request/proxy"
)

func NewClient() *Request {
	r := &Request{
		Client:  http.DefaultClient,
		headers: make(map[string]string),
	}

	// 初始化公众header信息
	r.initHeaders()

	return r
}

type Request struct {
	Client        *http.Client
	debug         bool              // 开启调试 打印更多信息
	request       *http.Request     // 每次请求都会重新生成
	headers       map[string]string // 请求 header
	timeout       time.Duration     // 超时时间
	retry         int               // 请求重试次数
	retryWaitTime time.Duration     // 重试等待时间
	cookies       []*http.Cookie    // cookies 数据
	err           error
}

func (r *Request) initHeaders() *Request {
	r.headers["Accept-Encoding"] = "gzip, deflate, br, zstd"
	r.headers["Accept-Language"] = "en"

	r.headers["Connection"] = "keep-alive"

	r.headers["Pragma"] = "no-cache"
	r.headers["Cache-Control"] = "no-cache"

	return r
}

// SetRetry 设置重试次数
func (r *Request) SetRetry(num int) *Request {
	r.retry = num
	return r
}

// SetRetryWaitTime 设置重试等待时间
func (r *Request) SetRetryWaitTime(waitTime time.Duration) *Request {
	r.retryWaitTime = waitTime
	return r
}

// HasHeader 检查请求头是否存在
func (r *Request) HasHeader(key string) bool {
	_, exists := r.headers[key]
	return exists
}

// SetHeader 设置请求头
func (r *Request) SetHeader(key string, value string) *Request {
	r.headers[key] = value
	return r
}

// SetHeaders 批量设置请求头
func (r *Request) SetHeaders(headers map[string]string) *Request {
	for k, v := range headers {
		r.headers[k] = v
	}

	return r
}

// SetUserAgent 设置浏览器/客户端标识
func (r *Request) SetUserAgent(userAgent string) *Request {
	r.headers["User-Agent"] = userAgent
	return r
}

// SetCookies 设置 cookie 数据
func (r *Request) SetCookies(cookies []*http.Cookie) *Request {
	r.cookies = cookies

	return r
}

// 将 Header 写入到请求中
func (r *Request) writeHeader() {
	for k, v := range r.headers {
		r.request.Header.Add(k, v)
	}
}

// SetTransport 设置Transport
func (r *Request) SetTransport(transport *http.Transport) *Request {
	r.Client.Transport = transport
	return r
}

// Proxy 设置代理
func (r *Request) Proxy(config *proxy.Proxy) *Request {
	if config == nil {
		return r
	}

	r.Client.Transport, r.err = new(proxy.ProxyFactory).Get(config)
	return r
}

// Head 请求
func (r *Request) Head(url string) *Request {
	if r.err != nil {
		return r
	}

	if r.request, r.err = http.NewRequest("HEAD", url, nil); r.err != nil {
		return r
	}

	return r
}

// Get 请求
func (r *Request) Get(url string, params url.Values) *Request {

	if r.err != nil {
		return r
	}

	if r.request, r.err = http.NewRequest("GET", url, nil); r.err != nil {
		return r
	}

	if len(params) > 0 {
		var originParams string
		if r.request.URL.RawQuery != "" {
			originParams = r.request.URL.RawQuery + "&"
		}

		r.request.URL.RawQuery = originParams + params.Encode()
	}

	return r
}

func (r *Request) Post(url string, body interface{}) *Request {
	if r.err != nil {
		return r
	}

	// 需要发送数据
	bodyReader := r.handleRequestBody(body)
	if r.request, r.err = http.NewRequest("POST", url, bodyReader); r.err != nil {
		return r
	}

	return r
}

// handleRequestBody
func (r *Request) handleRequestBody(body interface{}) (reader io.Reader) {
	if body == nil {
		return
	}

	// url values 类型
	if values, ok := body.(url.Values); ok {
		reader = strings.NewReader(values.Encode())
		r.request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		return
	}

	// 默认为JSON数据
	r.SetHeader("Content-Type", "application/json")
	reqBody, _ := json.Marshal(body)
	reader = bytes.NewReader(reqBody)

	return
}

func (r *Request) Do() (response *http.Response, err error) {
	if r.err != nil {
		return nil, r.err
	}

	r.writeHeader()

	// 处理 Cookie 数据
	if len(r.cookies) > 0 {
		for i := 0; i < len(r.cookies); i++ {
			r.request.AddCookie(r.cookies[i])
		}
	}

	// 增加自定义处理错误回调
	for retry := r.retry; retry > 0; retry-- {
		// 重试等待时间
		if err != nil && r.retryWaitTime > 0 {
			time.Sleep(r.retryWaitTime)
		}

		if response, err = r.Client.Do(r.request); err == nil {
			return
		}

		// 打印调试信息
		r.Debug(err.Error())

		if errors.Is(err, io.EOF) {
			continue
		}

		if errors.Is(err, syscall.ETIMEDOUT) {
			continue
		}

	}

	return
}

func (r *Request) DoWithStruct(data interface{}) (*http.Response, error) {
	if r.err != nil {
		return nil, r.err
	}

	response, err := r.Do()
	if err != nil {
		return nil, err
	}

	if data != nil {
		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			return response, err
		}

		if err := json.Unmarshal(body, &data); err != nil {
			return response, err
		}
	}

	return response, r.err
}

func (r *Request) Debug(message string) {
	if !r.debug {
		return
	}

	content := fmt.Sprintf("请求方式：%s; 请求地址: %s; ", r.request.Method, r.request.URL)
	fmt.Println(content + message)
}
