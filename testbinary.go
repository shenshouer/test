package main

import (
	//	"encoding/binary"
	//	"encoding/json"
	"fmt"
	"net"
)

type (
	// interface
	Jsoner interface {
		Marshaler() ([]byte, error)
		//		Unmarshal(data []byte) error
	}

	A struct {
		A_str string
	}

	B struct {
		*A
		B_str string
	}
)

func main() {
	//	buf := make([]byte, 100)
	//	binary.BigEndian.PutUint32(buf[0:4], 10)
	//	binary.BigEndian.PutUint32(buf[4:8], 100)
	//	strbyte := []byte("Hello World!")
	//	copy(buf[8:], strbyte)
	//	fmt.Println(string(buf[0:4]), string(buf[4:8]), string(buf[8:]))
	//	fmt.Println("==>>", binary.BigEndian.Uint32(buf[0:4]), binary.BigEndian.Uint32(buf[4:8]))

	//	a := &A{A_str: "a"}
	//	ab, err := a.Marshaler()
	//	fmt.Println(string(ab[:]), err)

	//	b := &B{A: a, B_str: "b"}

	//	bb, err := b.Marshaler()
	//	fmt.Println("===", b, string(bb[:]), err)

	//	cc, err := json.Marshal(b)
	//	fmt.Println("===", b, string(cc[:]), err)

	//	var v interface{}
	//	err = json.Unmarshal(cc, &v)
	//	fmt.Println(err, v)
	//	a2, ok := v.(A)
	//	fmt.Println(a2, ok)

	//	err = json.Unmarshal(cc, a)
	//	fmt.Println(err, a)

	fmt.Println(net.JoinHostPort("", "001"))
}
