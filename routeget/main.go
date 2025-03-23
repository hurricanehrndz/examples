package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/hurricanehrndz/examples/rtmprint"
	"golang.org/x/net/route"
	"golang.org/x/sys/unix"
)

func main() {
	// Open the AF_ROUTE socket directly using unix.Socket
	fd, err := unix.Socket(unix.AF_ROUTE, unix.SOCK_RAW, unix.AF_UNSPEC)
	if err != nil {
		os.Exit(1)
	}
	defer unix.Close(fd)

	pid := os.Getpid()
	// Create a RouteMessage with RTM_GET type
	rtm := route.RouteMessage{
		Version: syscall.RTM_VERSION,
		Type:    unix.RTM_GET,
		Flags:   unix.RTF_IFSCOPE,
		ID:      uintptr(pid),
		Seq:     0,
		// Index:   15,
		Addrs: []route.Addr{
			&route.Inet4Addr{IP: [4]byte{0, 0, 0, 0}},
			nil,
			&route.Inet4Addr{IP: [4]byte{0, 0, 0, 0}},
			nil,
			&route.LinkAddr{},
		},
	}

	// Marshal the message into bytes
	msgBytes, err := rtm.Marshal()
	if err != nil {
		// returns no such process when interface doesn't exist
		// returns no such process when route not in table
		fmt.Println("Error marshaling RouteMessage:", err)
		os.Exit(1)
	}

	// Send the message over the routing socket
	_, err = unix.Write(fd, msgBytes)
	if err != nil {
		fmt.Println("Error writing to AF_ROUTE socket:", err)
		os.Exit(1)
	}

	// wait for reply
	var buf [2 << 10]byte
	for {
		// Read the response from the routing socket
		n, err := unix.Read(fd, buf[:])
		if err != nil {
			fmt.Println("Error reading from AF_ROUTE socket:", err)
			os.Exit(1)
		}
		// Parse the response messages
		msgs, err := route.ParseRIB(route.RIBTypeRoute, buf[:n])
		if err != nil {
			fmt.Println("Error parsing RIB:", err)
			os.Exit(1)
		}
		if len(msgs) == 0 {
			continue
		}
		routeMsg, ok := msgs[0].(*route.RouteMessage)
		// confirm it is a reply to our query
		if !ok || (routeMsg.ID != uintptr(pid) && routeMsg.Seq != 0 && routeMsg.Type != unix.RTM_GET) {
			continue
		}
		rtmprint.LogRouteMessages(msgs)
		break
	}
}
