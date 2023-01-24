package utils

import (
	"fmt"
	"os"
)

type fileType struct{}

var File fileType

func (f fileType) recursion(folder string, urlArr []string, prefix string) []string {
	files, err := os.ReadDir(prefix + folder)
	if err != nil {
		fmt.Println(err, prefix+folder)
		return urlArr
	}
	if folder != "./" {
		prefix += folder + "/"
	}
	for _, file := range files {
		if file.IsDir() {
			arr := f.recursion(file.Name(), urlArr, prefix)
			urlArr = append(urlArr, arr...)
		} else {
			urlArr = append(urlArr, prefix+file.Name())
		}
	}
	return urlArr
}

// 获取文件夹下的所有文件路径
func (f *fileType) GetFilesUrl(folder string) []string {
	slice := make([]string, 0)
	fileInfo, err := os.Stat(folder)
	if err != nil {
		return slice
	}

	if fileInfo.IsDir() { // 是目录
		return f.recursion(folder, make([]string, 0), "")
	} else { // 是文件
		return append(slice, folder)
	}
}
