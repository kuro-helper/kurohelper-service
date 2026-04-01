package proxy

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"

	"golang.org/x/net/proxy"
)

var (
	// ErrCreateSOCKS5DialerFailed 建立 SOCKS5 Dialer 失敗
	ErrCreateSOCKS5DialerFailed = errors.New("kurohelperproxy: failed to create SOCKS5 dialer")
)

var (
	dialerInstance proxy.Dialer
	mu             sync.RWMutex
)

// get proxy dialer instance
func GetProxyDialer(addr, port string, auth *proxy.Auth) (proxy.Dialer, error) {
	mu.RLock()
	if dialerInstance != nil {
		defer mu.RUnlock()
		return dialerInstance, nil
	}
	mu.RUnlock()

	mu.Lock()
	defer mu.Unlock()

	if dialerInstance != nil {
		return dialerInstance, nil
	}

	fullAddr := net.JoinHostPort(addr, port)

	dialer, err := proxy.SOCKS5("tcp", fullAddr, auth, &net.Dialer{
		Timeout: 10 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCreateSOCKS5DialerFailed, err)
	}

	dialerInstance = dialer
	slog.Info("Proxy已成功設置", "address", fullAddr)

	return dialerInstance, nil
}

func GenerateProxyAuth(user, pwd string) *proxy.Auth {
	if user == "" && pwd == "" {
		return nil
	}
	return &proxy.Auth{User: user, Password: pwd}
}
