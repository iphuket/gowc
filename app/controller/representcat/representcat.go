package representcat

import (
	"fmt"
	"net/http"
	"net/url"

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

	c.HTML(http.StatusOK, "representcat.html", map[string]interface{}{"test": "test"})
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
	db := new(config.DB)
	engine, err := db.NewEngine()
	if err != nil {
		fmt.Println(err)
		x := fmt.Sprintf("%s", err)
		c.JSON(200, gin.H{"errCode": "error", "info": "new engine error " + x})
		return
	}
	rcc := new(RcConfig)
	err = engine.Sync2(rcc)
	if err != nil {
		x := fmt.Sprintf("%s", err)
		c.JSON(200, gin.H{"errCode": "error", "info": "engine Sync2 error " + x})
		return
	}
	// 获取配置信息
	bool, err := engine.Where("uuid = ?", rccUUID).Get(rcc)
	if err != nil {
		x := fmt.Sprintf("%s", err)
		c.JSON(200, gin.H{"errCode": "error", "info": "engine.Where error " + x})
		return
	} else if !bool {
		c.JSON(200, gin.H{"errCode": "error", "info": "not find config_uuid"})
		return
	}

	state := rccUUID
	ComponentAppID := cfg.ComponentAppID
	url := coauth.GetCodeURL(AppID, url.QueryEscape(RedirctURL), scope, state, ComponentAppID)
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
	user := new(RcUser)
	err = engine.Sync2(user)
	if err != nil {
		x := fmt.Sprintf("%s", err)
		c.JSON(200, gin.H{"errCode": "error", "info": "engine Sync2 error " + x})
		return
	}
	// 防止重复投票 需要先查询
	bool, err := engine.Where("openid = ? AND config_uuid = ?", userInfo.OpenID, rccUUID).Get(user)
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
	rURL := []byte("<script>window.location.href='https://gowc.iuu.pub/representcat/page?rcu_uuid=" + user.UUID + "&rcc_uuid=" + rccUUID + "';</script>")
	c.Writer.Write(rURL)
}

// RepresentCatConfig 配置
func (rc *RepresentCat) RepresentCatConfig(c *gin.Context) {
	do := c.Request.FormValue("do")
	db := new(config.DB)
	engine, err := db.NewEngine()
	if err != nil {
		fmt.Println(err)
		x := fmt.Sprintf("%s", err)
		c.JSON(200, gin.H{"errCode": "error", "info": "new engine error " + x})
		return
	}
	if do == "create" {
		cfg := new(RcConfig)
		err = engine.Sync2(cfg)
		if err != nil {
			x := fmt.Sprintf("%s", err)
			c.JSON(200, gin.H{"errCode": "error", "info": "engine Sync2 error " + x})
			return
		}
		cfg.UUID = uuid.New().String()
		cfg.Name = c.Request.FormValue("cfg_name")
		_, err = engine.Insert(cfg)
		if err != nil {
			x := fmt.Sprintf("%s", err)
			c.JSON(200, gin.H{"errCode": "error", "info": "RepresentCatConfig Insert error " + x})
			return
		}
		c.JSON(200, gin.H{"errCode": "success", "info": "RepresentCatConfig Insert is ok", "uuid": cfg.UUID})
		return
	}
	if do == "update" {
		rccUUID := c.Request.FormValue("rcc_uuid")
		if len(rccUUID) < 1 {
			fmt.Println("not rcc_uuid")
		}
		rcv := new(RcVote)

		err = engine.Sync2(rcv)
		if err != nil {
			x := fmt.Sprintf("%s", err)
			c.JSON(200, gin.H{"errCode": "error", "info": "engine Sync2 error " + x})
			return
		}
		rcv.ConfigUUID = rccUUID
		rcv.Count = 0
		rcv.UUID = uuid.New().String()
		rcvName := c.Request.FormValue("rcv_name")
		rcv.Name = rcvName
		_, err = engine.Insert(rcv)
		if err != nil {
			x := fmt.Sprintf("%s", err)
			c.JSON(200, gin.H{"errCode": "error", "info": "RepresentCatVote Insert error " + x})
			return
		}
		c.JSON(200, gin.H{"errCode": "success", "info": "RepresentCatVote Insert is ok", "uuid": rcv.UUID})
		return
	}
	c.JSON(200, gin.H{"errCode": "error", "info": "Param 'do' not"})
}

