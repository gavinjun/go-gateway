package proxy

import (
	"net/http"
	"net/http/httputil"
)

var DefaultProxy *httputil.ReverseProxy

func init() {
	DefaultProxy = NewDefaultReverseProxy()
}

func NewDefaultReverseProxy() *httputil.ReverseProxy {
	director := func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = "www.baidu.com"
		req.Host = ""
		req.URL.Path = "/api/panel/getDeleteVehicleList/"
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}
	return &httputil.ReverseProxy{Director: director}
}


