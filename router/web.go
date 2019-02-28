package router

import (
	"github.com/gin-gonic/gin"
	"github.com/iphuket/gowc/app/controller"
)

var wctr = new(controller.WeChat)

// WEB ...
func WEB(app *gin.Engine) {
	app.Any("SaDcNh3pRG.txt", wctr.WeChat.VerifyFile)
	wc := app.Group("/wechat/public/account")
	{
		wc.Any("auth_event", wctr.WeChat.AuthEvent)
		wc.Any("message_with_event/:appid", wctr.WeChat.MessageWithEvent)
		wc.Any("auth_call", wctr.WeChat.AuthCall)
		wc.Any("auth_url", wctr.WeChat.AuthURL)
		wc.Any("test", wctr.WeChat.Test)
	}
	app.Run(":80")
}
