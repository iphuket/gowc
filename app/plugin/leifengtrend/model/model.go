package model

import "time"

// LeifengtrendConfig 活动配置
type LeifengtrendConfig struct {
	UUID      string    `xorm:"varchar(36) notnull unique 'uuid'"` // uuid
	Name      string    `xorm:"varchar(512) 'name'"`               // 活动名称
	Tel       string    `xorm:"varchar(64) 'tel'"`                 // 活动描述
	StartTime time.Time `xorm:"start_time"`                        // 开始时间
	EndTime   time.Time `xorm:"end_time"`                          // 结束时间
	CreatedAt int64     `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	DeletedAt time.Time `xorm:"deleted"`
	Version   int64     `xorm:"version"`
}

// LeifengtrendUser 用户数据
type LeifengtrendUser struct {
	UUID       string    `xorm:"varchar(36) notnull unique 'uuid'"` // uuid
	ConfigUUID string    `xorm:"varchar(36) 'config_uuid'"`         // 所属活动的uuid
	OpenID     string    `xorm:"varchar(64) 'openid'"`              // openid -wx
	NickName   string    `xorm:"varchar(64) 'nickname'"`            // 昵称 -wx
	Name       string    `xorm:"varchar(64) 'name'"`                // 姓名
	Tel        string    `xorm:"varchar(64) 'tel'"`                 // 电话
	Sex        string    `xorm:"varchar(16) 'sex'"`                 // 性别 -wx
	Avatar     string    `xorm:"varchar(512) 'avatar'"`             // 头像 -wxurl
	CreatedAt  int64     `xorm:"created"`
	UpdatedAt  time.Time `xorm:"updated"`
	DeletedAt  time.Time `xorm:"deleted"`
	Version    int64     `xorm:"version"`
}

// LeifengtrendTask 任务表
type LeifengtrendTask struct {
	UUID       string    `xorm:"varchar(36) notnull unique 'uuid'"` // uuid
	ConfigUUID string    `xorm:"varchar(36) 'config_uuid'"`         // 所属活动的uuid
	Name       string    `xorm:"varchar(64) 'name'"`                // 任务名字
	Desc       string    `xorm:"varchar(64) 'desc'"`                // 任务描述
	CreatedAt  int64     `xorm:"created"`
	UpdatedAt  time.Time `xorm:"updated"`
	DeletedAt  time.Time `xorm:"deleted"`
	Version    int64     `xorm:"version"`
}

// LeifengtrendUserTask 用户领取的任务表
type LeifengtrendUserTask struct {
	UUID       string    `xorm:"varchar(36) notnull unique 'uuid'"` // uuid
	ConfigUUID string    `xorm:"varchar(36) 'config_uuid'"`         // 所属活动的uuid
	TaskUUID   string    `xorm:"varchar(36) 'task_uuid'"`           // 任务UUID
	UserUUID   string    `xorm:"varchar(36) 'user_uuid'"`           // 任务UUID
	CreatedAt  int64     `xorm:"created"`
	UpdatedAt  time.Time `xorm:"updated"`
	DeletedAt  time.Time `xorm:"deleted"`
	Version    int64     `xorm:"version"`
}
