package util

import "fmt"

// FormatData prints an amount of bytes in a readable rounded amount. The total
// number of digits before and after the decimal point will always be 3.
func FormatData(size int64) string {
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
		return fmtSize(float64(size)/1e18, "EB")
	} else if size >= 1e15 {
		return fmtSize(float64(size)/1e15, "PB")
	} else if size >= 1e12 {
		return fmtSize(float64(size)/1e12, "TB")
	} else if size >= 1e9 {
		return fmtSize(float64(size)/1e9, "GB")
	} else if size >= 1e6 {
		return fmtSize(float64(size)/1e6, "MB")
	} else if size >= 1e3 {
		return fmtSize(float64(size)/1e3, "kB")
	}
	return fmt.Sprintf("%d B", size)
}
