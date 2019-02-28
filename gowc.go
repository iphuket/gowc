package gowc

import (
	"github.com/gin-gonic/gin"
	"github.com/iphuket/gowc/app/middleware"

	"github.com/iphuket/gowc/router"
)

// Start GinApp
func Start() {
	GinApp := gin.Default()
	GinApp.Use(middleware.GOWC)
	router.WEB(GinApp)
}
