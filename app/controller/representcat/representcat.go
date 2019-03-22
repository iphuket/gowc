package representcat

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"github.com/iphuket/gowc/config"
	"github.com/iphuket/wechat/coauth"
	"github.com/iphuket/wechat/platforms"
)

// RepresentCat ...
type RepresentCat struct {
}

// Config struct
type Config struct {
	ComponentAppID          string
	ComponentAppSecret      string
	ComponentToken          string
	ComponentEncodingAESKey string
}

// 配置微信参数
var cfg = &Config{
	ComponentAppID:          "wx45fc784e7eb235cf",
	ComponentAppSecret:      "22afb1d94c9c7ea0e8f403b8c7133305",
	ComponentToken:          "HekMkeMY4tM5jeO40jTOaTtSPn5slMEg",
	ComponentEncodingAESKey: "L9LSFLlTj8lYLy8Y8Q2D6668d6y982ef6Q0G26jdetb",
}

// AppID ...
// RedirctURL ...
const (
	AppID      = "wx4285b8f2757f9774"
	RedirctURL = "https://gowc.iuu.pub/representcat/call"
)

// RepresentCatStart ...
func (rc *RepresentCat) RepresentCatStart(c *gin.Context) {

	c.HTML(http.StatusOK, "index.html", map[string]interface{}{"test": "test"})
	//c.Writer.WriteString("ok")
}

// RepresentCatURL 登陆页面
func (rc *RepresentCat) RepresentCatURL(c *gin.Context) {
	// c.JSON(200, gin.H{"errCode": "error", "info": "not find UUID error "})
	scope := "snsapi_userinfo"
	rccUUID := c.Request.FormValue("rcc_uuid")
	if len(rccUUID) < 1 {
		c.JSON(200, gin.H{"errCode": "error", "info": "not find rcc_uuid error "})
		return
	}
	state := rccUUID
	ComponentAppID := cfg.ComponentAppID
	url := coauth.GetCodeURL(AppID, RedirctURL, scope, state, ComponentAppID)
	rURL := []byte("<script>window.location.href='" + url + "';</script>")
	c.Writer.Write(rURL)
}

// WeChatLoginCall ... coauth call
func (rc *RepresentCat) WeChatLoginCall(c *gin.Context) {
	Code := c.Request.FormValue("code")
	if len(Code) < 1 {
		c.JSON(200, gin.H{"errCode": "error", "info": "not find code error "})
		return
	}
	rccUUID := c.Request.FormValue("state")
	if len(rccUUID) < 1 {
		c.JSON(200, gin.H{"errCode": "error", "info": "not find state error "})
		return
	}
	// 获取 component_verify_ticket

	cat := componentaccesstoken()
	if cat == "error" {
		cvt := componentverifyticket()
		if cvt == "error" {
			fmt.Println("获取不到缓存中的 component_verify_ticket")
			c.JSON(200, gin.H{"errCode": "error", "info": "not find component_verify_ticket error "})
			return
		}
		resComponentAccessToken, err := platforms.ComponentAccessToken(cfg.ComponentAppID, cfg.ComponentAppSecret, cvt)
		if err != nil {
			fmt.Println(err)
			c.JSON(200, gin.H{"errCode": "error", "info": "resComponentAccessToken.ErrMsg " + fmt.Sprintf("%s", err)})
			return
		}
		cat = resComponentAccessToken.ComponentAccessToken
		return
	}
	accessToken, err := coauth.GetAccessToken(AppID, Code, cfg.ComponentAppID, cat)
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{"errCode": "error", "info": "accessToken.ErrMsg " + fmt.Sprintf("%s", err)})
		return
	}
	userInfo, err := coauth.GetUserInfo(accessToken.AccessToken, accessToken.OpenID)
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{"errCode": "error", "info": "userInfo.ErrMsg " + fmt.Sprintf("%s", err)})
		return
	}
	db := new(config.DB)
	engine, err := db.NewEngine()
	if err != nil {
		fmt.Println(err)
		x := fmt.Sprintf("%s", err)
		c.JSON(200, gin.H{"errCode": "error", "info": "new engine error " + x})
		return
	}
	user := new(RCUser)
	err = engine.Sync2(user)
	if err != nil {
		x := fmt.Sprintf("%s", err)
		c.JSON(200, gin.H{"errCode": "error", "info": "engine Sync2 error " + x})
		return
	}
	// 防止重复投票 需要先查询
	bool, err := engine.Where("openid=?", userInfo.OpenID).Get(user)
	if err != nil {
		x := fmt.Sprintf("%s", err)
		c.JSON(200, gin.H{"errCode": "error", "info": "engine.Where error " + x})
		return
	} else if bool {
		// 存在直接跳转
		rURL := []byte("<script>window.location.href='https://gowc.iuu.pub/representcat/page?rcu_uuid=" + user.UUID + "&rcc_uuid=" + rccUUID + "';</script>")
		c.Writer.Write(rURL)
		return
	}
	user.UUID = uuid.New().String()
	user.VoteUUID = "0"
	user.ConfigUUID = rccUUID
	user.OpenID = userInfo.OpenID
	user.NickName = userInfo.NickName
	user.Sex = userInfo.Sex
	user.Avatar = userInfo.HeadImgURL
	_, err = engine.Insert(user)
	if err != nil {
		x := fmt.Sprintf("%s", err)
		c.JSON(200, gin.H{"errCode": "error", "info": "Insert error " + x})
		return
	}
	rURL := []byte("<script>window.location.href='https://gowc.iuu.pub/representcat/page?rcu_uuid=" + user.UUID + "';</script>")
	c.Writer.Write(rURL)
}

var ca = new(config.Cache)

// componentverifyticket
func componentverifyticket() string {
	str, err := ca.Get("ComponentVerifyTicket")
	if err != nil {
		return "error"
	}
	return str
}

// componentaccesstoken
func componentaccesstoken() string {
	str, err := ca.Get("ComponentAccessToken")
	if err != nil {
		return "error"
	}
	return str
}
