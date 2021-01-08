package predicatefactory

import (
	"context"
	//"github.com/pkg/errors"
	. "go-gateway/pkg/httpgateway"
	"net/http"
	"time"
)

var BetweenRoute WarpPredicateFunc = func(req *http.Request, config PredicateConfig, ctx context.Context) PredicateFunc {
	return func(req *http.Request) bool {
		timeNow := time.Now()
		if config.DateTime.IsZero() || config.DateTime2.IsZero() {
			return false;
		}
		return timeNow.After(config.DateTime) && timeNow.Before(config.DateTime2)
	}
}
