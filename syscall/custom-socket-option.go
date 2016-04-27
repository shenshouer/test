package main

import (
	"log"
	"net"
	"os"
	"syscall"
)

// copy from ipsock_posix.go
func ipToSockaddr(family int, ip net.IP, port int) (syscall.Sockaddr, error) {
	switch family {
	case syscall.AF_INET:
		if len(ip) == 0 {
			ip = net.IPv4zero
		}
		if ip = ip.To4(); ip == nil {
			return nil, net.InvalidAddrError("non-IPv4 address")
		}
		s := new(syscall.SockaddrInet4)
		for i := 0; i < net.IPv4len; i++ {
			s.Addr[i] = ip[i]
		}
		s.Port = port
		return s, nil
	case syscall.AF_INET6:
		if len(ip) == 0 {
			ip = net.IPv6zero
		}
		// IPv4 callers use 0.0.0.0 to mean "announce on any available address".
		// In IPv6 mode, Linux treats that as meaning "announce on 0.0.0.0",
		// which it refuses to do.  Rewrite to the IPv6 unspecified address.
		if ip.Equal(net.IPv4zero) {
			ip = net.IPv6zero
		}
		if ip = ip.To16(); ip == nil {
			return nil, net.InvalidAddrError("non-IPv6 address")
		}
		s := new(syscall.SockaddrInet6)
		for i := 0; i < net.IPv6len; i++ {
			s.Addr[i] = ip[i]
		}
		s.Port = port
		return s, nil
	}
	return nil, net.InvalidAddrError("unexpected socket family")
}

// copy from ipsock_posix.go
func probeIPv6Stack() (supportsIPv6, supportsIPv4map bool) {
	var probes = []struct {
		la net.TCPAddr
		ok bool
	}{
		// IPv6 communication capability
		{net.TCPAddr{IP: net.ParseIP("::1")}, false},
		// IPv6 IPv4-mapped address communication capability
		{net.TCPAddr{IP: net.IPv4(127, 0, 0, 1)}, false},
	}

	for i := range probes {
		s, err := syscall.Socket(syscall.AF_INET6, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
		if err != nil {
			continue
		}
		defer closesocket(s)
		syscall.SetsockoptInt(s, syscall.IPPROTO_IPV6, syscall.IPV6_V6ONLY, 0)
		sa, err := ipToSockaddr(syscall.AF_INET6, probes[i].la.IP, probes[i].la.Port)
		if err != nil {
			continue
		}
		err = syscall.Bind(s, sa)
		if err != nil {
			continue
		}
		probes[i].ok = true
	}

	return probes[0].ok, probes[1].ok
}

// copy from ipsock_posix.go
var supportsIPv6, supportsIPv4map = probeIPv6Stack()

// copy from fd_unix.go - this is different for windows and should support both
func closesocket(s int) error {
	return syscall.Close(s)
}

// copy from ipsock_posix.go

func isWildcard(ip net.IP) bool {
	if ip == nil {
		return true
	}
	return ip.IsUnspecified()
}

// copy from tcpsock_posix.go
func family(ip net.IP) int {
	if len(ip) <= net.IPv4len || ip.To4() != nil {
		return syscall.AF_INET
	}
	return syscall.AF_INET6
}

// copy from ipsock_posix.go
func favouriteAddrFamily(net string, laddr *net.TCPAddr, mode string) (int, bool) {
	switch net[len(net)-1] {
	case '4':
		return syscall.AF_INET, false
	case '6':
		return syscall.AF_INET6, true
	}

	if mode == "listen" && isWildcard(laddr.IP) {
		if supportsIPv4map {
			return syscall.AF_INET6, false
		}
		return family(laddr.IP), false
	}

	return family(laddr.IP), false
}

// copy from sock_cloexec.go
func sysSocket(f, t, p int) (int, error) {
	//s, err := syscall.Socket(f, t|syscall.SOCK_NONBLOCK|syscall.SOCK_CLOEXEC, p)
	s, err := syscall.Socket(f, t|syscall.SOCK_RAW, p)
	// The SOCK_NONBLOCK and SOCK_CLOEXEC flags were introduced in
	// Linux 2.6.27.  If we get an EINVAL error, fall back to
	// using socket without them.
	if err == nil || err != syscall.EINVAL {
		return s, err
	}

	// See ../syscall/exec_unix.go for description of ForkLock.
	syscall.ForkLock.RLock()
	s, err = syscall.Socket(f, t, p)
	if err == nil {
		syscall.CloseOnExec(s)
	}
	syscall.ForkLock.RUnlock()
	if err != nil {
		return -1, err
	}
	if err = syscall.SetNonblock(s, true); err != nil {
		syscall.Close(s)
		return -1, err
	}
	return s, nil
}

// copy from sockopt_linux.go - different for windows
func setDefaultSockopts(s, f int, ipv6only bool) error {
	switch f {
	case syscall.AF_INET6:
		if ipv6only {
			syscall.SetsockoptInt(s, syscall.IPPROTO_IPV6, syscall.IPV6_V6ONLY, 1)
		} else {
			// Allow both IP versions even if the OS default
			// is otherwise.  Note that some operating systems
			// never admit this option.
			syscall.SetsockoptInt(s, syscall.IPPROTO_IPV6, syscall.IPV6_V6ONLY, 0)
		}
	}
	// Allow broadcast.
	err := syscall.SetsockoptInt(s, syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1)
	if err != nil {
		return os.NewSyscallError("setsockopt", err)
	}
	return nil
}

// copy from sockopt_linux.go - different for windows
func setDefaultListenerSockopts(s int) error {
	// Allow reuse of recently-used addresses.
	err := syscall.SetsockoptInt(s, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	if err != nil {
		return os.NewSyscallError("setsockopt", err)
	}
	return nil
}

// set up our custom socket
func CreateListenerSocket(netType string, laddr *net.TCPAddr, maxbacklog int) (net.Listener, error) {
	var socketAddr syscall.Sockaddr
	var s int
	var err error

	family, ipv6only := favouriteAddrFamily(netType, laddr, "listen")

	if socketAddr, err = ipToSockaddr(family, laddr.IP, laddr.Port); err != nil {
		return nil, err
	}
	log.Printf("Found local socket address of [%v]", socketAddr)

	if s, err = sysSocket(family, syscall.SOCK_STREAM, 0); err != nil {
		return nil, err
	}

	// set the socket options we want
	if err = setDefaultSockopts(s, family, ipv6only); err != nil {
		closesocket(s)
		return nil, err
	}
	// set the reuse addr
	if err = setDefaultListenerSockopts(s); err != nil {
		closesocket(s)
		return nil, err
	}
	// set the custom reusePort option
	/*	if err = syscall.SetsockoptInt(s,syscall.SOL_SOCKET,15,1); err != nil {
			closesocket(s)
			return 0, err
		}
	*/
	// call bind
	if err = syscall.Bind(s, socketAddr); err != nil {
		closesocket(s)
		return nil, err
	}

	// set up the listener - would have to reimplement  maxbacklog  here for real code
	if err = syscall.Listen(s, maxbacklog); err != nil {
		closesocket(s)
		return nil, err
	}

	file := os.NewFile(uintptr(s), "listener-"+laddr.String())

	var socketListener net.Listener
	if socketListener, err = net.FileListener(file); err != nil {
		file.Close()
		return nil, err
	}
	log.Printf("Got file  handle %s", file.Name())
	file.Close()
	return socketListener, nil

}