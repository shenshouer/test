package main

import (
	"fmt"
	"github.com/shirou/gopsutil/disk"
	"strings"
)

func main() {
	v, err := disk.DiskPartitions(false)
	if err != nil {
		fmt.Println(err)
	}

	for _, k := range v {
		if strings.HasPrefix(k.Device, "/dev/") {
			u, e := disk.DiskUsage(k.Mountpoint)
			if e != nil {
				fmt.Println(e)
			}

			fmt.Println(k.Device, u.Total, u.UsedPercent)
		}
	}
}
