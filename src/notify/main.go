// 一个定时任务，在指定的时间通知一些消息
package main

import (
	"fmt"
	"os"
	"time"

	"command/src/utils"

	"github.com/0xAX/notificator"
	"gopkg.in/yaml.v2"
)

var notify *notificator.Notificator
var record string

type NotifyType struct {
	Title string
	Text  string
}

type ConfigType struct {
	Times map[string]NotifyType
}

var config ConfigType

var template = `times:
  '09:00':
    title: '打卡打卡'
  '11:30':
    title: 'Go! Go! Go! 干饭干饭!'
    text: '手里活儿停一下

`

func init() {
	notify = notificator.New(notificator.Options{
		AppName: "时间提醒",
	})

	configFile := "./notify.yml"
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
	timer := time.NewTicker(time.Duration(time.Second) * 30)
	fmt.Println("启动成功：" + utils.Date.DateFormater(time.Now(), ""))
	for range timer.C {
		nowDate := utils.Date.DateFormater(time.Now(), "hh:mm")
		if record != nowDate {
			val, ok := config.Times[nowDate]
			if ok {
				notify.Push(val.Title, val.Text, "", notificator.UR_CRITICAL)
			}
			record = nowDate
		}
	}
}
