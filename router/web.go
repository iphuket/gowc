package router

import (
	"github.com/gin-gonic/gin"
	"github.com/iphuket/gowc/app/controller"
)

var wctr = new(controller.WeChat)
var lt = new(controller.LeifengTrend)

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
	rlt := app.Group("leifengtrend")
	rlt.Use(mugdeda)
	{
		rlt.Any("oauth/wechat/url", lt.WeChatLoginURL)
		rlt.Any("oauth/wechat/call", lt.WeChatLoginCall)
		rlt.Any("updataimage", lt.UpdataImage)
		rlt.Any("registertask", lt.RegisterTask)
		rlt.Any("checkstate", lt.CheckState)
		// https://gowc.iuu.pub/leifengtrend/img/
		rlt.StaticFS("img", gin.Dir("./img", true))
	}
	go app.RunTLS(":443", "./1885284_gowc.iuu.pub.pem", "./1885284_gowc.iuu.pub.key")
	app.Run(":80")
}

func mugdeda(c *gin.Context) {
	c.Header("access-control-allow-origin", "*")
	c.Next()
}
