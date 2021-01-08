package predicatefactory

import "go-gateway/pkg/httpgateway"

// 对应的谓词方法
var SupportPredicateFunc = make(map[httpgateway.EnumPredicae]httpgateway.WarpPredicateFunc)

func init()  {
	SupportPredicateFunc[httpgateway.PATH] = PathRoute
	SupportPredicateFunc[httpgateway.WEIGHT] = WeightRoute
}
