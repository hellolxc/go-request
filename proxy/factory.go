package proxy

import (
	"errors"
	"net/http"
)

type ProxyFactory struct {
	proxyOrigin string
}

func (p *ProxyFactory) Get(config *Proxy) (httpTransport *http.Transport, err error) {
	if config.Type == SOCKS || config.Type == HTTPAndSOCKS {
		return new(Socks5Proxy).GetHttpTransport(config)
	}

	if config.Type == HTTP {
		return new(HttpProxy).GetHttpTransport(config)
	}

	err = errors.New("invalid proxy type")
	return
}
