package middlewares

import (
	"fmt"
	"github.com/ConfWatch/confwatchd/config"
	"github.com/ConfWatch/confwatchd/log"
	"github.com/gin-gonic/gin"
	"gopkg.in/unrolled/secure.v1"
	"strings"
)

func Security() gin.HandlerFunc {
	var rules *secure.Secure

	if config.Conf.Dev {
		rules = secure.New(secure.Options{
			FrameDeny:          true,
			ContentTypeNosniff: true,
			BrowserXssFilter:   true,
			ReferrerPolicy:     "same-origin",
		})
	} else {
		rules = secure.New(secure.Options{
			AllowedHosts:          config.Conf.Hosts,
			SSLRedirect:           true,
			SSLHost:               config.Conf.Hosts[0],
			STSSeconds:            315360000,
			STSIncludeSubdomains:  true,
			STSPreload:            true,
			FrameDeny:             true,
			ContentTypeNosniff:    true,
			BrowserXssFilter:      true,
			ContentSecurityPolicy: "default-src 'self'",
			ReferrerPolicy:        "same-origin",
		})
	}

	return func(c *gin.Context) {
		err := rules.Process(c.Writer, c.Request)
		if err != nil {
			who := strings.Split(c.Request.RemoteAddr, ":")[0]
			req := fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path)
			log.Warningf("%s > %s | Security exception: %s", who, req, err)
			c.Abort()
			return
		}

		if status := c.Writer.Status(); status > 300 && status < 399 {
			c.Abort()
		}
	}
}
