package util

import (
	"runtime"
	"time"

	"fornaxian.tech/log"
)

// DetectPauses runs a continuous loop which detects stalls in the runtime and
// garbage collection cycles
func DetectPauses() {
	var (
		ticker = time.NewTicker(time.Second)

		// GC data
		mstat     runtime.MemStats
		prevPause uint64
	)
	for range ticker.C {
		runtime.ReadMemStats(&mstat)
		if mstat.PauseTotalNs != prevPause {
			if mstat.PauseTotalNs-prevPause > 100e6 { // 100ms
				log.Warn(
					"Long GC detected: pause: %.2fms, total time spent collecting garbage: %0.2fms, heap size: %s",
					float64(mstat.PauseTotalNs-prevPause)/1e6,
					float64(mstat.PauseTotalNs)/1e6,
					FormatData(int64(mstat.HeapInuse)),
				)
			}
			prevPause = mstat.PauseTotalNs
		}
	}
}
