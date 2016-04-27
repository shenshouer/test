package main

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"os"
	"os/signal"
	"syscall"
	"time"
	//"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

func main() {
	for {
		time.Sleep(1 * time.Second)
		v, _ := mem.VirtualMemory()
		datas, _ := cpu.CPUPercent(1*time.Second, true)
		diskData, _ := disk.DiskUsage("/")

		fmt.Printf("Memory:%f%% CPU:%v Disk:%f%% \n", v.UsedPercent, datas, diskData.UsedPercent)
	}

	handleSignal()
}

func handleSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	for sig := range c {
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM:
			return
		}
	}
}
