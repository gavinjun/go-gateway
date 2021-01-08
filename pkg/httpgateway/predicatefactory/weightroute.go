package predicatefactory

import (
	"context"
	"fmt"
	. "go-gateway/pkg/httpgateway"
	"net/http"
)

var (
	WeightRandomKey = "WeightCalculatorRandom"
)

// 权重路由 x%几率通过
var WeightRoute WarpPredicateFunc = func(req *http.Request, config PredicateConfig, ctx context.Context) PredicateFunc {
	var randomKey float64 = -1.0
	if randomKeytmp,ok := ctx.Value(WeightRandomKey).(float64); ok {
		randomKey = randomKeytmp
	}
	var afterFunc PredicateAfterFunc = func(rs bool) {
		GetLogger().Debug(fmt.Sprintf("%s WeightCalculator predicate: %v , weightRanges:%v,index:%d,randomKey:%v", config.Id, rs, config.WeightRanges, config.WeightGroupIndex, randomKey))
	}
	return func(req *http.Request) bool {
		// 从ctx中获取随机数
		rs := false
		if randomKey >= 0 {
			// 获取到通过middleware设置进ctx的随机数 在当前区间内则表示通过路由断言
			if config.WeightRanges[config.WeightGroupIndex] <= randomKey && config.WeightRanges[config.WeightGroupIndex+1] > randomKey {
				rs = true
			}
		}
		afterFunc(rs)
		return rs
	}
}