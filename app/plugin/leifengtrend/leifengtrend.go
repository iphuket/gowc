package leifengtrend

import (
	"github.com/gin-gonic/gin"
	"github.com/iphuket/gowc/app/plugin/leifengtrend/controller"
	"github.com/iphuket/gowc/app/plugin/leifengtrend/middleware"
)

// Controller ... ctr
type Controller struct {
}

// Middleware ... mid
type Middleware struct {
}

// Register ...
func (ctr *Controller) Register(c *gin.Context) {
	controller.Register(c)
}

// Auth ...
func (mid *Middleware) Auth(c *gin.Context) {
	middleware.Auth(c)
}
