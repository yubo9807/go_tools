package main

import (
	"command/src/utils"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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
	Port     int
	PathName string
	LogsUrl  string
	Commands []commandType
}
type commandType struct {
	Repository string
	Branch     string
	Event      string
	Sh         string
}

// 默认配置
var Config = configType{
	Port:     3738,
	PathName: "/deploy",
	LogsUrl:  "logs/",
	Commands: []commandType{},
}
var template = `
port: 3738          # 服务端口
pathName: "/deploy" # 请求路径
logsUrl: "logs/"    # 日志路径
commands:
  -
    repository: "practical"
    branch: "main"
		event: "push"
		sh: "practical.sh"
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
	headerString := ""
	for name, headers := range r.Header {
		for _, h := range headers {
			headerString += name + ": " + h + "\n"
		}
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	defer r.Body.Close()

	if r.Method != "POST" {
		w.Write([]byte("Method not allowed"))
		return
	}
	writeLog(headerString + "body: " + string(body))

	bodyStruct := map[string]interface{}{}
	if err := json.Unmarshal(body, &bodyStruct); err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	repository := bodyStruct["repository"].(map[string]interface{})
	query := utils.Find(Config.Commands, func(v commandType, i int) bool {
		return v.Repository == repository["name"]
	})
	if query.Repository == "" {
		w.Write([]byte("Not found config"))
		return
	}
	if repository["default_branch"] != query.Branch || r.Header.Get("X-GitHub-Event") != query.Event {
		w.Write([]byte("Error branch or event"))
		return
	}

	bl := verifySignature(r.Header.Get("X-Hub-Signature-256"), string(body))

	go func() {
		writeLog("verifySignature: " + strconv.FormatBool(bl))
		result := execCmd(query.Sh)
		writeLog("cmd: " + result)
	}()
	w.Write([]byte("ok"))
}

func verifySignature(payload string, signature string) bool {
	// 创建一个基于SHA-256的新HMAC实例
	mac := hmac.New(sha256.New, []byte("j1odi"))
	mac.Write([]byte(payload))
	expectedMAC := mac.Sum(nil)

	// 解码GitHub发送的签名（不包含' sha256='前缀）
	signature = signature[len("sha256="):]

	// 使用恒定时间比较防止时序攻击
	b, _ := hex.DecodeString(signature)
	return hmac.Equal(expectedMAC, b)
}

// 执行脚本
func execCmd(sh string) string {
	cmd := exec.Command("sh", sh)
	buf, err := cmd.Output()
	if err != nil {
		return err.Error()
	}
	return string(buf)
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
