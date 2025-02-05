package rtmprint

import (
	"fmt"
	"net/netip"
	"log"

	"golang.org/x/net/route"
	"golang.org/x/sys/unix"
)

func LogRouteMessages(msgs []route.Message) {
	for i, msg := range msgs {
		switch msg := msg.(type) {
		default:
			log.Printf("  [%d] %T\n", i, msg)
		case *route.InterfaceMulticastAddrMessage:
			log.Printf("  [%d] InterfaceMulticastAddrMessage: ver=%d, type=%v, flags=0x%x, idx=%v\n",
				i, msg.Version, msg.Type, msg.Flags, msg.Index)
			LogRouteAddrs(msg.Addrs)
		case *route.RouteMessage:
			log.Printf("  [%d] RouteMessage: ver=%d, type=%v, flags=0x%x, idx=%v, id=%v, seq=%v, err=%v\n",
				i, msg.Version, msg.Type, msg.Flags, msg.Index, msg.ID, msg.Seq, msg.Err)
			LogRouteAddrs(msg.Addrs)
		}
	}
}

func LogRouteAddrs(addrs []route.Addr) {
	for i, a := range addrs {
		if a == nil {
			continue
		}
		log.Printf("      %v = %v\n", rtaxName(i), fmtAddr(a))
	}
}

func fmtAddr(a route.Addr) interface{} {
	if a == nil {
		return nil
	}
	switch a := a.(type) {
	case *route.Inet4Addr:
		return netip.AddrFrom4(a.IP)
	case *route.Inet6Addr:
		ip := netip.AddrFrom16(a.IP)
		if a.ZoneID != 0 {
			ip = ip.WithZone(fmt.Sprint(a.ZoneID)) // TODO: look up net.InterfaceByIndex? but it might be changing?
		}
		return ip
	case *route.LinkAddr:
		return fmt.Sprintf("[LinkAddr idx=%v name=%q addr=%x]", a.Index, a.Name, a.Addr)
	default:
		return fmt.Sprintf("%T: %+v", a, a)
	}
}

func rtaxName(i int) string {
	switch i {
	case unix.RTAX_DST:
		return "dst"
	case unix.RTAX_GATEWAY:
		return "gateway"
	case unix.RTAX_NETMASK:
		return "netmask"
	case unix.RTAX_GENMASK:
		return "genmask"
	case unix.RTAX_IFP:
		return "IFP"
	case unix.RTAX_IFA:
		return "IFA"
	case unix.RTAX_AUTHOR:
		return "author"
	case unix.RTAX_BRD:
		return "BRD"
	}
	return fmt.Sprint(i)
}

