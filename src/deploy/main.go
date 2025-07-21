package main

import (
	"bufio"
	"command/src/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"gopkg.in/yaml.v2"
)

type configType struct {
	Port       int
	PathName   string
	LogsUrl    string
	CommandUrl string `yaml:"commandUrl"`
	DeployKey  string `yaml:"deployKey"`
}

// 默认配置
var Config = configType{
	Port:       3738,
	PathName:   "/deploy",
	LogsUrl:    "logs/",
	CommandUrl: ".",
	DeployKey:  "",
}
var template = `
port: 3738          # 服务端口
pathName: "/deploy" # 请求路径
logsUrl: "logs/"    # 日志路径
commandUrl: "."     # 执行命令目录
deployKey: ""       # 部署秘钥
`

func init() {
	configFile := "./deploy.yml"
	data, err := os.ReadFile(configFile)
	if err != nil {
		os.Create(configFile)
		os.WriteFile(configFile, []byte(template), 0777)
		data, _ = os.ReadFile(configFile)
	}

	if err := yaml.Unmarshal([]byte(data), &Config); err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc(Config.PathName, handleFuncFunc)
	port := ":" + strconv.Itoa(Config.Port)
	fmt.Println("http://localhost" + port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println(err)
	}
}

func handleFuncFunc(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if msg := recover(); msg != nil {
			message := fmt.Sprintf("%v", msg)
			errorLog(r, message)
			http.Error(w, message, http.StatusInternalServerError)
		}
	}()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	if r.Method != "POST" {
		panic("Method not allowed")
	}

	// 校验部署秘钥
	if Config.DeployKey != r.Header.Get("Deploy-Key") {
		panic("Deploy key error")
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		panic("Streaming unsupported")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	bodyStruct := map[string]interface{}{}
	if err := json.Unmarshal(body, &bodyStruct); err != nil {
		panic(err)
	}

	type Params struct {
		Command string `json:"command" binding:"required"`
	}
	params := Params{}
	if err := json.Unmarshal(body, &params); err != nil {
		panic(err)
	}

	entries, err := os.ReadDir(Config.CommandUrl)
	if err != nil {
		panic(err)
	}

	// 找到对应的执行文件
	command := ""
	for _, entry := range entries {
		if params.Command+".sh" == entry.Name() {
			command = "./" + entry.Name()
			break
		}
	}
	if command == "" {
		panic("Exec file not found")
	}

	result := ""
	stream := func(str string) {
		result += str + "\n"
		fmt.Fprintf(w, "data: %s\n\n", str)
		flusher.Flush()
	}
	stream("Exec command: " + command)

	cmd := exec.Command("sh", command)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	// 实时读取 stdout
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			stream(scanner.Text())
		}
	}()

	// 实时读取 stderr
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			stream(scanner.Text())
		}
	}()

	// 等待命令执行完成
	if err := cmd.Wait(); err != nil {
		panic(err)
	} else {
		writeLog(result)
		stream("Deploy successfully.")
	}
}

// 错误日志记录
func errorLog(r *http.Request, errStr string) {
	headerString := "header:\n"
	for name, headers := range r.Header {
		for _, h := range headers {
			headerString += "\t" + name + ": " + h + "\n"
		}
	}
	writeLog("url: " + r.URL.Path + "\n" + headerString + "error: " + errStr)
}

// 写入日志
func writeLog(msg string) {
	_, err := os.Stat(Config.LogsUrl)
	if err != nil {
		os.MkdirAll(Config.LogsUrl, os.ModePerm)
	}
	filename := Config.LogsUrl + utils.Date.DateFormater(time.Now(), "YYYY-MM-DD") + ".log"
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	file.WriteString(utils.Date.DateFormater(time.Now(), "YYYY-MM-DD hh:mm:ss") + "\n" + msg + "\n\n")
}
