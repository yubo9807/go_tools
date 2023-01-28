// 在当前目录启动一个静态资源服务器
// 此服务器对当前目录下的前端框架打包后的文件做了处理
// 框架打包后的文件路由都由 js 控制
// 所以，找不到相应的文件，则会去找 index.html，依然交给 js 处理
// 另外，对前端框架多页面应用可进行配置，static.yml
package main

import (
	"command/src/utils"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"gopkg.in/yaml.v2"
)

type configType struct {
	Port   int
	Public string
	Routes map[string]string
}

// 默认配置
var config = configType{
	Port:   8000,
	Public: "./",
}
var template = `
port: 8000    # 服务端口
public: './'  # 静态资源地址
routes:       # 自定义路由，针对多页面应用
  "/admin/": "/admin.html"  # 这里的地址会拼接 public
  "/www/": "/www.html"
`

func init() {
	configFile := "./static.yml"
	data, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Println("多页面应用请新建 static.yml 进行配置：")
		fmt.Println(template)
	}

	if err := yaml.Unmarshal([]byte(data), &config); err != nil {
		panic(err)
	}
}

// 启动服务
func main() {
	// 找不到文件夹则创建一个
	_, err := os.Open(config.Public)
	if err != nil {
		os.Mkdir(config.Public, 0666)
	}

	// 静态资源服务
	http.Handle("/", http.FileServer(http.Dir(config.Public)))

	files, err := os.ReadDir(config.Public)
	if err != nil {
		fmt.Println(err)
	}

	collect := make([]string, 0)

	// 给配置文件中的路由注册访问规则
	for key := range config.Routes {
		val, ok := config.Routes[key]
		if ok {
			collect = append(collect, key)
			http.HandleFunc(key, func(w http.ResponseWriter, r *http.Request) {
				http.ServeFile(w, r, config.Public+val)
			})
		}
	}

	// 给目录下的每个文件夹都重置下文件访问规则
	for _, file := range files {
		path := "/" + file.Name()
		isRegister := utils.Slice.Includes(collect, path+"/")
		if !isRegister && file.IsDir() {
			http.HandleFunc(path+"/", handler(path))
		}
	}

	// 追加端口
	port := ":" + strconv.Itoa(utils.Server.PortResult(config.Port))
	fmt.Println("http://localhost" + port)
	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Println(err.Error())
	}
}

// 处理前端包
// 如果没找到指定的文件会去找目录下的 index.html
func handler(pattern string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		_, err := os.Open(config.Public + path)
		if err != nil { // 文件没找到
			http.ServeFile(w, r, config.Public+pattern+"/index.html")
			return
		}

		http.ServeFile(w, r, config.Public+path)
	}
}
