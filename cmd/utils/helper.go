package utils

import (
	"fmt"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

func GetSystemStat() (string, string) {
	var memUsage string
	var cpuUsage string

	// Memory usage
	v, err := mem.VirtualMemory()
	if err != nil {
		memUsage = "~"
	} else {
		memUsage = fmt.Sprintf("%.2f%%", v.UsedPercent)
	}

	// CPU usage
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		cpuUsage = "~"
	} else {
		cpuUsage = fmt.Sprintf("%.2f%%", cpuPercent[0])
	}

	return memUsage, cpuUsage
}
