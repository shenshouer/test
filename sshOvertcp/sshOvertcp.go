package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"test/sshOvertcp/client"
	"test/sshOvertcp/server"
)

var (
	gIsServer   bool   // 是否作为服务器端启动
	gPort       int64  // 服务器端端口
	gServerAddr string // 服务器ip地址，如作为客户端启动，此参数必须配置

	s server.Server
)

func main() {
	flag.BoolVar(&gIsServer, "asServer", false, "是否将服务作为服务器端启动，默认false")
	flag.Int64Var(&gPort, "port", 8888, "服务器端服务端口，默认8888")
	flag.StringVar(&gServerAddr, "serverIp", "127.0.0.1", "服务器端服务端口，默认8888")
	flag.Parse()

	fmt.Println(gIsServer, gServerAddr, gPort)

	if gIsServer {
		s = server.NewServer(gPort)
		var input string
		go func() {
			err := s.Start()
			if err != nil {
				panic(fmt.Errorf("启动server失败:%v", err))
			}
		}()
		for {
			fmt.Scanln(&input)
			s.SendMessage(input)
		}
	} else {
		c := client.NewClient(gServerAddr, gPort)
		c.Connect()
	}

	//handleSignal()
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
