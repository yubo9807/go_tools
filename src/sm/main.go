package main

import (
	"command/src/sm/m2"
	"fmt"
)

func main() {

	// length := len(os.Args)
	// if length <= 1 {
	// 	fmt.Print("请输入命令行参数：\n", "gen | encrypt | decrypt\n")
	// 	return
	// }

	// switch os.Args[1] {
	// case "gen":
	// 	pr, pu, _ := m2.GenerateKeyPair()
	// 	fmt.Print("私钥： ", *pr, "\n公钥： ", *pu, "\n")

	// case "encrypt":
	// 	var filename, newFilename, pu string
	// 	fmt.Print("加密文件路径：")
	// 	fmt.Scanln(&filename)
	// 	fmt.Print("生成文件路径：")
	// 	fmt.Scanln(&newFilename)
	// 	fmt.Print("请输入公钥：")
	// 	fmt.Scanln(&pu)
	// 	m2.EncryptFile(filename, newFilename, pu)

	// case "decrypt":
	// 	var filename, newFilename, pr string
	// 	fmt.Print("解密文件路径：")
	// 	fmt.Scanln(&filename)
	// 	fmt.Print("生成文件路径：")
	// 	fmt.Scanln(&newFilename)
	// 	fmt.Print("请输入私钥：")
	// 	fmt.Scanln(&pr)
	// 	m2.DecryptFile(filename, newFilename, pr)
	// }

	text := []byte("hello")
	pr, pu, err := m2.GenerateKey()
	if err != nil {
		panic(err)
	}
	data, err := m2.Encrypt(text, pu)
	if err != nil {
		panic(err)
	}
	result, err := m2.Decrypt(data, pr)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(result))
}
