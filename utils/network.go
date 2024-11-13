package utils

import (
	"errors"
	"net"
)

type NetworkInfo struct {
	Name string   `json:"name"`
	Mac  string   `json:"mac"`
	Ipv4 []string `json:"ipv4"`
	Ipv6 []string `json:"ipv6"`
}

// NetworkInfoList 获取当前主机网络信息
//
// 注: 每张(启用/非环回/非虚拟)网卡的所有(全局单播)网络
func NetworkInfoList() (list []*NetworkInfo, err error) {
	infs, err := net.Interfaces()
	if err != nil {
		return
	}
	var infNets []*net.Interface
	for _, inf := range infs {
		condition1 := inf.Flags&net.FlagUp != 0       //网卡启用状态
		condition2 := inf.Flags&net.FlagLoopback == 0 //非环回网络
		condition3 := len(inf.HardwareAddr) > 0       //非虚拟网络
		if condition1 && condition2 && condition3 {
			infNets = append(infNets, &inf)
		}
	}
	if len(infNets) == 0 {
		err = errors.New("cannot find local network interface")
		return
	}

	for _, infNet := range infNets {
		addrs, addErr := infNet.Addrs()
		if addErr != nil {
			continue
		}
		info := &NetworkInfo{
			Name: infNet.Name,
			Mac:  infNet.HardwareAddr.String(),
		}
		list = append(list, info)
		for _, addr := range addrs {
			var ipNet *net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ipNet = &v.IP
			case *net.IPAddr:
				ipNet = &v.IP
			}
			if ipNet == nil || !ipNet.IsGlobalUnicast() {
				continue //非全局单播地址(非公网+私网地址) 则跳下一项
			}
			if ipNet.To4() != nil {
				info.Ipv4 = append(info.Ipv4, ipNet.String())
			} else {
				info.Ipv6 = append(info.Ipv6, ipNet.String())
			}
		}
	}
	return
}
