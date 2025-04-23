package request

import (
	"net/url"
	"testing"
)

// TestURLParamsOverride 测试 Get 请求中URL原始参数是否会丢失
func TestURLParamsOverride(t *testing.T) {
	testCases := []struct {
		name          string
		urlWithParams string
		params        url.Values
		expectedQuery string
	}{
		{
			name:          "URL参数被params覆盖",
			urlWithParams: "https://example.com/api?key1=url1&key2=url2",
			params: url.Values{
				"key3": []string{"params3"},
				"key4": []string{"params4"},
			},
			expectedQuery: "key1=url1&key2=url2&key3=params3&key4=params4",
		},
		{
			name:          "仅URL参数",
			urlWithParams: "https://example.com/api?key1=url1&key2=url2",
			params:        nil,
			expectedQuery: "key1=url1&key2=url2",
		},
		{
			name:          "仅params参数",
			urlWithParams: "https://example.com/api",
			params: url.Values{
				"key1": []string{"params1"},
				"key2": []string{"params2"},
			},
			expectedQuery: "key1=params1&key2=params2",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := NewClient().Get(tc.urlWithParams, tc.params)
			if req.err != nil {
				t.Fatalf("创建请求失败: %v", req.err)
			}

			gotQuery := req.request.URL.RawQuery
			if gotQuery != tc.expectedQuery {
				t.Errorf("期望的查询参数为 %s，但得到 %s", tc.expectedQuery, gotQuery)
			}
		})
	}
}
