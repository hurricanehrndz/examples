package main

/*
#include <sys/socket.h>
#include <net/route.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
	"log"

	"golang.org/x/net/route"
	"golang.org/x/sys/unix"
	"github.com/hurricanehrndz/examples/rtmprint"
)

func main() {
	fd, err := unix.Socket(unix.AF_ROUTE, unix.SOCK_RAW, unix.AF_UNSPEC) // unix.Socket(unix.AF_ROUTE, unix.SOCK_RAW, 0)
	// Filter for only gateway messages
	// flagFitler := unix.RTF_GATEWAY | unix.RTF_GLOBAL
	if err != nil {
		log.Fatal("failed to open socket")
	}
	var buf [2 << 10]byte
	for {
		n, err := unix.Read(fd, buf[:])
		if err != nil {
			log.Printf("%s\n", err)
		}
		if n < unix.SizeofRtMsghdr {
			log.Printf("Network monitor: read from routing socket returned less than expected: %d bytes\n", n)
			continue
		}

		// msg := (*unix.RtMsghdr)(unsafe.Pointer(&buf[0]))
		msgs, err := route.ParseRIB(route.RIBTypeRoute, buf[:n])
		if err != nil {
			log.Printf("read %d bytes (% 02x), failed to parse RIB: %v\n", n, buf[:n], err)
			msg := (*unix.RtMsghdr)(unsafe.Pointer(&buf[0]))
			// ifconfig -v en11 | grep -o -E "index.*"
			log.Println("type", msg.Type, "associated ifp index", msg.Index)
			fmt.Println(" bitmask identifying sockaddrs in msg", msg.Addrs)
		}
		if len(msgs) == 0 {
			log.Printf("read %d bytes with no messages (% 02x)\n", n, buf[:n])
			continue
		}
		log.Printf("read %d bytes, %d messages\n", n, len(msgs))
		rtmprint.LogRouteMessages(msgs)
	}
}
