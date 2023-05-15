package binarytrie

import (
	"net"
)

func parseIpNet(ipNet *net.IPNet) (net.IP, int, error) {
	ip := ipNet.IP.To16()
	if ip == nil {
		return nil, 0, ErrInvalidIPAddress
	}
	subnetSize, bits := ipNet.Mask.Size()
	if bits == net.IPv4len*8 {
		subnetSize += (net.IPv6len - net.IPv4len) * 8
	}
	return ip, subnetSize, nil
}

func extractBits(ip net.IP, position, length int) uint32 {
	if length < 1 || length > 32 || position < 0 || position+length-1 >= len(ip)*8 {
		return 0
	}

	lastBit := position + length - 1
	firstByte := position / 8
	lastByte := lastBit / 8

	// Extract the right bytes.
	rightShift := 7 - lastBit%8
	bits := uint32(ip[lastByte]) >> rightShift
	for i := 1; firstByte <= lastByte-i; i++ {
		bits |= uint32(ip[lastByte-i]) << (8*i - rightShift)
	}

	// Mask unnecessary bits.
	bits &= uint32(0xFFFFFFFF) >> (32 - length)

	return bits
}
