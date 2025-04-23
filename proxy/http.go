package proxy

import (
	"fmt"
	"net/http"
	"net/url"
)

type HttpProxy struct {
}

func (h *HttpProxy) GetHttpTransport(p *Proxy) (httpTransport *http.Transport, err error) {
	var proxyURL *url.URL
	proxyURL, err = url.Parse(fmt.Sprintf("http://%s:%s@%s:%d", p.User, p.Password, p.Host, p.Port))
	if err != nil {
		fmt.Println("Invalid proxy URL:", err)
		return
	}

	httpTransport = &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	return
}
