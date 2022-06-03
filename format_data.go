package util

import "fmt"

type Number interface {
	int | int64 | uint64 | float64
}

// FormatData prints an amount of bytes in a readable rounded amount. The total
// number of digits before and after the decimal point will always be 3.
func FormatData[T Number](size T) string {
	var sizef = float64(size)

	var fmtSize = func(n float64, u string) string {
		var f string
		if n >= 100 {
			f = "%.1f"
		} else if n >= 10 {
			f = "%.2f"
		} else {
			f = "%.3f"
		}
		return fmt.Sprintf(f+" "+u, n)
	}

	if size >= 1e18 {
		// An exabyte is the largest volume of data you can express in a signed
		// 64-bit integer
		return fmtSize(sizef/1e18, "EB")
	} else if sizef >= 1e15 {
		return fmtSize(sizef/1e15, "PB")
	} else if sizef >= 1e12 {
		return fmtSize(sizef/1e12, "TB")
	} else if sizef >= 1e9 {
		return fmtSize(sizef/1e9, "GB")
	} else if sizef >= 1e6 {
		return fmtSize(sizef/1e6, "MB")
	} else if sizef >= 1e3 {
		return fmtSize(sizef/1e3, "kB")
	}
	return fmt.Sprintf("%.0f B", sizef)
}
