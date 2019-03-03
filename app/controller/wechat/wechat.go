package wechat

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/iphuket/gowc/app/model"

	"github.com/iphuket/gowc/config"

	"github.com/iphuket/wechat/cache"

	"github.com/gin-gonic/gin"
	"github.com/iphuket/wechat"
	"github.com/iphuket/wechat/context"
	"github.com/iphuket/wechat/message"
)

// WeChat 控制器
type WeChat struct {
}

// 接口信息
const (
	RedirectURL           = "http://am.jyacad.cc/wechat/public/account/auth_call"
	ComponentloginPageURL = "https://mp.weixin.qq.com/cgi-bin/componentloginpage?component_appid=%s&pre_auth_code=%s&redirect_uri=%s&auth_type=3"
)

var memCache = cache.NewMemcache("127.0.0.1:11211")

// cvt 专有缓存
var ca = new(config.Cache)

var aesKey = "MY4tM5jSFLlTj8l35cf"
var authRedirectURL = "http://am.oovmi.com/wechat/public/account/auth_call"

//配置微信参数
var cfg = &wechat.Config{
	AppID:          "wx45fc784e7eb235cf",
	AppSecret:      "22afb1d94c9c7ea0e8f403b8c7133305",
	Token:          "HekMkeMY4tM5jeO40jTOaTtSPn5slMEg",
	EncodingAESKey: "L9LSFLlTj8lYLy8Y8Q2D6668d6y982ef6Q0G26jdetb",
	Cache:          memCache,
}

// AuthCall ... 授权后回调地址
func (w *WeChat) AuthCall(c *gin.Context) {
	db := new(config.DB)
	engine, err := db.NewEngine()
	if err != nil {
		fmt.Println(err)
		x := fmt.Sprintf("%s", err)
		c.JSON(200, gin.H{"errCode": "error", "info": "new engine error " + x})
		return
	}
	dbPlatforms := new(model.WeChat)
	err = engine.Sync2(dbPlatforms)
	if err != nil {
		x := fmt.Sprintf("%s", err)
		c.JSON(200, gin.H{"errCode": "error", "info": "engine Sync2 error " + x})
		return
	}
	ctx := new(context.Context)
	// 必须先设置 ComponentAccessToken()
	_, err = ctx.GetComponentAccessToken()
	if err != nil {
		fmt.Println(err, "下一步正在进行 获取 ctv")
		cvt := componentverifyticket()
		_, err := ctx.SetComponentAccessToken(cvt)
		if err != nil {
			fmt.Println(err, "设置ComponentAccessToken 错误")
			return
		}
	}
	ac := c.Request.FormValue("auth_code")
	if len(ac) < 1 {
		fmt.Println("not find auth_code")
		c.Writer.WriteString("not find auth_code")
		return
	}
	abi, err := ctx.QueryAuthCode(ac)
	dbPlatforms.UUID = uuid.New().String()
	dbPlatforms.AppID = abi.Appid
	dbPlatforms.Name = "财湘俱乐部"
	dbPlatforms.AccessToken = abi.AccessToken
	dbPlatforms.RefreshToken = abi.RefreshToken
	dbPlatforms.ExpiresIn = abi.ExpiresIn
	_, err = engine.Insert(dbPlatforms)
	if err != nil {
		x := fmt.Sprintf("%s", err)
		c.JSON(200, gin.H{"errCode": "error", "info": "Insert error " + x})
		return
	}
	c.JSON(200, gin.H{"errCode": "success", "AppID": dbPlatforms.AppID})
}

// AuthURL ... 授权地址
func (w *WeChat) AuthURL(c *gin.Context) {
	ctx := new(context.Context)
	// 获取预授权码 必须先设置 ComponentAccessToken()
	_, err := ctx.GetComponentAccessToken()
	if err != nil {
		fmt.Println(err, "下一步正在进行 获取 ctv")
		cvt := componentverifyticket()
		_, err := ctx.SetComponentAccessToken(cvt)
		if err != nil {
			fmt.Println(err, "设置ComponentAccessToken 错误")
			return
		}
	}
	PreCode, err := ctx.GetPreCode()
	if err != nil {
		fmt.Println(err)
		c.Writer.WriteString("GetPreCode() error")
		return
	}
	URL := fmt.Sprintf(ComponentloginPageURL, cfg.AppID, PreCode, RedirectURL)
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(200, "<a href='"+URL+"'>点击授权</a>")

}

// MessageWithEvent ... 消息与事件接收url
func (w *WeChat) MessageWithEvent(c *gin.Context) {
	wc := wechat.NewWechat(cfg)
	server := wc.GetServer(c.Request, c.Writer)

	//设置接收消息的处理方法
	server.SetMessageHandler(messageWithEventHandler)
	//处理消息接收以及回复
	err := server.Serve()
	if err != nil {
		fmt.Println(errors.New("server.Serve() error: "))
		fmt.Println(err)
		return
	}
	//发送回复的消息
	err = server.Send()
	if err != nil {
		fmt.Println(errors.New("server.Send() error: "))
		fmt.Println(err)
		return
	}
}

// messageWithEventHandler  消息与事件处理
func messageWithEventHandler(msg message.MixMessage) *message.Reply {
	return nil
}

// Test 执行测试任务
func (w *WeChat) Test(c *gin.Context) {

}

// AuthEvent 授权事件接收URL
func (w *WeChat) AuthEvent(c *gin.Context) {

	wc := wechat.NewWechat(cfg)
	server := wc.GetServer(c.Request, c.Writer)
	//设置接收消息的处理方法
	server.SetMessageHandler(authEventHandler)
	//处理消息接收以及回复
	err := server.Serve()
	if err != nil {
		fmt.Println(errors.New("server.Serve() error: "))
		fmt.Println(err)
		return
	}
	//发送回复的消息
	err = server.Send()
	if err != nil {
		fmt.Println(errors.New("server.Send() error: "))
		fmt.Println(err)
		return
	}
}

// authEventHandler 授权事件处理
func authEventHandler(msg message.MixMessage) *message.Reply {
	if len(msg.ComponentVerifyTicket) > 0 {
		// 缓存这个参数
		err := ca.Set("ComponentVerifyTicket", msg.ComponentVerifyTicket)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		return nil
	}
	return nil
}

// VerifyFile 微信文件效验
func (w *WeChat) VerifyFile(c *gin.Context) {
	c.Writer.WriteString("65e10e6cda0f37d81cdffaf5a3441979")
}

// componentverifyticket
func componentverifyticket() string {

	str, err := ca.Get("ComponentVerifyTicket")
	if err != nil {
		return "error"
	}
	return str
}
