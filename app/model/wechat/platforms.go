package wechat

import "time"

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
