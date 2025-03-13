package util

// CountingSemaphore is a utility which limits the concurrent execution of a
// function. When it is initialized it creates a channel with x capacity and
// fills it with x slots. Every time the Acquire() function is called a slot is
// removed from the channel. When the channel is empty the function will block
// until the Unlock function is called, which puts a new slot into the channel.
type CountingSemaphore struct {
	channel chan struct{}
	slots   int
}

// NewCountingSemaphore creates a new semaphore. The slots parameter is how many
// threads can concurrently execute the function
func NewCountingSemaphore(slots int) (cs *CountingSemaphore) {
	cs = &CountingSemaphore{
		channel: make(chan struct{}, slots),
		slots:   slots,
	}

	// Fill the channel with slots. When the channel is empty the Acquire function
	// will block until a new struct is fed into the channel through Release()
	for range slots {
		cs.channel <- struct{}{}
	}
	return cs
}

// Wait takes all execution slots and releases them again. This ensures that no
// other threads are using the semaphore anymore. This essentially functions as
// the Wait function of a WaitGroup.
func (cs *CountingSemaphore) Wait() {
	// Take all the slots
	for range cs.slots {
		cs.Acquire()
	}

	// Release all the slots
	for range cs.slots {
		cs.Release()
	}
}

// Take a slot
func (cs *CountingSemaphore) Try() (ok bool) {
	select {
	case <-cs.channel:
		return true
	default:
		return false
	}
}

// Take a slot
func (cs *CountingSemaphore) Acquire() { <-cs.channel }

// Release a slot
func (cs *CountingSemaphore) Release() { cs.channel <- struct{}{} }

// Exec obtains an execution slot, runs the provided function concurrently and
// then returns the slot
func (cs *CountingSemaphore) Exec(f func()) {
	cs.Acquire()
	go func() {
		f()
		cs.Release()
	}()
}
