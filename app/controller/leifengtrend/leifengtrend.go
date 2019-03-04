package leifengtrend

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/iphuket/gowc/config"
	"github.com/iphuket/wechat/coauth"
	"github.com/iphuket/wechat/platforms"
)

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
	RedirctURL = "http://gowc.iuu.pub/leifengtrend/oauth/wechat/call"
)

// LeifengTrend 控制器
type LeifengTrend struct {
}

// WeChatLoginURL 微信登陆地址
func (lt *LeifengTrend) WeChatLoginURL(c *gin.Context) {
	uuid := c.Request.FormValue("guuid")
	if len(uuid) < 1 {
		fmt.Println("not uuid")
		// 不存在UUID
		// c.JSON(200, gin.H{"errCode": "error", "info": "not find UUID error "})
		scope := "snsapi_userinfo"
		state := "state"
		ComponentAppID := cfg.ComponentAppID
		url := coauth.GetCodeURL(AppID, RedirctURL, scope, state, ComponentAppID)
		c.Redirect(302, url)
		return
	}
	c.Redirect(302, "https://6.u.mgd5.com/c/1x5z/gru2/index.html?guuid="+uuid)
}

// WeChatLoginCall 微信登陆回调
func (lt *LeifengTrend) WeChatLoginCall(c *gin.Context) {
	Code := c.Request.FormValue("code")
	if len(Code) < 1 {
		c.JSON(200, gin.H{"errCode": "error", "info": "not find code error "})
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
	// c.JSON(200, gin.H{"errCode": "success", "OpenID": userInfo.OpenID, "NickName": userInfo.Nickname, "Sex": userInfo.Sex, "HeadImgURL": userInfo.HeadImgURL})

	db := new(config.DB)
	engine, err := db.NewEngine()
	if err != nil {
		fmt.Println(err)
		x := fmt.Sprintf("%s", err)
		c.JSON(200, gin.H{"errCode": "error", "info": "new engine error " + x})
		return
	}
	user := new(LfUser)
	err = engine.Sync2(user)
	if err != nil {
		x := fmt.Sprintf("%s", err)
		c.JSON(200, gin.H{"errCode": "error", "info": "engine Sync2 error " + x})
		return
	}
	// 防止重复注册 需要先查询
	bool, err := engine.Where("uuid=?", userInfo.OpenID).Get(user)
	if err != nil {
		x := fmt.Sprintf("%s", err)
		c.JSON(200, gin.H{"errCode": "error", "info": "engine.Where error " + x})
		return
	} else if bool {
		// 存在直接跳转
		c.Redirect(302, "https://6.u.mgd5.com/c/1x5z/gru2/index.html?guuid="+user.UUID)
		return
	}
	// 不存在进行注册后跳转
	user.UUID = userInfo.OpenID
	user.ConfigUUID = "77a3460a-668c-4145-b20e-2e5dce7dae81"
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
	c.Redirect(302, "https://6.u.mgd5.com/c/1x5z/gru2/index.html?guuid="+user.UUID)
	// c.SetCookie("openid", userInfo.OpenID, 99999, "/", "gowc.iuu.pub", false, true)
}

// CheckState 查询
func (lt *LeifengTrend) CheckState(c *gin.Context) {
	db := new(config.DB)
	engine, err := db.NewEngine()
	if err != nil {
		x := fmt.Sprintf("%s", err)
		c.JSON(200, gin.H{"errCode": "error", "info": "new database engine error " + x})
		return
	}
	type ResBody struct {
		UUID  string `json:"uuid"`
		CUUID string `json:"cuuid"`
	}
	resBody := new(ResBody)
	// var resBody ResBody
	err = c.BindJSON(resBody)
	if err != nil {
		x := fmt.Sprintf("%s", err)
		c.JSON(200, gin.H{"errCode": "error", "info": "bind json error " + x})
		return
	}
	ld := new(LfData)
	// sync2 database struct
	Sync2 := engine.Sync2(ld)
	if Sync2 != nil {
		x := fmt.Sprintf("%s", err)
		c.JSON(200, gin.H{"errCode": "error", "info": "engine Sync2 error " + x})
		return
	}
	// 查询用户状态
	// c.Writer.WriteString(resBody.UUID)
	bool, err := engine.Where("uuid=?", resBody.UUID).Get(ld)
	if err != nil {
		x := fmt.Sprintf("%s", err)
		c.JSON(200, gin.H{"errCode": "error", "info": "engine.Where error " + x})
		return
	} else if bool {
		// 存在任务 继续检查是否已经完成
		bool, err = pathExists("./img/" + resBody.UUID + ".jpg")
		if err == nil && bool {
			// 已经完成
			fmt.Println(err)
			c.JSON(200, gin.H{"errCode": "success", "state": 2, "task": ld.Task, "image": "https://gowc.iuu.pub/leifengtrend/img/" + resBody.UUID + ".jpg"})
			return
		}
		// 没有完成完成
		fmt.Println(err)
		c.JSON(200, gin.H{"errCode": "success", "state": 1, "task": ld.Task})
		return

	}
	// 没有领取任务
	c.JSON(200, gin.H{"errCode": "success", "state": 0, "o": ld.UUID})

}

// RegisterTask 注册任务
func (lt *LeifengTrend) RegisterTask(c *gin.Context) {
	db := new(config.DB)
	engine, err := db.NewEngine()
	if err != nil {
		c.JSON(200, gin.H{"errCode": "error", "info": "new database engine error " + fmt.Sprintf("%s", err)})
		return
	}
	type ResBody struct {
		UUID  string `json:"uuid"`
		CUUID string `json:"cuuid"`
		Task  int32  `json:"task"`
	}
	resBody := new(ResBody)
	err = c.BindJSON(resBody)
	if err != nil {
		c.JSON(200, gin.H{"errCode": "error", "info": "bind json error " + fmt.Sprintf("%s", err)})
		return
	}
	ld := new(LfData)
	// sync2 database struct
	Sync2 := engine.Sync2(ld)
	if Sync2 != nil {
		c.JSON(200, gin.H{"errCode": "error", "info": "engine Sync2 error " + fmt.Sprintf("%s", err)})
		return
	}

	// lf  谨慎重复注册 lf_user_uuid
	bool, err := engine.Where("uuid=?", resBody.UUID).Get(ld)
	if err != nil {
		c.JSON(200, gin.H{"errCode": "error", "info": "engine.Where error " + fmt.Sprintf("%s", err)})
		return
	} else if bool {
		c.JSON(200, gin.H{"errCode": "error", "info": "请不要重复领取任务"})
		return
	}
	ld.UUID = resBody.UUID
	ld.ConfigUUID = "77a3460a-668c-4145-b20e-2e5dce7dae81"
	ld.Task = resBody.Task
	_, err = engine.Insert(ld)
	if err != nil {
		x := fmt.Sprintf("%s", err)
		c.JSON(200, gin.H{"errCode": "error", "info": "Insert error " + x})
		return
	}
	c.JSON(200, gin.H{"errCode": "success"})
}

// UpdataImage 上传数据接口
func (lt *LeifengTrend) UpdataImage(c *gin.Context) {
	db := new(config.DB)
	engine, err := db.NewEngine()
	if err != nil {
		c.JSON(200, gin.H{"errCode": "error", "info": "new database engine error " + fmt.Sprintf("%s", err)})
		return
	}
	type ResBody struct {
		UUID      string `json:"uuid"`
		CUUID     string `json:"cuuid"`
		ImageData string `json:"image_data"`
		Task      int32  `json:"task"`
	}
	resBody := new(ResBody)
	err = c.BindJSON(resBody)
	if err != nil {
		c.JSON(200, gin.H{"errCode": "error", "info": "bind json error " + fmt.Sprintf("%s", err)})
		return
	}
	ld := new(LfData)
	// sync2 database struct
	Sync2 := engine.Sync2(ld)
	if Sync2 != nil {
		c.JSON(200, gin.H{"errCode": "error", "info": "engine Sync2 error " + fmt.Sprintf("%s", err)})
		return
	}
	// ld.Task = resBody.Task
	// ld.ImageData = "resBody.ImageData"
	/*
		_, err = engine.Where("lf_user_uuid", resBody.UUID).And("config_uuid", "77a3460a-668c-4145-b20e-2e5dce7dae81").Update(ld)
		if err != nil {
			x := fmt.Sprintf("%s", err)
			c.JSON(200, gin.H{"errCode": "error", "info": "Insert error " + x})
			return
		}
	*/
	ddd, err := base64.StdEncoding.DecodeString(resBody.ImageData[23:]) //成图片文件并把文件写入到buffer
	if err != nil {
		c.JSON(200, gin.H{"errCode": "error", "info": "base64 error " + fmt.Sprintf("%s", err)})
		return
	}
	err = ioutil.WriteFile("./img/"+resBody.UUID+".jpg", ddd, 0666)
	if err != nil {
		c.JSON(200, gin.H{"errCode": "error", "info": "wi error " + fmt.Sprintf("%s", err)})
		return
	}
	c.JSON(200, gin.H{"errCode": "success"})
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

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
