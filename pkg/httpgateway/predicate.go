package httpgateway

import (
	"context"
	"net/http"
	"time"
)

// 谓词
type Predicate interface {
	Test(req *http.Request) bool
}

/**
	比较特殊的权重路由分组依赖参数 为WeightGroup，WeightRanges

    比如 Id 为 weightTest的路由有如下的配置
    group3,30

    比如 Id 为 weightTest2的路由有如下的配置
    group3, 60

    默认配置 100-30-60
    group3,10

	WeightGroup = group3
	WeightGroupIndex = 0
    WeightRanges = [0, 0.3, 0.9, 1.0]
 */

type PredicateConfig struct {
	Id string // 路由id
	DateTime time.Time // date1 日期 进行大于 小于判断
	DateTime2 time.Time // date2 日期 between 时候才需要 表示2个日期之前
	WeightGroup string // 权重分组分组名
	WeightGroupIndex int // 在同组内的index
	WeightRanges []float64 //权重分组区域
}

type PredicateFunc func(req *http.Request) bool

func (fn PredicateFunc) Test(req *http.Request) bool {
	return fn(req)
}

// 谓词方法
type WarpPredicateFunc func(req *http.Request, config PredicateConfig, ctx context.Context) PredicateFunc

// 谓词判断后执行 一般做些日志操作
type PredicateAfterFunc func(predicateTestRs bool)


