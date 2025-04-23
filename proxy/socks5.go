package proxy

import (
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/net/proxy"
	"net"
	"net/http"
)

type Socks5Proxy struct {
}

func (s *Socks5Proxy) GetHttpTransport(config *Proxy) (httpTransport *http.Transport, err error) {

	dialer, err := s.getDialer(config)
	httpTransport = &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		},
	}

	return
}

func (s *Socks5Proxy) getDialer(config *Proxy) (dialer proxy.Dialer, err error) {
	var auth *proxy.Auth

	// 需要密码登陆
	if config.User != "" && config.Password != "" {
		auth = &proxy.Auth{
			User:     config.User,
			Password: config.Password,
		}
	}

	socks5Proxy := fmt.Sprintf("%s:%d", config.Host, config.Port)

	if dialer, err = proxy.SOCKS5("tcp", socks5Proxy, auth, proxy.Direct); err != nil {
		fmt.Println("Error creating SOCKS5 dialer:", err)
	}

	return
}
