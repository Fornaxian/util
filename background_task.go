package util

import (
	"math/rand"
	"sync"
	"time"

	"fornaxian.tech/log"
)

// NewBackgroundTask runs a function at an interval with some randomness. If the
// interval is 30 seconds and the randomness is 10 seconds the function will run
// every 25 to 35 seconds. The function will register itself to the waitgroup
// when starting and mark itself as done when the stopChannel is closed. When
// the stopChannel is closed the task will run one last time regardless of how
// long ago it ran.
func NewBackgroundTask(
	task func(),
	taskName string,
	wg *sync.WaitGroup,
	stopChannel chan bool,
	interval time.Duration,
	randomness time.Duration,
) {
	wg.Add(1)
	defer wg.Done()

	var calcDuration = func() time.Duration {
		if randomness == 0 {
			return time.Until(time.Now().Add(interval).Truncate(interval))
		}
		return interval - (randomness / 2) + time.Duration(rand.Int63n(randomness.Milliseconds()))*time.Millisecond
	}

	var timer = time.NewTimer(calcDuration())
	for {
		select {
		case <-timer.C:
			task()
			timer.Reset(calcDuration())
		case <-stopChannel:
			log.Info("Stopping task %s", taskName)
			if !timer.Stop() {
				<-timer.C
			}
			task()
			log.Info("Stopped task %s", taskName)
			return
		}
	}
}
