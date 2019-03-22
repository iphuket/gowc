package representcat

import "time"

// RCConfig 活动配置
type RCConfig struct {
	UUID      string    `xorm:"varchar(36) pk notnull unique 'uuid'"` // uuid
	Name      string    `xorm:"varchar(512) 'name'"`                  // 活动名称
	StartTime time.Time `xorm:"start_time"`                           // 开始时间
	EndTime   time.Time `xorm:"end_time"`                             // 结束时间
	CreatedAt int64     `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	DeletedAt time.Time `xorm:"deleted"`
	Version   int64     `xorm:"version"`
}

// RCVote 用户数据
type RCVote struct {
	UUID       string    `xorm:"varchar(36) pk notnull unique 'uuid'"` // uuid
	ConfigUUID string    `xorm:"varchar(36) 'config_uuid'"`            // 所属活动的uuid
	Name       string    `xorm:"varchar(64) 'name'"`                   // 名称
	Count      int32     `xorm:"varchar(16) 'count'"`                  // 所选任务ID
	CreatedAt  int64     `xorm:"created"`
	UpdatedAt  time.Time `xorm:"updated"`
	DeletedAt  time.Time `xorm:"deleted"`
	Version    int64     `xorm:"version"`
}

// RCUser 用户信息
type RCUser struct {
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
