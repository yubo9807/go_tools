// 服务代理
// 我会生成一个 proxy.yml 的配置文件，在里面可以修改代理地址
package main

import (
	"command/src/utils"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strconv"

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
var template = `https: false
# 生成证书
# openssl req -newkey rsa:2048 -nodes -keyout server.key -x509 -days 365 -out server.crt

port: 9000  # 启动端口
proxy:
  - prefix: "^/"
    target: "http://hpyyb.cn"
  - prefix: "^/api"
    target: "http://hicky.hpyyb.cn"

`

func init() {
	configFile := "./proxy.yml"
	data, err := os.ReadFile(configFile)
	if err != nil {
		os.Create(configFile)
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

	// 设置端口，占用后 ++
	port := ":" + strconv.Itoa(utils.Server.PortResult(config.Port))
	fmt.Println(utils.If(config.Https, "https", "http") + "://localhost" + port)

	if config.Https {
		// 启动 https 服务
		if err := http.ListenAndServeTLS(port, "server.crt", "server.key", proxy); err != nil {
			fmt.Println(err.Error())
		}
	} else {
		// 启动 http 服务
		if err := http.ListenAndServe(port, proxy); err != nil {
			fmt.Println(err.Error())
		}
	}
}

// 创建代理
func newMultipleHostsReverseProxy(targets []*url.URL) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		for i, val := range config.Proxy {
			reg, _ := regexp.Compile(val.Prefix)
			matched := reg.MatchString(req.RequestURI)
			if matched {
				target := targets[i]
				req.URL.Scheme = target.Scheme
				req.URL.Host = target.Host
				req.Host = target.Host
			} else {
				fmt.Println("匹配不到任何相关代理地址")
			}
		}
	}
	return &httputil.ReverseProxy{Director: director}
}
