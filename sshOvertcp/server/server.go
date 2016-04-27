package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

type (
	Server interface {
		Start() error                         // 启动服务
		StartClientSSH(client net.Conn) error // 启动客户端的ssh
		SendMessage(msg string)               // 广播消息到客户端
	}
	myServer struct {
		listener net.Listener
		port     int64
		client   []net.Conn
	}
)

func NewServer(port int64) Server {
	ms := &myServer{}
	ms.init(port)
	return ms
}

func (this *myServer) init(port int64) {
	this.port = port
}

func (this *myServer) Start() (err error) {
	addr := fmt.Sprintf(":%d", this.port)

	this.listener, err = net.Listen("tcp", addr)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		clientConn, err := this.listener.Accept()
		if err != nil {
			fmt.Println("error:" + err.Error())
			continue
		}

		go this.handleConnection(clientConn)
	}
}

func (this *myServer) StartClientSSH(clientConn net.Conn) error {
	w := bufio.NewWriter(clientConn)
	if _, err := w.WriteString("ssh"); err != nil {
		fmt.Println("error:", err)
		return err
	}

	return nil
}

func (this *myServer) SendMessage(msg string) {
	for _, client := range this.client {
		w := bufio.NewWriter(client)
		_, err := w.Write([]byte(msg))

		if err != nil {
			fmt.Println(client.LocalAddr(), "Send Message Error:", err)
		}
		w.Flush()
	}
}

func (this *myServer) handleConnection(conn net.Conn) {
	this.client = append(this.client, conn)

	r := bufio.NewReader(conn)
	for {
		line, _, err := r.ReadLine()
		if err == io.EOF {
			continue
		} else {
			fmt.Println("error:" + err.Error())
			// TODO 判断连接是否中断，如果中断，则需要从缓存中将连接删除
		}
		fmt.Printf("%s\n", line)
	}
}
