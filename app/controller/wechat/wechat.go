package wechat

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/iphuket/gowc/config"

	"github.com/iphuket/wechat/cache"

	"github.com/gin-gonic/gin"
	"github.com/iphuket/wechat"
	"github.com/iphuket/wechat/message"
	"github.com/iphuket/wechat/platforms"
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

// PlatformsInfo ... 开发平台数据库
type PlatformsInfo struct {
	UUID         string    `xorm:"varchar(36) pk notnull unique 'uuid'"`  // uuid
	AppID        string    `xorm:"varchar(64) pk notnull unique 'appid'"` // appid
	Name         string    `xorm:"varchar(64) 'name'"`                    // 名称
	AccessToken  string    `xorm:"varchar(64) 'access_token'"`            // 令牌
	RefreshToken string    `xorm:"varchar(64) 'refresh_token'"`           // 刷新令牌
	ExpiresIn    int64     `xorm:"varchar(64) 'expires_in'"`              // 令牌过期时间
	CreatedAt    int64     `xorm:"created"`
	UpdatedAt    time.Time `xorm:"updated"`
	DeletedAt    time.Time `xorm:"deleted"`
	Version      int64     `xorm:"version"`
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
	dbPlatforms := new(PlatformsInfo)
	err = engine.Sync2(dbPlatforms)
	if err != nil {
		x := fmt.Sprintf("%s", err)
		c.JSON(200, gin.H{"errCode": "error", "info": "engine Sync2 error " + x})
		return
	}
	// 获取 component_verify_ticket
	cvt := componentverifyticket()
	if cvt == "" {
		fmt.Println("获取不到缓存中的 component_verify_ticket")
		c.Writer.WriteString("获取不到缓存中的 component_verify_ticket")
		return
	}

	resComponentAccessToken, err := platforms.ComponentAccessToken(cfg.AppID, cfg.AppSecret, cvt)
	if err != nil {
		fmt.Println(err)
		c.Writer.WriteString("获取 ComponentAccessToken 出现错误, 错误信息：" + resComponentAccessToken.ErrMsg)
		return
	}

	abi, err := platforms.AuthBaseInfo(cfg.AppID, resComponentAccessToken.ComponentAccessToken, c.Request.FormValue("auth_code"))
	if err != nil {
		c.Writer.WriteString("platforms.AuthBaseInfo error ")
		fmt.Println(err)
		return
	}
	if len(abi.ErrMsg) > 0 {
		c.Writer.WriteString("获取 AuthBaseInfo 出现错误, 错误信息：" + abi.ErrMsg)
		return
	}
	dbPlatforms.UUID = uuid.New().String()
	dbPlatforms.AppID = abi.AuthorizationInfo.AuthorizerAppID
	dbPlatforms.Name = "财湘俱乐部"
	dbPlatforms.AccessToken = abi.AuthorizationInfo.AuthorizerAccessToken
	dbPlatforms.RefreshToken = abi.AuthorizationInfo.AuthorizerRefreshToken
	dbPlatforms.ExpiresIn = abi.AuthorizationInfo.ExpiresIn
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
	// 获取 component_verify_ticket
	cvt := componentverifyticket()
	if cvt == "" {
		fmt.Println("获取不到缓存中的 component_verify_ticket")
		c.Writer.WriteString("获取不到缓存中的 component_verify_ticket")
		return
	}

	resComponentAccessToken, err := platforms.ComponentAccessToken(cfg.AppID, cfg.AppSecret, cvt)
	if err != nil {
		fmt.Println(err)
		c.Writer.WriteString("获取 ComponentAccessToken 出现错误, 错误信息：" + resComponentAccessToken.ErrMsg)
		return
	}

	resPreAuthCode, err := platforms.PreAuthCode(cfg.AppID, resComponentAccessToken.ComponentAccessToken)
	if err != nil {
		fmt.Println(err)
		c.Writer.WriteString("获取 PreAuthCode 出现错误, 错误信息：" + resPreAuthCode.ErrMsg)
		return
	}

	URL := fmt.Sprintf(platforms.AuthURL, cfg.AppID, resPreAuthCode.PreAuthCode, RedirectURL)
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
	db := new(config.DB)
	engine, err := db.NewEngine()
	if err != nil {
		fmt.Println(err)
		x := fmt.Sprintf("%s", err)
		c.JSON(200, gin.H{"errCode": "error", "info": "new engine error " + x})
		return
	}
	dbPlatforms := new(PlatformsInfo)
	err = engine.Sync2(dbPlatforms)
	if err != nil {
		x := fmt.Sprintf("%s", err)
		c.JSON(200, gin.H{"errCode": "error", "info": "engine Sync2 error " + x})
		return
	}
	bool, err := engine.Where("appid=?", c.Request.FormValue("appid")).Get(dbPlatforms)
	if err != nil {
		c.JSON(200, gin.H{"errCode": "error", "info": "engine.Where error"})
		return
	} else if !bool {
		c.JSON(200, gin.H{"errCode": "error", "info": "not find appid"})
		return
	}
	c.JSON(200, gin.H{"errCode": "success", "name": dbPlatforms.Name, "uuid": dbPlatforms.UUID})
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
