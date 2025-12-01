package utils

import (
	"github.com/likexian/whois"
	"github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
)

func PerformWhoisLookup(ip string) (string, error) {
	res, err := whois.Whois(ip)
	if err != nil {
		return "", err
	}
	return res, nil
}

func GetServiceByPort(port uint32) string {
	pid := FindPIDByPort(port)
	if pid == -1 {
		return "Unknown Service"
	}

	proc, err := process.NewProcess(pid)

	if err != nil {
		return "Unknown Service"
	}

	name, _ := proc.Name()
	exe, _ := proc.Exe()
	username, _ := proc.Username()

	info := "Process Name: " + name + "\n"
	info += "Executable Path: " + exe + "\n"
	info += "User: " + username + "\n"
	return info
}

func FindPIDByPort(port uint32) int32 {
	conns, _ := net.Connections("tcp")
	for _, conn := range conns {
		if conn.Laddr.Port == port {
			return conn.Pid
		}
	}

	udpConns, _ := net.Connections("udp")
	for _, conn := range udpConns {
		if conn.Laddr.Port == port {
			return conn.Pid
		}
	}

	return -1
}
