package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// Must be root for AF_PACKET
	if os.Geteuid() != 0 {
		fmt.Println("run with sudo")
		os.Exit(1)
	}

	//net interfae
	iface := "wlp2s0"

	fd, err := openRawSocket(iface)
	if err != nil {
		fmt.Println("socket error:", err)
		os.Exit(1)
	}
	defer closeSocket(fd)

	fmt.Println("🔍 packet sniffer running on", iface)
	fmt.Println(strings.Repeat("=", 60))

	buf := make([]byte, 65536)

	for {
		n, err := readPacket(fd, buf)
		if err != nil {
			fmt.Println("read error:", err)
			continue
		}

		parsePacket(buf[:n])
	}
}

func parsePacket(data []byte) {
	if len(data) < 14 {
		return
	}

	dstMAC := net.HardwareAddr(data[0:6])
	srcMAC := net.HardwareAddr(data[6:12])
	etherType := binary.BigEndian.Uint16(data[12:14])

	offset := 14

	// 802.1Q VLAN tag
	if etherType == 0x8100 && len(data) >= 18 {
		etherType = binary.BigEndian.Uint16(data[16:18])
		offset = 18
	}

	fmt.Printf("\n📦 Ethernet\n")
	fmt.Printf("   Src MAC: %s\n", srcMAC)
	fmt.Printf("   Dst MAC: %s\n", dstMAC)
	fmt.Printf("   Type:    0x%04X", etherType)

	switch etherType {
	case 0x0800:
		fmt.Println(" (IPv4)")
		parseIPv4(data[offset:])
	case 0x0806:
		fmt.Println(" (ARP)")
	case 0x86DD:
		fmt.Println(" (IPv6)")
	default:
		fmt.Println(" (unknown)")
	}
}

func parseIPv4(data []byte) {
	if len(data) < 20 {
		return
	}

	version := data[0] >> 4
	ihl := int(data[0]&0x0F) * 4

	if version != 4 || len(data) < ihl {
		return
	}

	totalLen := binary.BigEndian.Uint16(data[2:4])
	ttl := data[8]
	proto := data[9]
	src := net.IPv4(data[12], data[13], data[14], data[15])
	dst := net.IPv4(data[16], data[17], data[18], data[19])

	fmt.Printf("   🌐 IPv4\n")
	fmt.Printf("      Src: %s\n", src)
	fmt.Printf("      Dst: %s\n", dst)
	fmt.Printf("      TTL: %d\n", ttl)
	fmt.Printf("      Len: %d\n", totalLen)

	switch proto {
	case 6:
		fmt.Println("      Proto: TCP")
		parseTCP(data[ihl:])
	case 17:
		fmt.Println("      Proto: UDP")
		parseUDP(data[ihl:])
	case 1:
		fmt.Println("      Proto: ICMP")
	default:
		fmt.Println("      Proto: other")
	}
}

func parseTCP(data []byte) {
	if len(data) < 20 {
		return
	}

	srcPort := binary.BigEndian.Uint16(data[0:2])
	dstPort := binary.BigEndian.Uint16(data[2:4])
	seq := binary.BigEndian.Uint32(data[4:8])
	ack := binary.BigEndian.Uint32(data[8:12])
	flags := data[13]

	fmt.Printf("      🔌 TCP\n")
	fmt.Printf("         %d → %d\n", srcPort, dstPort)
	fmt.Printf("         Seq: %d Ack: %d\n", seq, ack)
	fmt.Printf("         Flags:")

	if flags&0x02 != 0 {
		fmt.Print(" SYN")
	}
	if flags&0x10 != 0 {
		fmt.Print(" ACK")
	}
	if flags&0x01 != 0 {
		fmt.Print(" FIN")
	}
	if flags&0x04 != 0 {
		fmt.Print(" RST")
	}
	if flags&0x08 != 0 {
		fmt.Print(" PSH")
	}

	fmt.Println()
}

func parseUDP(data []byte) {
	if len(data) < 8 {
		return
	}

	srcPort := binary.BigEndian.Uint16(data[0:2])
	dstPort := binary.BigEndian.Uint16(data[2:4])
	length := binary.BigEndian.Uint16(data[4:6])

	fmt.Printf("      📡 UDP\n")
	fmt.Printf("         %d → %d (%d bytes)\n", srcPort, dstPort, length)
}

