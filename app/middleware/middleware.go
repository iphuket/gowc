package middleware

import (
	"github.com/gin-gonic/gin"
)

// GOWC ...
func GOWC(c *gin.Context) {
	c.Header("sever", "gowc/1.0")
	c.Next()
}
