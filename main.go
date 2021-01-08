package main

import (
	"context"
	"go-gateway/pkg/httpgateway/autoconfig"
	_ "go-gateway/pkg/httpgateway/autoconfig"
	"go-gateway/pkg/httpgateway/predicatefactory"
	"math/rand"
	"net/http"
	"time"
)


func sayHello(w http.ResponseWriter, r *http.Request) {
	//proxy.ServeHTTP(w, r)
}
/**
var index http.HandlerFunc = sayHello
	err := http.ListenAndServe("127.0.0.1:9000", index)
	if err != nil {
		fmt.Printf("http.ListenAndServe()函数执行错误,错误为:%v\n", err)
		return
	}
 */

func main() {
	// 配置
	gateWayConfig := autoconfig.ConfigParse()

	predicateConfig, _ := autoconfig.PredicatesParse(gateWayConfig)

	ctx := context.Background()
	ctx = context.WithValue(ctx, predicatefactory.WeightRandomKey, rand.New(rand.NewSource(time.Now().UnixNano())).Float64())

	for key, predicateCfgs := range predicateConfig {
		if key == "weight" {
			for _, cfg := range predicateCfgs {
				predicatefactory.WeightRoute(nil, cfg, ctx)(nil)
			}
		}
	}

	//config := httpgateway.PredicateConfig{
	//}
	//config.Id = "routeTest"
	//config.DateTime, _ = time.ParseInLocation("2006-01-02 15:04:05", "2021-01-06 16:20:00", time.Local)
	//config.WeightGroup = "group3"
	//config.WeightGroupIndex = 0
	//config.WeightRanges = []float64{0, 0.3, 0.9, 1.0}
	//
	//fmt.Println(predicatefactory.WeightRoute(nil, config, ctx)(nil))
}