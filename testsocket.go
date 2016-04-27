package main

import (
	"net"
	"net/http"
)

func main() {
	net.Listen()
	http.ListenAndServe()
}
