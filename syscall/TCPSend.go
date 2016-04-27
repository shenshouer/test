package main

import (
	"log"
	"net"
	"syscall"
	"golang.org/x/net/ipv4"
	"math/rand"
)

func main() {
	log.SetFlags(log.Flags()|log.Lshortfile)
	var err error
	//fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
	if err != nil{
		log.Fatal(err)
	}
	addr := syscall.SockaddrInet4{
		Port: 0,
		Addr: [4]byte{127, 0, 0, 1},
	}
	p := pkt()
	err = syscall.Sendto(fd, p, 0, &addr)
	if err != nil {
		log.Fatal("Sendto:", err)
	}
}

func pkt() []byte {
	packet := TCPHeader{
		Source: 0xaa47, // Random ephemeral port
		Destination: 80,
		SeqNum: rand.Uint32(),
		AckNum: 0,
		DataOffset: 5, // 4 bits
		Reserved: 0, // 3 bits
		ECN: 0, // 3 bits
		Ctrl: 2, // 6 bits (000010, SYN bit set)
		Window: 0xaaaa, // size of your receive window
		Checksum: 0, // Kernel will set this if it's 0
		Urgent: 0,
		Options: []TCPOption{},
	}

	data := packet.Marshal()
	packet.Checksum = csum(data, to4byte(laddr), to4byte(raddr))
	data = packet.Marshal()

	conn, err := net.Dial("ip4:tcp", raddr)
	if err != nil {
		log.Fatalf("Dial: %s\n", err)
	}

	conn.Write(data)

	//h := ipv4.Header{
	//	Version:  4,
	//	Len:      20,
	//	TotalLen: 20 + 10, // 20 bytes for IP, 10 for ICMP
	//	TTL:      64,
	//	Protocol: 1, // ICMP
	//	Dst:      net.IPv4(127, 0, 0, 1),
	//	// ID, Src and Checksum will be set for us by the kernel
	//}
	//
	//icmp := []byte{
	//	8, // type: echo request
	//	0, // code: not used by echo request
	//	0, // checksum (16 bit), we fill in below
	//	0,
	//	0, // identifier (16 bit). zero allowed.
	//	0,
	//	0, // sequence number (16 bit). zero allowed.
	//	0,
	//	0xC0, // Optional data. ping puts time packet sent here
	//	0xDE,
	//}
	//cs := csum(icmp)
	//icmp[2] = byte(cs)
	//icmp[3] = byte(cs >> 8)
	//
	//out, err := h.Marshal()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//return append(out, icmp...)
}

func csum(b []byte) uint16 {
	var s uint32
	for i := 0; i < len(b); i += 2 {
		s += uint32(b[i+1])<<8 | uint32(b[i])
	}
	// add back the carry
	s = s>>16 + s&0xffff
	s = s + s>>16
	return uint16(^s)
}