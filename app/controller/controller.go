package controller

import (
	"github.com/iphuket/gowc/app/controller/leifengtrend"
	"github.com/iphuket/gowc/app/controller/representcat"
	"github.com/iphuket/gowc/app/controller/wechat"
)

/*
// Controller ...
type Controller struct {
	*wechat.WeChat
}
*/

// WeChat ... Controller ...
type WeChat struct {
	*wechat.WeChat // 微信控制器

}

// LeifengTrend ... Controller ...
type LeifengTrend struct {
	*leifengtrend.LeifengTrend
}

// RepresentCat ... Controller ...
type RepresentCat struct {
	*representcat.RepresentCat
}
