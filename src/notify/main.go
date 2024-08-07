// 一个定时任务，在指定的时间通知一些消息
// 通过配置 notify.yml 修改时间段提示
package main

import (
	"command/src/utils"
	"fmt"
	"os"
	"os/exec"
	"time"

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
	Interval  int
	SoundFile string `yaml:"soundFile"`
	Times     map[string]NotifyType
}

var config ConfigType
var template = `interval: 30  # 定时器间隔，秒
soundFile: '/System/Library/Sounds/Purr.aiff'  # 通知声音文件位置

times:
  '09:00':
    title: '打卡打卡'
  '11:28':
    title: 'Go! Go! Go! 干饭干饭!'
    text:  '手里活儿停一下'

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
		panic(err.Error())
	}
}

func main() {
	timer := time.NewTicker(time.Duration(config.Interval) * time.Second)
	fmt.Println("启动成功：" + utils.Date.DateFormater(time.Now(), ""))
	for range timer.C {
		nowDate := utils.Date.DateFormater(time.Now(), "hh:mm")
		if record != nowDate {
			val, ok := config.Times[nowDate]
			if ok {
				notify.Push(val.Title, val.Text, "", notificator.UR_CRITICAL)
				cmd := exec.Command("afplay", config.SoundFile)
				cmd.Run()
			}
			record = nowDate
		}
	}
}
