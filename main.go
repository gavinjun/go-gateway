package main

import (
	"context"
	"fmt"
	"go-gateway/pkg/httpgateway/autoconfig"
	"go-gateway/pkg/httpgateway/predicatefactory"
	"go-gateway/pkg/httpgateway/proxy"
	"math/rand"
	"net/http"
	"time"
)


func sayHello(w http.ResponseWriter, r *http.Request) {
	// 配置
	gateWayConfig := autoconfig.ConfigParse()
	predicateConfig, bizRoutes := autoconfig.PredicatesParse(gateWayConfig)
	ctx := context.Background()
	ctx = context.WithValue(ctx, predicatefactory.WeightRandomKey, rand.New(rand.NewSource(time.Now().UnixNano())).Float64())
	hitsRoute := make([]autoconfig.BizRoute, 0)
	for key, predicateCfgs := range predicateConfig {
		var pass = true
		for _, predicateCfg := range predicateCfgs {
			if fn, ok := predicatefactory.SupportPredicateFunc[predicateCfg.PredicateType]; ok {
				rs := fn(r, predicateCfg, ctx)(r)
				if !rs {
					pass = false
					break
				}
			}
		}
		if pass {
			hitsRoute = append(hitsRoute, bizRoutes[key])
		}
	}
	fmt.Println(hitsRoute)
	proxy.DefaultProxy.ServeHTTP(w, r)
}
/**

 */

func main() {
	var index http.HandlerFunc = sayHello
	err := http.ListenAndServe("127.0.0.1:9000", index)
	if err != nil {
		fmt.Printf("http.ListenAndServe()函数执行错误,错误为:%v\n", err)
		return
	}
}