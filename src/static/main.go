// 在当前目录启动一个静态资源服务器
// 此服务器对当前目录下的前端框架打包后的文件做了处理
// 框架打包后的文件路由都由 js 控制
// 所以，找不到相应的文件，则会去找 index.html，依然交给 js 处理
package main

import (
	"command/src/utils"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

type configType struct {
	Port   int
	Public string
}

// 默认配置
var config = configType{
	Port:   8000,
	Public: "./",
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

	// 给目录下的每个文件夹都重置下文件访问规则
	for _, file := range files {
		path := "/" + file.Name()
		if file.IsDir() {
			http.HandleFunc(path+"/", handler(path))
		}
	}

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
