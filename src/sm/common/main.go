package common

import (
	"os"
)

func ForwardText(text string, callback func([]byte) ([]byte, error)) string {
	data, err := callback([]byte(text))
	if err != nil {
		panic(err)
	}
	return string(data)
}

func ForwardFile(filename, newFilename string, callback func([]byte) ([]byte, error)) {
	// 读取原始文件内容
	content, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	data, err := callback(content)
	if err != nil {
		panic(err)
	}

	// 将加密后的内容写入新文件
	if err := os.WriteFile(newFilename, data, 0644); err != nil {
		panic(err)
	}
}
