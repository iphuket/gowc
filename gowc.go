package gowc

import (
	"github.com/gin-gonic/gin"
	"github.com/iphuket/gowc/app/middleware"
	"github.com/iphuket/gowc/router"
)

// Start GinApp
func Start(port string) {
	GinApp := gin.Default()
	GinApp.Use(middleware.GOWC)
	router.WEB(GinApp)
	GinApp.Run(port)
}
