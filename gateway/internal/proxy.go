package gateway

import (
	"fmt"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func ReverseProxy(targetUrl string) gin.HandlerFunc {
	target, err := url.Parse(targetUrl)
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(target)

	return func(c *gin.Context) {
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
