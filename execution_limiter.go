package util

// ExecutionLimiter is a utility which limits the concurrent execution of a
// function. When it is initialized it creates a channel with x capacity and
// fills it with x booleans. Every time the Lock() function is called a bool is
// removed from the channel. When the channel is empty the function will block
// until the Unlock function is called, which puts a new bool into the channel.
type ExecutionLimiter struct {
	channel chan struct{}
}

// NewExecutionLimiter creates a new Exection Limiter. The numThreads parameter
// is how many threads can concurrently execute the function
func NewExecutionLimiter(numThreads int) (el *ExecutionLimiter) {
	el = &ExecutionLimiter{channel: make(chan struct{}, numThreads)}

	// Fill the channel with bools. Each boolean essentially acts as an
	// execution slot. When the channel is empty the Lock function will block
	// until a new bool is fed into the channel through Unlock
	for i := 0; i < numThreads; i++ {
		el.channel <- struct{}{}
	}
	return el
}

// Stop the ExecutionLimiter. This destroys the channel. Calling Unlock after
// Stop will panic
func (el *ExecutionLimiter) Stop() { close(el.channel) }

// Lock the ExecutionLimiter
func (el *ExecutionLimiter) Lock() { <-el.channel }

// Unlock the ExecutionLimiter
func (el *ExecutionLimiter) Unlock() { el.channel <- struct{}{} }
