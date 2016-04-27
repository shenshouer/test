package client

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
	"time"
)

type (
	Client interface {
		Connect()
	}

	myClient struct {
		conn     net.Conn
		serverIp string
		port     int64
		clientId int
		reader   *bufio.Reader
		writer   *bufio.Writer
	}
)

func NewClient(ip string, port int64) Client {
	return &myClient{clientId: time.Now().Nanosecond(), serverIp: ip, port: port}
}

func (this *myClient) Connect() {
	addr := fmt.Sprintf("%s:%d", this.serverIp, this.port)
	var err error
	this.conn, err = net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}

	this.reader = bufio.NewReader(this.conn)
	this.writer = bufio.NewWriter(this.conn)
	buf := make([]byte, 1024)
	for {
		//line, _, err := this.reader.ReadLine()

		n, err := this.conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				continue
			}
			fmt.Println("error:", err)
		}
		msg := string(buf[:n])
		fmt.Println("Message From Server:", msg)
		if msg == "ssh" {
			fmt.Println("准备启动本地ssh服务")
			this.startSSH()
		}
	}
}

func (this *myClient) startSSH() error {
	sshConfig := &ssh.ClientConfig{
		User: "sope",
		Auth: []ssh.AuthMethod{ssh.Password("123456")},
	}
	sshClient, err := ssh.Dial("tcp", "localhost:22", sshConfig)
	if err != nil {
		fmt.Println("连接本地ssh失败:", err)
		return err
	}

	session, err := sshClient.NewSession()
	defer session.Close()
	if err != nil {
		fmt.Println("创建ssh会话失败:", err)
		return err
	}

	go func() {
		_, err = io.Copy(session.Stdout, this.conn)
		if err != nil {
			fmt.Println("执行结构发送失败:", err)
		}
	}()
	go func() {
		_, err = io.Copy(this.conn, session.Stdin)
		if err != nil {
			fmt.Println("获取指令失败:", err)
		}
	}()

	return nil
}
