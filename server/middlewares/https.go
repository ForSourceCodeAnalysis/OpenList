package middlewares

import (
	"fmt"
	"strings"

	"github.com/OpenListTeam/OpenList/v4/internal/conf"
	"github.com/gin-gonic/gin"
)

func ForceHttps(c *gin.Context) {
	if c.Request.TLS == nil {
		host := c.Request.Host
		// change port to https port
		host = strings.Replace(host, fmt.Sprintf(":%d", conf.Conf.Scheme.HTTPPort), fmt.Sprintf(":%d", conf.Conf.Scheme.HTTPSPort), 1)
		c.Redirect(302, "https://"+host+c.Request.RequestURI)
		c.Abort()
		return
	}
	c.Next()
}
