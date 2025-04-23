## 安装
```bash
go get github.com/hellolxc/go-request@latest
```

## 使用
```go
type Data struct {
    Code int    `json:"code"`
    Msg  string `json:"msg"`
}

var response Data

// 设置重试和超时时间
request := NewClient()SetRetry(3).SetRetryWaitTime(time.Minute)

// 设置 header
request.SetHeader("Origin", "https://test.com")

// Get 请求
params := url.Values{}
params.Add("page", "1")
_, err := request.Get("test.com?key=param", params).DoWithStruct(&response)

// Post 请求
_, err = request.Post("test.com", Data{Code: 1}).DoWithStruct(&response)
```