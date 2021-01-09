package proxy

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"go-gateway/pkg/httpgateway"
	"go-gateway/pkg/httpgateway/autoconfig"
	"go-gateway/pkg/httpgateway/predicatefactory"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"strings"
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
		ctx := context.Background()
		ctx = context.WithValue(ctx, predicatefactory.WeightRandomKey, rand.New(rand.NewSource(time.Now().UnixNano())).Float64())
		hitsRoute := make([]autoconfig.BizRoute, 0)
		for key, predicateCfgs := range autoconfig.PredicateConfig {
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
				hitsRoute = append(hitsRoute, autoconfig.BizRoutes[key])
			}
		}
		httpgateway.GetLogger().Debug(fmt.Sprintf("hitsRoute:%#v", hitsRoute))
		hits := len(hitsRoute)
		if hits == 1{
			scheme, host, err := UriParse(hitsRoute[0].Uri)
			if err != nil {
				req.URL.Host = ""
			} else {
				req.URL.Scheme = scheme
				req.URL.Host = host
				req.Host = ""
			}
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

func UriParse(uri string) (scheme string, host string, err error) {
	if uri != "" {
		splictRs := strings.Split(uri, "://")
		if len(splictRs) == 2 {
			scheme = splictRs[0]
			host = splictRs[1]
		} else {
			err = errors.New("uri parse, uri was empty string")
		}
	} else {
		err = errors.New("uri parse, uri was empty string")
	}
	return
}


