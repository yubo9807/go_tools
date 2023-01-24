package utils

import (
	"net"
	"os/exec"
	"strconv"
)

type serverType struct{}

var Server serverType

// 端口是否被占用
func (s *serverType) PortIsOccupy(port int) bool {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return true
	} else {
		listener.Close()
		return false
	}
}

// 返回一个新的端口
func (s *serverType) PortResult(port int) int {
	isOccupy := s.PortIsOccupy(port)
	if isOccupy {
		port++
		return s.PortResult(port)
	}
	return port
}

// 获取 IPv4 地址
func (s *serverType) IPv4() (string, error) {
	addr, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, val := range addr {
		ipnet, ok := val.(*net.IPNet)
		if !(ok && !ipnet.IP.IsLoopback()) {
			continue
		}
		if ipnet.IP.To4() != nil {
			return ipnet.IP.String(), nil
		}
	}
	return "", err
}

func (s *serverType) IPPublic() (string, error) {
	cmd := exec.Command("curl", "ifconfig.me")
	buf, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(buf), nil
}
