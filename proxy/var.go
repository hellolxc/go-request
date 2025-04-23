package proxy

import "fmt"

type ProxyType int

const (
	HTTP         ProxyType = iota + 1 // 0: HTTP
	SOCKS                             // 1: SOCKS
	HTTPAndSOCKS                      // 2: HTTP&SOCKS
)

func (p ProxyType) String() string {
	switch p {
	case HTTP:
		return "HTTP"
	case SOCKS:
		return "SOCKS"
	case HTTPAndSOCKS:
		return "HTTP&SOCKS"
	default:
		return "Unknown"
	}
}

func (p ProxyType) Int() *int {
	switch p {
	case HTTP:
		http := int(HTTP)
		return &http
	case SOCKS:
		socks := int(SOCKS)
		return &socks
	case HTTPAndSOCKS:
		httpAndSocks := int(HTTPAndSOCKS)
		return &httpAndSocks
	default:
		return nil
	}
}

type Proxy struct {
	Type     ProxyType
	Host     string
	Port     int
	User     string
	Password string
}

// IsNeedAuth 是否需要验证
func (p *Proxy) IsNeedAuth() bool {
	return len(p.User) > 0
}

// IsSocks 是否Socks
func (p *Proxy) IsSocks() bool {
	return p.Type == SOCKS || p.Type == HTTPAndSOCKS
}

func (p *Proxy) Address() string {
	if p.IsSocks() {
		return fmt.Sprintf("%s:%d", p.Host, p.Port)
	}

	if p.IsSocks() == false && p.IsNeedAuth() == false {
		return fmt.Sprintf("%s:%d", p.Host, p.Port)
	}

	// TODO 修改为HTTP 需要密码验证方式
	return fmt.Sprintf("%s:%d", p.Host, p.Port)
}
