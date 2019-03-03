package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/iphuket/gowc/app/plugin/leifengtrend/jwt"
)

// Login use jwt login
func Login(c *gin.Context) {
	// 是否登陆
	uuid, err := jwt.Chcek(c.ClientIP(), c.GetHeader("token"))
	if err != nil {
		// ... user not login
		c.Writer.WriteString("user not login")
		return
	}
	c.Writer.WriteString(uuid)
}
