package autoconfig

import "go-gateway/pkg/httpgateway"

// 标准的配置文件示例
var GateWayC = `{
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
					"uri": "https://www.baidu.com",
					"predicates": [{
						"name": "Weight",
						"args": {
							"group": "s1",
							"value": "30"
						}
					},{
						"name": "Path",
						"args": {
							"value": "/index.php"
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
					"uri": "https://www.baidu.com",
					"predicates": [{
						"name": "Weight",
						"args": {
							"group": "s1",
							"value": "20"
						}
					},{
						"name": "Path",
						"args": {
							"value": "/index.php"
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

// 基础配置转换加载到内存
var (
	GateWaycfgInstance = &GateWayCfg{}
	PredicateConfig map[routeId][]httpgateway.PredicateConfig
	BizRoutes map[routeId]BizRoute
)

type GateWayCfg struct {

}

// 初始化
func (cfg *GateWayCfg) Init()  {
	gateWayConfig := ConfigParse(GateWayC)
	var parseErr error
	PredicateConfig, BizRoutes, parseErr = PredicatesParse(gateWayConfig)
	if parseErr != nil {
		panic(parseErr.Error())
	}
}

// 监听配置变化W
func (cfg *GateWayCfg) Watch()  {

}


