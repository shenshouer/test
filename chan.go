package main

import (
	"fmt"
)

func main() {
	var a chan string = make(chan string)
	go func() {
		//close(a)
		//a <- "1"
	}()
	v := <-a
	fmt.Println(v)
}
