package utils

import (
	"strconv"
	"strings"
	"time"
)

type dateType struct{}

var Date dateType

type DateObjType struct {
	Year   int
	Mouth  int
	Day    int
	Hours  int
	Minute int
	Second int
}

// 时间戳转结构体
func (d *dateType) DateToObj(t time.Time) DateObjType {
	return DateObjType{
		Year:   t.Local().Year(),
		Mouth:  int(t.Month()),
		Day:    t.Day(),
		Hours:  t.Hour(),
		Minute: t.Minute(),
		Second: t.Second(),
	}
}

// 十位补全
func (d *dateType) DateZeroize(num int) string {
	str := "0" + strconv.Itoa(num)
	length := len(str)
	return str[length-2 : length]
}

// 时间格式化
func (d *dateType) DateFormater(time time.Time, formater string) string {
	if formater == "" {
		formater = "YYYY-MM-DD hh:mm:ss"
	}
	t := d.DateToObj(time)

	str1 := strings.Replace(formater, "YYYY", strconv.Itoa(t.Year), -1)
	str2 := strings.Replace(str1, "MM", d.DateZeroize(t.Mouth), -1)
	str3 := strings.Replace(str2, "DD", d.DateZeroize(t.Day), -1)
	str4 := strings.Replace(str3, "hh", d.DateZeroize(t.Hours), -1)
	str5 := strings.Replace(str4, "mm", d.DateZeroize(t.Minute), -1)
	str6 := strings.Replace(str5, "ss", d.DateZeroize(t.Second), -1)
	return str6
}
