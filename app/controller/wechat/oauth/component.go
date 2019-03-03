package oauth

import (
	"encoding/json"
	"fmt"

	"github.com/iphuket/wechat/oauth"
	"github.com/iphuket/wechat/util"
)

const (
	codeURL        = "https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=%s&scope=%s&state=%s&component_appid=%s#wechat_redirect"
	accessTokenURL = "https://api.weixin.qq.com/sns/oauth2/component/access_token?appid=%s&code=%s&grant_type=authorization_code&component_appid=%s&component_access_token=%s"
	userInfoURL    = "https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN"
)

// GetCodeURL ... 换取codeURL
func GetCodeURL(AppID, RedirectURL, ResponseType, scope, state, ComponentAppID string) string {
	url := fmt.Sprintf(codeURL, AppID, RedirectURL, ResponseType, scope, state, ComponentAppID)
	return url
}

// GetAccessToken ... 仅用以获取用户信息
func GetAccessToken(AppID, Code, ComponentAppID, ComponentAccessToken string) (result oauth.ResAccessToken, err error) {
	urlStr := fmt.Sprintf(accessTokenURL, AppID, Code, ComponentAppID, ComponentAccessToken)
	var response []byte
	response, err = util.HTTPGet(urlStr)
	if err != nil {
		return
	}

	err = json.Unmarshal(response, &result)
	if err != nil {
		return
	}
	if result.ErrCode != 0 {
		err = fmt.Errorf("GetUserInfo error : errcode=%v , errmsg=%v", result.ErrCode, result.ErrMsg)
		return
	}
	return
}

// GetUserInfo ... 获取用户信息
func GetUserInfo(accessToken, openID string) (result oauth.UserInfo, err error) {
	urlStr := fmt.Sprintf(userInfoURL, accessToken, openID)
	var response []byte
	response, err = util.HTTPGet(urlStr)
	if err != nil {
		return
	}

	err = json.Unmarshal(response, &result)
	if err != nil {
		return
	}
	if result.ErrCode != 0 {
		err = fmt.Errorf("GetUserInfo error : errcode=%v , errmsg=%v", result.ErrCode, result.ErrMsg)
		return
	}
	return
}
