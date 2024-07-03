package server

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSAllowAll sets up the CORS headers to allow all origins, methods and headers.
func CORSAllowAll(c *gin.Context) {
	c.Header("Access-Control-Allow-Methods", "POST,OPTIONS,GET,PUT,DELETE,PATCH")
	c.Header("Access-Control-Allow-Credentials", "true")

	origin := c.GetHeader("origin")
	if len(strings.TrimSpace(origin)) == 0 {
		origin = "*"
	}
	c.Header("Access-Control-Allow-Origin", origin)

	ac := c.GetHeader("Access-Control-Request-Headers")
	if len(strings.TrimSpace(ac)) == 0 {
		ac = "*"
	}
	c.Header("Access-Control-Allow-Headers", ac)
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

	c.Next()
}
