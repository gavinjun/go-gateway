package predicatefactory

import (
	"context"
	//"github.com/pkg/errors"
	. "go-gateway/pkg/httpgateway"
	"net/http"
	"time"
)

var BeforeRoute WarpPredicateFunc = func(req *http.Request, config PredicateConfig, ctx context.Context) PredicateFunc {
	return func(req *http.Request) bool {
		timeNow := time.Now()
		if config.DateTime.IsZero() {
			return false;
		}
		return timeNow.Before(config.DateTime)
	}
}
