package controller

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"github.com/iphuket/gowc/app/plugin/leifengtrend/db"
	"github.com/iphuket/gowc/app/plugin/leifengtrend/model"
)

// Register 活动注册
func Register(c *gin.Context) {
	config := new(model.LeifengtrendConfig)
	fmt.Println(config)

	user := new(model.LeifengtrendUser)
	fmt.Println(user)

	engine, err := db.NewEngine()
	if err != nil {
		x := fmt.Sprintf("%s", err)
		c.JSON(200, gin.H{"errCode": "error", "info": "new engine error " + x})
		return
	}

	err = engine.Sync2(user)
	if err != nil {
		x := fmt.Sprintf("%s", err)
		c.JSON(200, gin.H{"errCode": "error", "info": "engine Sync2 error " + x})
		return
	}
	// 注册信息
	type Body struct {
		ConfigUUID string `json:"config_uuid"`
		NickName   string `json:"nickname"`
		OpenID     string `json:"openid"`
		Name       string `json:"name"`
		Sex        string `json:"sex"`
		Avatar     string `json:"avatar"`
	}
	rBody := new(Body)

	err = c.BindJSON(rBody)
	if err != nil {
		x := fmt.Sprintf("%s", err)
		c.JSON(200, gin.H{"errCode": "error", "info": "bind json error " + x})
		return
	}
	user.ConfigUUID = c.Param("config_uuid")
	user.UUID = uuid.New().String()
	_, err = engine.Insert(user)
	if err != nil {
		x := fmt.Sprintf("%s", err)
		c.JSON(200, gin.H{"errCode": "error", "info": "username or email already exists " + x})
		return
	}
	c.JSON(200, gin.H{"errCode": "success", "data": user.UUID})
}

// GetCaseInfo 获取活动信息... config_uuid
func GetCaseInfo(c *gin.Context) {
	config := new(model.LeifengtrendConfig)
	fmt.Println(config)
	user := new(model.LeifengtrendUser)
	fmt.Println(user)
}

// GetUserInfo 获取用户信息... openid + config_uuid
func GetUserInfo(c *gin.Context) {
	config := new(model.LeifengtrendConfig)
	fmt.Println(config)
	user := new(model.LeifengtrendUser)
	fmt.Println(user)
}
