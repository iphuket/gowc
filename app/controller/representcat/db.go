package representcat

import "time"

// RcConfig 活动配置
type RcConfig struct {
	UUID      string    `xorm:"varchar(36) pk notnull unique 'uuid'"` // uuid
	Name      string    `xorm:"varchar(512) 'name'"`                  // 活动名称
	StartTime time.Time `xorm:"start_time"`                           // 开始时间
	EndTime   time.Time `xorm:"end_time"`                             // 结束时间
	CreatedAt int64     `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	DeletedAt time.Time `xorm:"deleted"`
	Version   int64     `xorm:"version"`
}

// RcVote 用户数据
type RcVote struct {
	UUID       string    `xorm:"varchar(36) pk notnull unique 'uuid'"` // uuid
	ConfigUUID string    `xorm:"varchar(36) 'config_uuid'"`            // 所属活动的uuid
	Name       string    `xorm:"varchar(64) 'name'"`                   // 选票名称
	Count      int64     `xorm:"varchar(32) 'count'"`                  // 投票计数
	CreatedAt  int64     `xorm:"created"`
	UpdatedAt  time.Time `xorm:"updated"`
	DeletedAt  time.Time `xorm:"deleted"`
	Version    int64     `xorm:"version"`
}

// RcUser 用户信息
type RcUser struct {
	UUID       string    `xorm:"varchar(36) pk notnull unique 'uuid'"` // uuid
	ConfigUUID string    `xorm:"varchar(36) 'config_uuid'"`            // 所属活动的uuid
	VoteUUID   string    `xorm:"varchar(36) 'vote_uuid'"`              // 所属活动的uuid
	OpenID     string    `xorm:"varchar(64) notnull unique 'openid'"`  // openid -wx
	NickName   string    `xorm:"varchar(64) 'nickname'"`               // 昵称 -wx
	Sex        int32     `xorm:"varchar(16) 'sex'"`                    // 性别 -wx
	Avatar     string    `xorm:"varchar(512) 'avatar'"`                // 头像 -wxurl
	CreatedAt  int64     `xorm:"created"`
	UpdatedAt  time.Time `xorm:"updated"`
	DeletedAt  time.Time `xorm:"deleted"`
	Version    int64     `xorm:"version"`
}
