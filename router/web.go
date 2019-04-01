package router

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iphuket/gowc/app/controller"
)

var wctr = new(controller.WeChat)
var lt = new(controller.LeifengTrend)
var rc = new(controller.RepresentCat)

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
	app.Delims("{{", "}}")
	app.LoadHTMLGlob("views/**/*")
	rrc := app.Group("representcat")
	{
		rrc.StaticFS("img", gin.Dir("statics/representcat/img", true))
		rrc.GET("page", rc.RepresentCatStart)
		rrc.GET("oauth", rc.RepresentCatURL)
		rrc.GET("call", rc.WeChatLoginCall)
		rrc.GET("config", rc.RepresentCatConfig)
		rrc.GET("vote", rc.RepresentCatVote)
		rrc.GET("getdata", rc.RepresentCatGetData)
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
func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d%02d/%02d", year, month, day)
}
