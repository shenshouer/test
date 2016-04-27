package main

import (
	"fmt"
	"strings"
	"time"
)

func main() {
	fmt.Println(time.Now().Format("20060102150405"))
	fmt.Println(time.Now().Add(60 * time.Minute).Format("20060102150405"))
	fmt.Println(time.Now().Add(-10 * time.Minute).Format("20060102150405"))
	fmt.Println(time.Now().Add(-5 * time.Minute).Format("20060102150405"))

	date := "20151122"
	sections := []string{"08:00:59-09:00:00", "13:30:00-14:30:00"}
	fmt.Println("==========================>>")
	for _, v := range sections {
		if svs := strings.Split(v, "-"); len(svs) == 2 {
			starttime := ""
			endtime := ""
			if t1, err := time.Parse("15:04:05", svs[0]); err != nil {
				fmt.Println(err)
				continue
			} else {
				starttime = fmt.Sprintf("%s%s", date, t1.Format("150405"))
			}

			if t1, err := time.Parse("15:04:05", svs[1]); err != nil {
				fmt.Println(err)
				continue
			} else {
				endtime = fmt.Sprintf("%s%s", date, t1.Format("150405"))
			}

			fmt.Println(starttime, endtime)
		}
	}
}
