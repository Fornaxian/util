package util

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// NetInfo contains information about the system's networking interfaces
type NetInfo struct {
	Interface    string
	RXBytes      int64
	RXPackets    int64
	RXErrors     int64
	RXDropped    int64
	RXFIFO       int64
	RXFrame      int64
	RXCompressed int64
	RXMulticast  int64
	TXBytes      int64
	TXPackets    int64
	TXErrors     int64
	TXDropped    int64
	TXFIFO       int64
	TXCollisions int64
	TXCarrier    int64
	TXCompressed int64
}

// GetNetInfo reads and parses Linux's /proc/net/dev file
func GetNetInfo(iface string) (inf NetInfo, err error) {
	file, err := os.Open("/proc/net/dev")
	if err != nil {
		return inf, fmt.Errorf("could not open netdev file: %w", err)
	}
	defer file.Close()

	var scanner = bufio.NewScanner(file)
	var split []string

	for linenr := 0; scanner.Scan(); linenr++ {
		if linenr < 2 {
			continue // Skip the first two header lines
		}

		// Split interface and stats
		if split = strings.SplitN(scanner.Text(), ": ", 2); len(split) != 2 {
			return inf, fmt.Errorf("invalid number of interface columns in netdev file. %d instead of 2", len(split))
		}

		// Save the interface name. If it's not the right interface we go to the
		// next one
		if inf.Interface = strings.TrimSpace(split[0]); inf.Interface != iface {
			continue
		}

		break // We found our interface, start parsing
	}

	if inf.Interface != iface {
		return inf, fmt.Errorf("network interface %s not found in netdev file", iface)
	}

	// Split the fields
	if split = strings.Fields(split[1]); len(split) != 16 {
		return inf, fmt.Errorf("invalid number of stat columns in netdev file. %d instead of 16", len(split))
	}

	if inf.RXBytes, err = strconv.ParseInt(split[0], 10, 64); err != nil {
		return inf, err
	}
	if inf.RXPackets, err = strconv.ParseInt(split[1], 10, 64); err != nil {
		return inf, err
	}
	if inf.RXErrors, err = strconv.ParseInt(split[2], 10, 64); err != nil {
		return inf, err
	}
	if inf.RXDropped, err = strconv.ParseInt(split[3], 10, 64); err != nil {
		return inf, err
	}
	if inf.RXFIFO, err = strconv.ParseInt(split[4], 10, 64); err != nil {
		return inf, err
	}
	if inf.RXFrame, err = strconv.ParseInt(split[5], 10, 64); err != nil {
		return inf, err
	}
	if inf.RXCompressed, err = strconv.ParseInt(split[6], 10, 64); err != nil {
		return inf, err
	}
	if inf.RXMulticast, err = strconv.ParseInt(split[7], 10, 64); err != nil {
		return inf, err
	}

	if inf.TXBytes, err = strconv.ParseInt(split[8], 10, 64); err != nil {
		return inf, err
	}
	if inf.TXPackets, err = strconv.ParseInt(split[9], 10, 64); err != nil {
		return inf, err
	}
	if inf.TXErrors, err = strconv.ParseInt(split[10], 10, 64); err != nil {
		return inf, err
	}
	if inf.TXDropped, err = strconv.ParseInt(split[11], 10, 64); err != nil {
		return inf, err
	}
	if inf.TXFIFO, err = strconv.ParseInt(split[12], 10, 64); err != nil {
		return inf, err
	}
	if inf.TXCollisions, err = strconv.ParseInt(split[13], 10, 64); err != nil {
		return inf, err
	}
	if inf.TXCarrier, err = strconv.ParseInt(split[14], 10, 64); err != nil {
		return inf, err
	}
	if inf.TXCompressed, err = strconv.ParseInt(split[15], 10, 64); err != nil {
		return inf, err
	}

	return inf, scanner.Err()
}
