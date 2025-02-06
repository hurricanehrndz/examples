package main

import (
	"fmt"
	"log"
	"os"
	"syscall"

	"golang.org/x/net/route"
	"golang.org/x/sys/unix"
)

func main() {
	// Open the AF_ROUTE socket directly using unix.Socket
	fd, err := unix.Socket(unix.AF_ROUTE, unix.SOCK_RAW, unix.AF_UNSPEC)
	if err != nil {
		fmt.Println("Error opening AF_ROUTE socket:", err)
		os.Exit(1)
	}
	defer unix.Close(fd)

	// Create a RouteMessage with RTM_GET type
	rtm := &route.RouteMessage{
		Version: syscall.RTM_VERSION,
		Type:    unix.RTM_GET,
		ID:      uintptr(os.Getpid()),
		Seq:     0,
		Addrs: []route.Addr{
			&route.Inet4Addr{IP: [4]byte{127, 0, 0, 0}},
		},
	}

	// Marshal the message into bytes
	msgBytes, err := rtm.Marshal()
	if err != nil {
		fmt.Println("Error marshaling RouteMessage:", err)
		os.Exit(1)
	}

	// Send the message over the routing socket
	_, err = unix.Write(fd, msgBytes)
	if err != nil {
		fmt.Println("Error writing to AF_ROUTE socket:", err)
		os.Exit(1)
	}

	// Read the response from the routing socket
	var buf [2 << 10]byte
	n, err := unix.Read(fd, buf[:])
	if err != nil {
		fmt.Println("Error reading from AF_ROUTE socket:", err)
		os.Exit(1)
	}
	log.Printf("read %d bytes (% 02x)\n", n, buf[:n])

	// Parse the response messages
	msgs, err := route.ParseRIB(route.RIBTypeRoute, buf[:n])
	if err != nil {
		fmt.Println("Error parsing RIB:", err)
		os.Exit(1)
	}
	routeMsg, ok := msgs[0].(*route.RouteMessage)
	if !ok {
		os.Exit(1)
	}
	netmask, ok := routeMsg.Addrs[2].(*route.Inet4Addr)
	if !ok {
		os.Exit(1)
	}
	fmt.Printf("netmask: %v\n", netmask.IP)
}
