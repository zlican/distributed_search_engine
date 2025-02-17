package utils

import (
	"fmt"
	"net"
)

// GetLocalIP 获取本机网卡IP(内网)
func GetLocalIP() (ipv4 string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", fmt.Errorf("获取网卡地址失败: %v", err)
	}

	for _, addr := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("没有找到有效的IPv4地址")
}
