package util

import (
	"fmt"
	"sync"
	"time"

	"fornaxian.tech/log"
)

// ChangeWatcher watches a database row for changes and relays the events to a
// list of listeners
type ChangeWatcher struct {
	watchers       map[string]*watcher
	totalListeners int
	fn             ChangeWatcherFunc
	mu             sync.Mutex
}

type watcher struct {
	listeners int
	addrem    chan listenerOp // Channel for adding and removing listeners
}

type listenerOp struct {
	add      bool
	listener chan interface{}
}

// ChangeWatcherFunc is the function which periodically checks if a value has
// changed. The ID is the ID of the thing to monitor. previousThing is the last
// value of the thing which was returned, use this to compare if the thing has
// changed. previousThing will be nil in the first run.
//
// If the thing has changed you should return true and the new thing. Else
// return false and the thing which was checked
//
// And an error occurs you should return changed=false and thing=nil
type ChangeWatcherFunc func(id string, previousThing interface{}) (changed bool, thing interface{})

// NewChangeWatcher creates a new change watcher. The changeFunc is used to
// check whether a change occurred
func NewChangeWatcher(changeFunc ChangeWatcherFunc) *ChangeWatcher {
	return &ChangeWatcher{
		watchers:       make(map[string]*watcher),
		totalListeners: 0,
		fn:             changeFunc,
		mu:             sync.Mutex{},
	}
}

// Open creates a new change listener for an item. Do not close the channel
// yourself because then the watcher thread will crash. Call Close() instead
func (s *ChangeWatcher) Open(id string) chan interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if the watcher already exists
	w, ok := s.watchers[id]

	if !ok {
		// Watcher does not exist yet. Create it
		w = &watcher{
			listeners: 0,
			addrem:    make(chan listenerOp, 4),
		}
		s.watchers[id] = w
		go s.watch(id, w.addrem)
	}

	// Create channel and add it to the watcher. Then return the channel to the
	// listener
	var c = make(chan interface{})
	w.listeners++
	s.totalListeners++
	w.addrem <- listenerOp{true, c}
	return c
}

// Close closes a channel and removes it from the list of change listeners. If
// this is the last listener for that feed the feed will be removed
func (s *ChangeWatcher) Close(id string, c chan interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	w, ok := s.watchers[id]
	if !ok {
		panic(fmt.Errorf(
			"Tried to close channel %v for watcher %s, but watcher doesn't exist",
			c, id,
		))
	}

	w.listeners--
	s.totalListeners--
	w.addrem <- listenerOp{false, c}

	if w.listeners == 0 {
		// There are no more listeners. Remove this watcher from the map and
		// close the operation channel. Closing the channel will cause the watch
		// thread to close all the listener channels
		delete(s.watchers, id)
		close(w.addrem)

		log.Debug(
			"No listeners left for watcher %s, stopping thread. %d watchers remain",
			id, len(s.watchers),
		)
	}
}

// Stats returns some statistics about the change watcher. Currently the only
// available stat is the number of watcher threads active
func (s *ChangeWatcher) Stats() (watchers int, listeners int) {
	s.mu.Lock()
	watchers = len(s.watchers)
	listeners = s.totalListeners
	s.mu.Unlock()
	return watchers, listeners
}

func (s *ChangeWatcher) watch(id string, addrem <-chan listenerOp) {
	var (
		listeners []chan interface{}
		changed   bool
		thing     interface{}
		timeout   = time.Second
		timer     = time.NewTimer(timeout)
	)
	for {
		select {
		case lop, ok := <-addrem:
			if !ok {
				// Drain the timer and stop
				if !timer.Stop() {
					<-timer.C
				}

				// Close all the remaining listeners
				for _, listener := range listeners {
					log.Warn("Cleaned up orphan listener %v from watcher %s", listener, id)
					close(listener)
				}

				log.Debug("Change watcher thread %s has stopped", id)
				return
			}

			if lop.add {
				// Add listener to the slice
				listeners = append(listeners, lop.listener)
				log.Debug(
					"Added listener %v to watcher %s. Total listeners %d",
					lop.listener, id, len(listeners),
				)
			} else {
				// Loop over the listeners to see if this one exists
				var found = false
				for k := range listeners {
					if listeners[k] == lop.listener {
						found = true

						// Remove listener from the slice
						listeners = append(listeners[:k], listeners[k+1:]...)
						close(lop.listener)

						log.Debug(
							"Removed listener %v from watcher %s. Total listeners %d",
							lop.listener, id, len(listeners),
						)
						break
					}
				}
				if !found {
					panic(fmt.Errorf(
						"Tried to remove channel %v from watcher %s but it doesn't exist",
						lop.listener, id,
					))
				}
			}
			continue
		case <-timer.C:
		}

		// Check if the thing has changed
		changed, thing = s.fn(id, thing)

		// Reset the timer
		timer.Reset(timeout)

		// If it did not change we increase the timeout and sleep again
		if !changed {
			if timeout < time.Second*10 {
				timeout += time.Millisecond * 100
			}
			continue
		}

		// If it did change we decrease the timeout and relay the new thing to
		// all our channels
		if timeout > time.Millisecond*100 {
			timeout -= time.Millisecond * 100
		}

		// Forward the update to all the listeners
		for _, listener := range listeners {
			// Try to send, but skip if the receiver blocks
			select {
			case listener <- thing:
			default:
			}
		}
	}
}
