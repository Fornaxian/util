package util

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// MemInfo contains information about the system's random access memory
type MemInfo struct {
	MemTotal          int64
	MemFree           int64
	MemAvailable      int64
	Buffers           int64
	Cached            int64
	SwapCached        int64
	Active            int64
	Inactive          int64
	ActiveAnon        int64
	InactiveAnon      int64
	ActiveFile        int64
	InactiveFile      int64
	Unevictable       int64
	Mlocked           int64
	SwapTotal         int64
	SwapFree          int64
	Dirty             int64
	Writeback         int64
	AnonPages         int64
	Mapped            int64
	Shmem             int64
	KReclaimable      int64
	Slab              int64
	SReclaimable      int64
	SUnreclaim        int64
	KernelStack       int64
	PageTables        int64
	NFSUnstable       int64
	Bounce            int64
	WritebackTmp      int64
	CommitLimit       int64
	CommittedAs       int64
	VmallocTotal      int64
	VmallocUsed       int64
	VmallocChunk      int64
	Percpu            int64
	HardwareCorrupted int64
	AnonHugePages     int64
	ShmemHugePages    int64
	ShmemPmdMapped    int64
	FileHugePages     int64
	FilePmdMapped     int64
	HugePagesTotal    int64
	HugePagesFree     int64
	HugePagesRsvd     int64
	HugePagesSurp     int64
	Hugepagesize      int64
	Hugetlb           int64
	DirectMap4k       int64
	DirectMap2M       int64
	DirectMap1G       int64
}

// GetMemInfo reads and parses Linux's /proc/meminfo file
func GetMemInfo() (mem MemInfo, err error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return mem, fmt.Errorf("could not open meminfo file: %w", err)
	}
	defer file.Close()

	var (
		scanner = bufio.NewScanner(file)
		split   []string
		key     string
		val     string
		ext     string
		num     int64
	)
	for scanner.Scan() {
		if split = strings.SplitN(scanner.Text(), ":", 2); len(split) != 2 {
			return mem, fmt.Errorf("could not split string '%s'", scanner.Text())
		}

		key = strings.TrimSpace(split[0])
		val = strings.TrimSpace(split[1])
		ext = ""

		// Split the number and the extension
		if valSplit := strings.SplitN(val, " ", 2); len(valSplit) == 2 {
			val = valSplit[0]
			ext = valSplit[1]
		}

		// Parse the number
		if num, err = strconv.ParseInt(val, 10, 64); err != nil {
			return mem, fmt.Errorf("failed to parse '%s' as int: %s", val, err)
		}

		if ext == "kB" {
			num *= 1024
		}

		switch key {
		case "MemTotal":
			mem.MemTotal = num
		case "MemFree":
			mem.MemFree = num
		case "MemAvailable":
			mem.MemAvailable = num
		case "Buffers":
			mem.Buffers = num
		case "Cached":
			mem.Cached = num
		case "SwapCached":
			mem.SwapCached = num
		case "Active":
			mem.Active = num
		case "Inactive":
			mem.Inactive = num
		case "Active(anon)":
			mem.ActiveAnon = num
		case "Inactive(anon)":
			mem.InactiveAnon = num
		case "Active(file)":
			mem.ActiveFile = num
		case "Inactive(file)":
			mem.InactiveFile = num
		case "Unevictable":
			mem.Unevictable = num
		case "Mlocked":
			mem.Mlocked = num
		case "SwapTotal":
			mem.SwapTotal = num
		case "SwapFree":
			mem.SwapFree = num
		case "Dirty":
			mem.Dirty = num
		case "Writeback":
			mem.Writeback = num
		case "AnonPages":
			mem.AnonPages = num
		case "Mapped":
			mem.Mapped = num
		case "Shmem":
			mem.Shmem = num
		case "KReclaimable":
			mem.KReclaimable = num
		case "Slab":
			mem.Slab = num
		case "SReclaimable":
			mem.SReclaimable = num
		case "SUnreclaim":
			mem.SUnreclaim = num
		case "KernelStack":
			mem.KernelStack = num
		case "PageTables":
			mem.PageTables = num
		case "NFS_Unstable":
			mem.NFSUnstable = num
		case "Bounce":
			mem.Bounce = num
		case "WritebackTmp":
			mem.WritebackTmp = num
		case "CommitLimit":
			mem.CommitLimit = num
		case "Committed_AS":
			mem.CommittedAs = num
		case "VmallocTotal":
			mem.VmallocTotal = num
		case "VmallocUsed":
			mem.VmallocUsed = num
		case "VmallocChunk":
			mem.VmallocChunk = num
		case "Percpu":
			mem.Percpu = num
		case "HardwareCorrupted":
			mem.HardwareCorrupted = num
		case "AnonHugePages":
			mem.AnonHugePages = num
		case "ShmemHugePages":
			mem.ShmemHugePages = num
		case "ShmemPmdMapped":
			mem.ShmemPmdMapped = num
		case "FileHugePages":
			mem.FileHugePages = num
		case "FilePmdMapped":
			mem.FilePmdMapped = num
		case "HugePages_Total":
			mem.HugePagesTotal = num
		case "HugePages_Free":
			mem.HugePagesFree = num
		case "HugePages_Rsvd":
			mem.HugePagesRsvd = num
		case "HugePages_Surp":
			mem.HugePagesSurp = num
		case "Hugepagesize":
			mem.Hugepagesize = num
		case "Hugetlb":
			mem.Hugetlb = num
		case "DirectMap4k":
			mem.DirectMap4k = num
		case "DirectMap2M":
			mem.DirectMap2M = num
		case "DirectMap1G":
			mem.DirectMap1G = num
		}
	}

	return mem, scanner.Err()
}
