package util

// ExecutionLimiter is a utility which limits the concurrent execution of a
// function. When it is initialized it creates a channel with x capacity and
// fills it with x booleans. Every time the Lock() function is called a bool is
// removed from the channel. When the channel is empty the function will block
// until the Unlock function is called, which puts a new bool into the channel.
type ExecutionLimiter struct {
	channel chan struct{}
	threads int
}

// NewExecutionLimiter creates a new Exection Limiter. The numThreads parameter
// is how many threads can concurrently execute the function
func NewExecutionLimiter(numThreads int) (el *ExecutionLimiter) {
	el = &ExecutionLimiter{
		channel: make(chan struct{}, numThreads),
		threads: numThreads,
	}

	// Fill the channel with slots. When the channel is empty the Lock function
	// will block until a new struct is fed into the channel through Unlock
	for i := 0; i < numThreads; i++ {
		el.channel <- struct{}{}
	}
	return el
}

// Stop the ExecutionLimiter. This destroys the channel. Calling Unlock after
// Stop will panic
func (el *ExecutionLimiter) Stop() { close(el.channel) }

// Drain drains the execution limiter of all slots. This essentially functions
// as the Wait function of a WaitGroup. After Drain the ExecutionLimiter cannot
// be used anymore
func (el *ExecutionLimiter) Drain() {
	for i := 0; i < el.threads; i++ {
		<-el.channel
	}
	el.Stop()
}

// Lock the ExecutionLimiter
func (el *ExecutionLimiter) TryLock() (ok bool) {
	select {
	case <-el.channel:
		return true
	default:
		return false
	}
}

// Lock the ExecutionLimiter
func (el *ExecutionLimiter) Lock() { <-el.channel }

// Unlock the ExecutionLimiter
func (el *ExecutionLimiter) Unlock() { el.channel <- struct{}{} }

// Exec obtains an execution slot, runs the provided function and then returns the slot
func (el *ExecutionLimiter) Exec(f func()) {
	el.Lock()
	f()
	el.Unlock()
}
