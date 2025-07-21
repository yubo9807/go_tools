package main

import (
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
	CommandUrl: "",
	DeployKey:  "",
}
var template = `
port: 3738          # 服务端口
pathName: "/deploy" # 请求路径
logsUrl: "logs/"    # 日志路径
commandUrl: ""      # 执行命令目录
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
			fmt.Println(message)
			errorLog(r, message)
			m := map[string]interface{}{"code": 500, "message": message}
			data, _ := json.Marshal(m)
			w.Write([]byte(data))
		}
	}()

	if r.Method != "POST" {
		panic("Method not allowed")
	}

	// 校验部署秘钥
	deployKey := r.Header.Get("Deploy-Key")
	if deployKey != Config.DeployKey {
		panic("DeployKey error")
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
	commandFile := ""
	for _, entry := range entries {
		if params.Command+".sh" == entry.Name() {
			commandFile = "./" + entry.Name()
			break
		}
	}

	if commandFile == "" {
		panic("command not found")
	}

	cmd := exec.Command("sh", commandFile)
	buf, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	cmdResult := string(buf)
	writeLog("cmdResult: " + cmdResult)
	m := map[string]interface{}{"code": 200, "message": cmdResult}
	data, _ := json.Marshal(m)
	w.Write(data)
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
