package utils

import (
	"fmt"
	"net"
)

type LocalNetwork struct {
	Ipv4 string `json:"ipv4"`
	Ipv6 string `json:"ipv6"`
	Mac  string `json:"mac"`
}

// LocalNetworkInfo 获取当前网络信息 (首个启用/非环回/非虚拟的网卡网络)
func LocalNetworkInfo() (info *LocalNetwork, err error) {
	infs, err := net.Interfaces()
	if err != nil {
		return
	}
	var firstNet *net.Interface //获取首个符合条件的网卡网络
	for _, inf := range infs {
		condition1 := inf.Flags&net.FlagUp != 0       //网卡启用状态
		condition2 := inf.Flags&net.FlagLoopback == 0 //非环回网络
		condition3 := len(inf.HardwareAddr) > 0       //非虚拟网络
		if condition1 && condition2 && condition3 {
			firstNet = &inf
			break
		}
	}
	if firstNet == nil {
		err = fmt.Errorf("cannot find local network interface")
		return
	}
	addrs, err := firstNet.Addrs()
	if err != nil {
		return
	}
	info = &LocalNetwork{Mac: firstNet.HardwareAddr.String()}
	for _, addr := range addrs {
		var ipNet *net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ipNet = &v.IP
		case *net.IPAddr:
			ipNet = &v.IP
		}
		if ipNet == nil || ipNet.IsLoopback() || !ipNet.IsGlobalUnicast() {
			continue //环回地址、非全局单播地址(非公网+私网地址) 则跳下一项
		}
		if ipNet.To4() != nil && info.Ipv4 == "" {
			info.Ipv4 = ipNet.String()
		} else if ipNet.To4() == nil && info.Ipv6 == "" {
			info.Ipv6 = ipNet.String()
		}
		if info.Ipv4 != "" && info.Ipv6 != "" {
			break
		}
	}
	return
}