// RepresentCatVote 投票
func (rc *RepresentCat) RepresentCatVote(c *gin.Context) {
	do, rcvUUID, rccUUID, rcuUUID := c.Request.FormValue("do"), c.Request.FormValue("rcv_uuid"), c.Request.FormValue("rcc_uuid"), c.Request.FormValue("rcu_uuid")
	if len(do) < 1 || len(rcvUUID) < 1 || len(rccUUID) < 1 || len(rcuUUID) < 1 {
		c.JSON(http.StatusOK, gin.H{"errCode": "error", "info": "not find do, rcv_uuid, rcc_uuid, rcu_uuid"})
		return
	}
	// 防止重复投票 需要先查询用户是否已投票
	db := new(config.DB)
	engine, err := db.NewEngine()
	if err != nil {
		fmt.Println(err)
		x := fmt.Sprintf("%s", err)
		c.JSON(200, gin.H{"errCode": "error", "info": "new engine error " + x})
		return
	}
	user := new(RcUser)
	err = engine.Sync2(user)
	if err != nil {
		x := fmt.Sprintf("%s", err)
		c.JSON(200, gin.H{"errCode": "error", "info": "engine Sync2 error " + x})
		return
	}
	// 防止重复投票 需要先查询
	bool, err := engine.Where("uuid = ? AND config_uuid = ?", rcuUUID, rccUUID).Get(user)
	if err != nil {
		x := fmt.Sprintf("%s", err)
		c.JSON(200, gin.H{"errCode": "error", "info": "engine.Where error " + x})
		return
	} else if bool {
		// 用户存在进入下一步 判断
		if len(user.VoteUUID) > 1 {
			// c.JSON(200, gin.H{"errCode": "error", "info": "您已经投票过了，请不要重复投票", "desk": "Don't repeat your vote"})
			data, err := voteCount(rccUUID)
			if err != nil {
				fmt.Println("res data error ", err)
				c.JSON(200, gin.H{"errCode": "error", "info": "res data error", "error": fmt.Sprintf("%s", err)})
				return
			}
			c.JSON(200, gin.H{"errCode": "error", "info": "您已经投票过了，请不要重复投票", "data": data})
			return
		}
		user.VoteUUID = rcvUUID
		affected, err := engine.Where("uuid = ? AND config_uuid = ?", rcuUUID, rccUUID).Update(user)
		if err != nil {
			x := fmt.Sprintf("%s", err)
			c.JSON(200, gin.H{"errCode": "error", "info": "engine.Where error " + x})
			return
		}
		if affected < 0 {
			c.JSON(200, gin.H{"errCode": "error", "info": "update error affected = 0"})
			return
		}
		voteP(c, rcvUUID, rccUUID)
		return
	}
	c.JSON(200, gin.H{"errCode": "error", "info": "unknown error or rcu_uuid not find"})
}
// RepresentCatGetData ...
func (rc *RepresentCat) RepresentCatGetData(c *gin.Context) {
	do, rcvUUID, rccUUID, rcuUUID := c.Request.FormValue("do"), c.Request.FormValue("rcv_uuid"), c.Request.FormValue("rcc_uuid"), c.Request.FormValue("rcu_uuid")
	if len(do) < 1 || len(rcvUUID) < 1 || len(rccUUID) < 1 || len(rcuUUID) < 1 {
		c.JSON(http.StatusOK, gin.H{"errCode": "error", "info": "not find do, rcv_uuid, rcc_uuid, rcu_uuid"})
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
	user := new(RcUser)
	err = engine.Sync2(user)
	if err != nil {
		x := fmt.Sprintf("%s", err)
		c.JSON(200, gin.H{"errCode": "error", "info": "engine Sync2 error " + x})
		return
	}
	bool, err := engine.Where("uuid = ? AND config_uuid = ?", rcuUUID, rccUUID).Get(user)
	if err != nil {
		x := fmt.Sprintf("%s", err)
		c.JSON(200, gin.H{"errCode": "error", "info": "engine.Where error " + x})
		return
	} else if bool {
		if len(user.VoteUUID) > 1 {
			data, err := voteCount(rccUUID)
			if err != nil {
				fmt.Println("res data error ", err)
				c.JSON(200, gin.H{"errCode": "error", "info": "res data error", "error": fmt.Sprintf("%s", err)})
				return
			}
			c.JSON(200, gin.H{"errCode": "success", "info": "您已经投票过了，请不要重复投票", "data": data})
			return
		}
		c.JSON(200, gin.H{"errCode": "error", "info": "还没有投票"})
		return
	}
	c.JSON(200, gin.H{"errCode": "error", "info": "出了点错误"})

}
func voteP(c *gin.Context, rcvUUID, rccUUID string) {
	rcv := new(RcVote)
	db := new(config.DB)
	engine, err := db.NewEngine()
	if err != nil {
		fmt.Println(err)
		x := fmt.Sprintf("%s", err)
		c.JSON(200, gin.H{"errCode": "error", "info": "new engine error " + x})
		return
	}
	err = engine.Sync2(rcv)
	if err != nil {
		x := fmt.Sprintf("%s", err)
		c.JSON(200, gin.H{"errCode": "error", "info": "engine Sync2 error " + x})
		return
	}
	bool, err := engine.Where("uuid = ? AND config_uuid = ?", rcvUUID, rccUUID).Get(rcv)
	if err != nil {
		x := fmt.Sprintf("%s", err)
		c.JSON(200, gin.H{"errCode": "error", "info": "engine.Where error " + x})
		return
	} else if bool {
		// 存在 即可以进行投票 更新票数
		rcv.Count = rcv.Count + 1
		affected, err := engine.Where("uuid = ? AND config_uuid = ?", rcvUUID, rccUUID).Update(rcv)
		if err != nil {
			x := fmt.Sprintf("%s", err)
			c.JSON(200, gin.H{"errCode": "error", "info": "engine.Where error " + x})
			return
		}
		if affected < 0 {
			c.JSON(200, gin.H{"errCode": "error", "info": "update error affected = 0"})
			return
		}
		data, err := voteCount(rccUUID)
		if err != nil {
			fmt.Println("res data error ", err)
			c.JSON(200, gin.H{"errCode": "error", "info": "res data error", "error": fmt.Sprintf("%s", err)})
			return
		}
		c.JSON(200, gin.H{"errCode": "success", "data": data})
		return
	}
	c.JSON(200, gin.H{"errCode": "error", "info": "not find rcv_uuid, rcc_uuid"})
	return
}
func voteCount(rccUUID string) ([]RcVote, error) {
	var rcv []RcVote
	db := new(config.DB)
	engine, err := db.NewEngine()
	if err != nil {
		return nil, err
	}
	err = engine.Where("config_uuid = ?", rccUUID).Find(&rcv)
	if err != nil {
		return nil, err
	}
	return rcv, nil
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
