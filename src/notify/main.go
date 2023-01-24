// 一个定时任务，在指定的时间通知一些消息
package main

import (
	"fmt"
	"time"

	"command/src/utils"

	"github.com/0xAX/notificator"
)

var notify *notificator.Notificator
var record string

type notifyType struct {
	title string
	text  string
}

var notifyMap = make(map[string]notifyType)

func init() {
	notify = notificator.New(notificator.Options{
		AppName: "时间提醒",
	})
	notifyMap["09:00"] = notifyType{"打卡打卡", ""}
	notifyMap["11:28"] = notifyType{"Go! Go! Go! 干饭干饭!", "手里活儿停一下"}
	notifyMap["11:30"] = notifyType{"干饭了兄嘚!", ""}
	notifyMap["14:00"] = notifyType{"继续摸鱼了，渔夫", ""}
	notifyMap["17:00"] = notifyType{"差不多该吃饭了", ""}
	notifyMap["18:00"] = notifyType{"下班下班，记得打卡", ""}
	notifyMap["18:30"] = notifyType{"还不下班吗？再不走生痔疮了", ""}
	notifyMap["19:00"] = notifyType{"这么拼，老板给加班费吗？", ""}
	notifyMap["21:30"] = notifyType{"莫要造轮子了，头发快没了", ""}
	notifyMap["23:00"] = notifyType{"这个点儿还不睡，妙！实在是妙", ""}
}

func main() {
	timer := time.NewTicker(time.Duration(time.Second) * 30)
	fmt.Println("启动成功：" + utils.Date.DateFormater(time.Now(), ""))
	for range timer.C {
		nowDate := utils.Date.DateFormater(time.Now(), "hh:mm")
		if record != nowDate {
			val, ok := notifyMap[nowDate]
			if ok {
				notify.Push(val.title, val.text, "", notificator.UR_CRITICAL)
			}
			record = nowDate
		}
	}
}
