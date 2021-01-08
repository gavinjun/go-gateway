package predicatefactory

import (
	"context"
	"fmt"
	. "go-gateway/pkg/httpgateway"
	"net/http"
	"regexp"
)

// path路由 x%几率通过
var PathRoute WarpPredicateFunc = func(req *http.Request, config PredicateConfig, ctx context.Context) PredicateFunc {


	var afterFunc PredicateAfterFunc = func(rs bool) {
		GetLogger().Debug(fmt.Sprintf("%s path predicate: %v , pathRegPattern:%s, path:%v", config.Id, rs, config.PathRegPattern, req.URL.Path))
	}

	return func(req *http.Request) bool {
		// 从ctx中获取随机数
		rs := false
		reqPath := req.URL.Path

		if reqPath != "" && config.PathRegPattern != "" {
			if compile, err := regexp.Compile(config.PathRegPattern); err == nil {
				rs = compile.MatchString(reqPath)
			}
		}
		afterFunc(rs)
		return rs
	}
}