package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	service := "127.0.0.1"
	protocol := "icmp"
	fmt.Println(1)
	IPAddr, err := net.ResolveIPAddr("ip4", service)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(2)
	IPconn, err := net.ListenIP("ip4:"+protocol, IPAddr)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(3)
	buffer := make([]byte, 1024)

	for { // display the bytes read from IP connection
		fmt.Println(4)
		num, clientAddr, _ := IPconn.ReadFrom(buffer)

		fmt.Println("Reading from : ", clientAddr)
		fmt.Printf("% X\n", buffer[:num])
	}
	fmt.Println(5)
}