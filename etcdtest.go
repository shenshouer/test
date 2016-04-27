package main

import (
	"fmt"
	eclient "github.com/coreos/go-etcd/etcd"
	"log"
	"time"
)

var (
	distRoot = "distRoot"
	machines = []string{"http://127.0.0.1:4001", "http://127.0.0.1:4002", "http://127.0.0.1:4003"}
)

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	client := eclient.NewClient(machines)

	receiver := make(chan *eclient.Response)
	stop := make(chan bool)

	go func() {
		fmt.Println("distRoot")
		_, err := client.Watch(distRoot, 0, true, receiver, stop)
		if err != nil {
			panic(err)
		}
	}()

	ftpreceiver := make(chan *eclient.Response)
	ftpstop := make(chan bool)
	go func() {
		fmt.Println("ftpRoot")
		_, err := client.Watch("ftpRoot", 0, true, ftpreceiver, ftpstop)
		if err != nil {
			panic(err)
		}
	}()

	/*
		resp, err := client.Get(distRoot, false, false)
		if err != nil {
			log.Println(err)
			return
		}

		for _, node := range resp.Node.Nodes {
			fmt.Println(node.Key, node.Value, node.TTL, node.Dir, node.Expiration, node.Nodes)
		}
	*/
	fmt.Println("sss")
	tick := time.Tick(2 * time.Second)
	for {
		select {
		case resp := <-receiver:
			{
				fmt.Println(resp)
			}
		case ftpResp := <-ftpreceiver:
			{
				fmt.Println("ftp", ftpResp)
			}
		case <-tick:
			{
				fmt.Println("超时")
			}
		}
	}
}
