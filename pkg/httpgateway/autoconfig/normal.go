package autoconfig

import (
	"encoding/json"
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

// 标准的配置文件示例
var demoConfigStr = `{
	"spring": {
		"application": {
			"name": "xxxx"
		},
		"cloud": {
			"consul": {
				"host": "localhost",
				"port": 8500,
				"discovery": {
					"enabled": true,
					"instance-id": "",
					"service-name": "xxxx",
					"prefer-ip-address": true
				}
			},
			"gateway": {
				"discovery": {
					"locator": {
						"enabled": true
					}
				},
				"routes": [{
					"id": "activity-route",
					"uri": "lb://activity",
					"predicates": [{
						"name": "Weight",
						"args": {
							"group": "s1",
							"value": "30"
						}
					}],
					"filters": [{
							"name": "AddRequestHeader",
							"args": {
								"name": "'foo'",
								"value": "'bar'"
							}
						},
						{
							"name": "RewritePath",
							"args": {
								"regexp": "'/' + serviceId + '/(?<remaining>.*)'",
								"replacement": "'/${remaining}'"
							}
						}
					]
				}, {
					"id": "activity-route-2",
					"uri": "lb://activity1",
					"predicates": [{
						"name": "Weight",
						"args": {
							"group": "s1",
							"value": "20"
						}
					}],
					"filters": [{
						"name": "AddRequestHeader",
						"args": {
							"name": "'foo'",
							"value": "'bar'"
						}
					}]
				}]
			}
		}
	}
}`

/**
  配置的形态
 */
type BizRoute struct {
	Id string
	Uri string
	Filters []Filter
}

func ConfigParse() *GateWayConfig {
	config := make(cmap)
	json.Unmarshal([]byte(demoConfigStr), &config)
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
type predicateName = string
type weightGroup = string
type weightGroupVal = int // group1 => routeId,value ps: weightGroup: 30

func PredicatesParse(gateWayConfig *GateWayConfig) (map[predicateName][]httpgateway.PredicateConfig, map[routeId]BizRoute) {
	var predicateConfigMap = make(map[predicateName][]httpgateway.PredicateConfig)
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
					Id: route.Id,
				}
				if string_util.StrCompareIgnoreLowerOrUpper(predicateName, "weight") {
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

					if _, ok := predicateConfigMap[predicateName]; !ok {
						predicateConfigMap[predicateName] = []httpgateway.PredicateConfig{}
					}
					predicateConfigMap[predicateName] = append(predicateConfigMap[predicateName], routeCfg)
				}
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
	if _, ok := predicateConfigMap["weight"]; ok {
		for k, vv := range predicateConfigMap["weight"] {
			predicateConfigMap["weight"][k].WeightRanges = weightGroupValuesMap[vv.WeightGroup]
		}
	}
	return predicateConfigMap, routeGlobMap
}

