// 打印目录下的所有文件
package main

import (
	"command/src/utils"
	"fmt"
)

func main() {
	arr := utils.File.GetFilesUrl("./")
	for _, val := range arr {
		fmt.Println(val)
	}
}
