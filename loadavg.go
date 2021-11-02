package util

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// LoadAvg contains information about the system load
type LoadAvg struct {
	Load1Min            float64
	Load5Min            float64
	Load15Min           float64
	CurrentlyScheduling int
	TotalScheduling     int
	LastPID             int
}

// GetLoadAvg reads and parses Linux's /proc/loadavg file
func GetLoadAvg() (load LoadAvg, err error) {
	raw, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return load, fmt.Errorf("could not open loadavg file: %w", err)
	}

	trimmed := strings.TrimSpace(string(raw))

	var split []string
	if split = strings.SplitN(trimmed, " ", 5); len(split) != 5 {
		return load, fmt.Errorf("could not split loadavg string '%s'", trimmed)
	}

	if load.Load1Min, err = strconv.ParseFloat(split[0], 64); err != nil {
		return load, fmt.Errorf("failed to parse 1min load: %w", err)
	}
	if load.Load5Min, err = strconv.ParseFloat(split[1], 64); err != nil {
		return load, fmt.Errorf("failed to parse 5min load: %w", err)
	}
	if load.Load15Min, err = strconv.ParseFloat(split[2], 64); err != nil {
		return load, fmt.Errorf("failed to parse 15min load: %w", err)
	}
	if load.LastPID, err = strconv.Atoi(split[4]); err != nil {
		return load, fmt.Errorf("failed to parse last pid: %w", err)
	}

	if split = strings.SplitN(split[3], "/", 2); len(split) != 2 {
		return load, fmt.Errorf("could not split scheduling string '%s'", split[3])
	}
	if load.CurrentlyScheduling, err = strconv.Atoi(split[0]); err != nil {
		return load, fmt.Errorf("failed to parse currently scheduling tasks: %w", err)
	}
	if load.TotalScheduling, err = strconv.Atoi(split[1]); err != nil {
		return load, fmt.Errorf("failed to parse total scheduling tasks: %w", err)
	}

	return load, nil
}
