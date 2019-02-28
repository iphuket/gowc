package router

import (
	"github.com/gin-gonic/gin"
	"github.com/iphuket/gowc/app/controller"
)

var ctr = new(controller.Controller)

// WEB ...
func WEB(app *gin.Engine) {
	app.Any("SaDcNh3pRG.txt", ctr.WeChat.VerifyFile)
	wc := app.Group("/wechat/public/account")
	{
		wc.Any("auth_event", ctr.WeChat.AuthEvent)
		wc.Any("message_with_event/:appid", ctr.WeChat.MessageWithEvent)
		wc.Any("auth_call", ctr.WeChat.AuthCall)
		wc.Any("auth_url", ctr.WeChat.AuthURL)
		wc.Any("test", ctr.WeChat.Test)
	}
	app.Run(":80")
}
