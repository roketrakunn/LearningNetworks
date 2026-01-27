//go:build linux

package main

import (
	"net"
	"syscall"
)

func openRawSocket(ifaceName string) (int, error) {
	fd, err := syscall.Socket(
		syscall.AF_PACKET,
		syscall.SOCK_RAW,
		int(htons(syscall.ETH_P_ALL)),
	)
	if err != nil {
		return 0, err
	}

	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		syscall.Close(fd)
		return 0, err
	}

	sll := &syscall.SockaddrLinklayer{
		Protocol: htons(syscall.ETH_P_ALL),
		Ifindex:  iface.Index,
	}

	if err := syscall.Bind(fd, sll); err != nil {
		syscall.Close(fd)
		return 0, err
	}

	return fd, nil
}

func closeSocket(fd int) {
	_ = syscall.Close(fd)
}

func readPacket(fd int, buf []byte) (int, error) {
	n, _, err := syscall.Recvfrom(fd, buf, 0)
	return n, err
}

func htons(i uint16) uint16 {
	return (i<<8)&0xff00 | i>>8
}

