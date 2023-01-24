package main

import (
	"command/src/utils"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
)

type proxyType struct {
	prefix string
	target string
}

var proxyList = []proxyType{}

func init() {
	proxyList = append(proxyList, proxyType{"/", "http://hpyyb.cn"})
	proxyList = append(proxyList, proxyType{"/api", "http://hicky.hpyyb.cn"})
}

func newMultipleHostsReverseProxy(targets []*url.URL) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		var target *url.URL
		for i, val := range proxyList {
			if strings.HasPrefix(req.RequestURI, val.prefix) {
				target = targets[i]
				continue
			} else {
				target = targets[0]
			}
		}
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host
	}
	return &httputil.ReverseProxy{Director: director}
}

func main() {
	slice := make([]*url.URL, 0)
	for _, val := range proxyList {
		addr, _ := url.Parse(val.target)
		slice = append(slice, addr)
	}
	proxy := newMultipleHostsReverseProxy(slice)

	port := ":" + strconv.Itoa(utils.Server.PortResult(9000))
	fmt.Println("http://localhost" + port)
	if err := http.ListenAndServe(port, proxy); err != nil {
		fmt.Println(err.Error())
	}
}
