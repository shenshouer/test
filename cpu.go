package main

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"time"
)

func main() {
	cps, err := cpu.CPUPercent(5*time.Second, false)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error for get disk stat :%v", err))
	}

	fmt.Println(cps)
}
