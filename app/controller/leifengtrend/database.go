package leifengtrend

import "time"

// LfConfig 活动配置
type LfConfig struct {
	UUID      string    `xorm:"varchar(36) pk notnull unique 'uuid'"` // uuid
	Name      string    `xorm:"varchar(512) 'name'"`                  // 活动名称
	Tel       string    `xorm:"varchar(64) 'tel'"`                    // 活动描述
	StartTime time.Time `xorm:"start_time"`                           // 开始时间
	EndTime   time.Time `xorm:"end_time"`                             // 结束时间
	CreatedAt int64     `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	DeletedAt time.Time `xorm:"deleted"`
	Version   int64     `xorm:"version"`
}

// LfUser 用户信息
type LfUser struct {
	UUID       string    `xorm:"varchar(36) pk notnull unique 'uuid'"` // uuid
	ConfigUUID string    `xorm:"varchar(36) 'config_uuid'"`            // 所属活动的uuid
	OpenID     string    `xorm:"varchar(64) notnull unique 'openid'"`  // openid -wx
	NickName   string    `xorm:"varchar(64) 'nickname'"`               // 昵称 -wx
	Sex        int32     `xorm:"varchar(16) 'sex'"`                    // 性别 -wx
	Avatar     string    `xorm:"varchar(512) 'avatar'"`                // 头像 -wxurl
	CreatedAt  int64     `xorm:"created"`
	UpdatedAt  time.Time `xorm:"updated"`
	DeletedAt  time.Time `xorm:"deleted"`
	Version    int64     `xorm:"version"`
}

// LfData 用户数据
type LfData struct {
	UUID       string    `xorm:"varchar(36) pk notnull unique 'uuid'"` // uuid
	ConfigUUID string    `xorm:"varchar(36) 'config_uuid'"`            // 所属活动的uuid
	Task       int32     `xorm:"varchar(16) 'task'"`                   // 所选任务ID
	ImageData  string    `xorm:"char 'image_data'"`                    // 图像数据
	CreatedAt  int64     `xorm:"created"`
	UpdatedAt  time.Time `xorm:"updated"`
	DeletedAt  time.Time `xorm:"deleted"`
	Version    int64     `xorm:"version"`
}
