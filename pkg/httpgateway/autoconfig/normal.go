package autoconfig

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"go-gateway/pkg/httpgateway"
	"go-gateway/pkg/string_util"
	"strconv"
	"strings"
)
type cmap = map[string]interface{}

type Predicate struct {
	Name string
	Args map[string]string
}

type Filter struct {
	Name string
	Args map[string]string
}

type Discovery struct {
	Locator Locator
}
type Locator struct {
	Enabled bool
}
type Route struct {
	Id string
	Uri string
	Predicates []Predicate
	Filters []Filter
}
type GateWayConfig struct {
	Discovery Discovery
	Routes []Route
}

/**
  配置的形态
 */
type BizRoute struct {
	Id string
	Uri string
	Filters []Filter
}

func ConfigParse(ConfigStr string) *GateWayConfig {
	config := make(cmap)
	json.Unmarshal([]byte(ConfigStr), &config)
	//剥离掉前缀spring.cloud.gateway.discovery.locator
	prefixSlice := strings.Split("spring.cloud.gateway", ".")
	for i:=0; i < len(prefixSlice); i++ {
		keytemp := prefixSlice[i]
		if _,ok := config[keytemp]; ok {
			if vvv,ok := config[keytemp].(cmap); ok {
				config = vvv
			} else {
				break
			}
		} else {
			// 未找到配置
			break
		}
		continue
	}
	locator ,_ := json.Marshal(config)
	gateWayConfig := GateWayConfig{}
	json.Unmarshal(locator, &gateWayConfig)
	return &gateWayConfig
}

type routeId = string
type weightGroup = string
type weightGroupVal = int // group1 => routeId,value ps: weightGroup: 30

func PredicatesParse(gateWayConfig *GateWayConfig) (map[routeId][]httpgateway.PredicateConfig, map[routeId]BizRoute, error) {
	var predicateConfigMap = make(map[routeId][]httpgateway.PredicateConfig)
	// 全局的配置
	var routeGlobMap = make(map[routeId]BizRoute)
	// 权重路由暂存
	var weightGroupList = make(map[weightGroup][]weightGroupVal)
	// 根据名称判断谓词
	if len(gateWayConfig.Routes) > 0 {
		for _, route := range gateWayConfig.Routes {
			// 对filter进行拷贝
			filters := make([]Filter, len(route.Filters))
			copy(filters, route.Filters)
			routeGlobMap[route.Id] = BizRoute{
				Id:      route.Id,
				Uri:     route.Uri,
				Filters: filters,
			}
			for _, predicate := range route.Predicates {
				predicateName := strings.ToLower(predicate.Name)
				routeCfg := httpgateway.PredicateConfig{
					PredicateType:httpgateway.StringToPredicate(predicateName),
					Id: route.Id,
				}
				if string_util.StrCompareIgnoreLowerOrUpper(predicateName, httpgateway.PredicateToString(httpgateway.WEIGHT)) {
					// 权重路由
					var group string
					var value string
					if _, ok := predicate.Args["group"]; ok {
						group = predicate.Args["group"]
					}
					if _, ok := predicate.Args["value"]; ok {
						value = predicate.Args["value"]
					}
					if group != "" && value != "" {
						routeCfg.WeightGroup = group
						valueInt,_ := strconv.Atoi(value)
						weightGroupIndex := 0
						if _, ok := weightGroupList[group]; ok {
							// append
							weightGroupIndex = len(weightGroupList[group])
							weightGroupList[group] = append(weightGroupList[group], valueInt)
						} else {
							// init
							weightGroupList[group] = []weightGroupVal{valueInt}
						}
						routeCfg.WeightGroupIndex = weightGroupIndex
					}

					//if _, ok := predicateConfigMap[predicateName]; !ok {
					//	predicateConfigMap[predicateName] = []httpgateway.PredicateConfig{}
					//}
					//predicateConfigMap[predicateName] = append(predicateConfigMap[predicateName], routeCfg)
				}
				if string_util.StrCompareIgnoreLowerOrUpper(predicateName, httpgateway.PredicateToString(httpgateway.PATH)) {
					// 权重路由
					var value string
					if _, ok := predicate.Args["value"]; ok {
						value = predicate.Args["value"]
					}
					if value != "" {
						routeCfg.PathRegPattern = value
					}
				}
				if _, ok := predicateConfigMap[route.Id]; !ok {
					predicateConfigMap[route.Id] = []httpgateway.PredicateConfig{}
				}
				predicateConfigMap[route.Id] = append(predicateConfigMap[route.Id], routeCfg)
			}
		}
	}

	// 权重比较特殊需要后续处理
	var weightGroupValuesMap = make(map[weightGroup][]float64)
	for weightGroupName, weights := range weightGroupList {
		rangesLen := len(weights) + 2
		var ranges = make([]float64, rangesLen)
		var weightConvert = make([]float64, len(weights) + 1)
		maxWeight := 100
		var total = 0
		for _, value := range weights {
			total += value
		}
		if total > maxWeight {
			// 有异常
			return nil, nil, errors.New(fmt.Sprintf("%s total more than %d", weightGroupName, maxWeight))
		}
		defaultWeight := maxWeight - total
		for key, value := range weights {
			weightConvert[key] = float64(value) / float64(maxWeight)
		}
		weightConvert[len(weights)] = float64(defaultWeight) / float64(maxWeight)

		// 计算ranges
		var start float64
		ranges[0] = start
		for i:=0; i<len(weightConvert); i++{
			ranges[i+1] = ranges[i] + weightConvert[i]
		}
		weightGroupValuesMap[weightGroupName] = ranges
	}
	for k, v := range predicateConfigMap {
		for kk, vv := range v {
			if _, ok := weightGroupValuesMap[vv.WeightGroup]; ok {
				predicateConfigMap[k][kk].WeightRanges = weightGroupValuesMap[vv.WeightGroup]
			}
		}
	}
	return predicateConfigMap, routeGlobMap, nil
}

