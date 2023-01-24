package main

import (
	"command/src/utils"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

type ProxyType struct {
	Prefix string
	Target string
}

type ConfigType struct {
	Https bool
	Port  int
	Proxy []ProxyType
}

var config ConfigType

func init() {
	configFile := "./proxy.yml"
	data, err := os.ReadFile(configFile)
	if err != nil {
		os.Create(configFile)
		template := `https: false
# 生成证书
# openssl req -newkey rsa:2048 -nodes -keyout server.key -x509 -days 365 -out server.crt

port: 9000  # 启动端口
proxy:
  - prefix: "/"
    target: "http://hpyyb.cn"
  - prefix: "/api"
    target: "http://hicky.hpyyb.cn"

`
		os.WriteFile(configFile, []byte(template), 0777)
		data, _ = os.ReadFile(configFile)
	}

	if err := yaml.Unmarshal([]byte(data), &config); err != nil {
		panic(err)
	}
}

func main() {
	slice := make([]*url.URL, 0)
	for _, val := range config.Proxy {
		addr, _ := url.Parse(val.Target)
		slice = append(slice, addr)
	}
	proxy := newMultipleHostsReverseProxy(slice)

	port := ":" + strconv.Itoa(utils.Server.PortResult(config.Port))
	agreement := "http"
	if config.Https {
		agreement = "https"
	}
	fmt.Println(agreement + "://localhost" + port)

	if config.Https {
		if err := http.ListenAndServeTLS(port, "server.crt", "server.key", proxy); err != nil {
			fmt.Println(err.Error())
		}
	} else {
		if err := http.ListenAndServe(port, proxy); err != nil {
			fmt.Println(err.Error())
		}
	}
}

func newMultipleHostsReverseProxy(targets []*url.URL) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		var target *url.URL
		for i, val := range config.Proxy {
			if strings.HasPrefix(req.RequestURI, val.Prefix) {
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
