// 获取本机的 ipv4 地址及公网地址
package main

import (
	"command/src/utils"
	"fmt"
)

func main() {

	// ipv4
	ipv4, err := utils.Server.IPv4()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(ipv4)

	// 公网ip
	ipPublic, err := utils.Server.IPPublic()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(ipPublic)

}
