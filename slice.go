package main

import (
	"fmt"
)

func main() {
	a := []string{"22:00:00-23:00:00", "12:00:00-12:30:00"}
	fmt.Println(fmt.Sprintf("%+v", a))
}
