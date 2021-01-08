package main

import (
	"fmt"
	"go-gateway/pkg/httpgateway/proxy"
	"net/http"
)


func sayHello(w http.ResponseWriter, r *http.Request) {
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