package utils

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"net"
	"time"
)

type MachineInfo struct {
	Ip  net.IP
	Mac net.HardwareAddr
}

func GetIpMac() []MachineInfo {
	var arr []MachineInfo
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		info := MachineInfo{}
		info.Mac = iface.HardwareAddr
		addrs, err := iface.Addrs()
		if err != nil {
			return nil
		}

		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			info.Ip = ip
		}
		arr = append(arr, info)
	}
	return arr
}

//获取ip
func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}

func CpuPercent() float32 {
	percent, _ := cpu.Percent(time.Second, false)
	cpuUsage := 0.0

	for _, v := range percent {
		cpuUsage += v
	}
	return float32(cpuUsage) / float32(len(percent))
}

func MemoryPercent() float32 {
	memData, _ := mem.VirtualMemory()
	// almost every return value is a struct
	return float32(memData.UsedPercent)
}
