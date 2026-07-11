package gateway

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

type LoadBalancedProxy struct {
	proxies []*httputil.ReverseProxy
	next    atomic.Uint64
}

func NewLoadBalancedProxy(targetURLs []string) *LoadBalancedProxy {
	proxies := make([]*httputil.ReverseProxy, 0, len(targetURLs))

	for _, targetURL := range targetURLs {
		target, err := url.Parse(targetURL)
		if err != nil {
			panic(err)
		}

		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, proxyErr error) {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusBadGateway)
			_, _ = writer.Write([]byte(`{"code":502,"msg":"Upstream service unavailable","data":null}`))
		}

		proxies = append(proxies, proxy)
	}

	if len(proxies) == 0 {
		panic("at least one proxy target is required")
	}

	return &LoadBalancedProxy{
		proxies: proxies,
	}
}

func (lb *LoadBalancedProxy) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		proxy := lb.nextProxy()

		userID, _ := c.Get("user_id")
		username, _ := c.Get("username")
		role, _ := c.Get("role")

		// Setheader
		c.Request.Header.Set("KIT-DEV-USER-ID", fmt.Sprint(userID))
		c.Request.Header.Set("KIT-DEV-USERNAME", fmt.Sprint(username))
		c.Request.Header.Set("KIT-DEV-ROLE-ID", fmt.Sprint(role))

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func (lb *LoadBalancedProxy) nextProxy() *httputil.ReverseProxy {
	index := lb.next.Add(1)
	return lb.proxies[(index-1)%uint64(len(lb.proxies))]
}
