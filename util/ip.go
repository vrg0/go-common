package util

import (
	"errors"
	"net"
)

func LocalIp() (string, error) {
	interfaces, err := net.InterfaceAddrs()

	if err != nil {
		return "", err
	}

	for _, address := range interfaces {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}
	return "", errors.New("未找到本机IP地址")

}
