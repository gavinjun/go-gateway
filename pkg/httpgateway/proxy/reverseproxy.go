package proxy

import (
	"context"
	"fmt"
	"go-gateway/pkg/httpgateway/autoconfig"
	"go-gateway/pkg/httpgateway/predicatefactory"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"time"
)

var DefaultProxy *httputil.ReverseProxy

func init() {
	DefaultProxy = NewDefaultReverseProxy()
}

func NewDefaultReverseProxy() *httputil.ReverseProxy {
	director := func(req *http.Request) {
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
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
					rs := fn(req, predicateCfg, ctx)(req)
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
		hits := len(hitsRoute)
		if hits == 1{
			req.URL.Scheme = "http"
			req.URL.Host = hitsRoute[0].Uri
			req.Host = ""
		} else if hits == 0 {
			// default
			req.URL.Host = ""
		} else {
			// error
			req.URL.Host = ""
		}

	}
	return &httputil.ReverseProxy{Director: director}
}


