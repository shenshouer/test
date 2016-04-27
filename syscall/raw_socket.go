package main

import (
	"fmt"
	"os"
	"syscall"
)

func main() {
	// NOTE : replace ... with your own config. for example : fd, err := syscall.Socket(syscall.AFINET, syscall.SOCKSTREAM, syscall.IPPROTO_TCP)
	//fd, err := syscall.Socket(...)
	//file := os.NewFile(uintptr(fd), "")
	//conn, err := net.FileConn(file)

	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)

	if err != nil {
		fmt.Println(err)
	}

	file := os.NewFile(uintptr(fd), "")

	for {
		buffer := make([]byte, 1024)
		num, _ := file.Read(buffer)

		fmt.Printf("% X\n", buffer[:num])
	}
}