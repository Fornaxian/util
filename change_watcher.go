package util

import (
	"fmt"
	"sync"
	"time"

	"fornaxian.tech/log"
)

// ChangeWatcher watches a database row for changes and relays the events to a
// list of listeners
type ChangeWatcher[T any] struct {
	watchers       map[string]*watcher[T]
	totalListeners int
	changeFunc     ChangeWatcherFunc[T]
	mu             sync.Mutex

	intervalStep time.Duration
	maxInterval  time.Duration
}

type watcher[T any] struct {
	listeners int
	addrem    chan listenerOp[T] // Channel for adding and removing listeners
}

type listenerOp[T any] struct {
	add      bool
	listener chan T
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
type ChangeWatcherFunc[T any] func(id string, previousThing T) (changed bool, thing T)

// NewChangeWatcher creates a new change watcher. The changeFunc is used to
// check whether a change occurred
func NewChangeWatcher[T any](
	changeFunc ChangeWatcherFunc[T],
	intervalStep time.Duration,
	maxInterval time.Duration,
) *ChangeWatcher[T] {
	return &ChangeWatcher[T]{
		watchers:       make(map[string]*watcher[T]),
		totalListeners: 0,
		changeFunc:     changeFunc,
		mu:             sync.Mutex{},

		intervalStep: intervalStep,
		maxInterval:  maxInterval,
	}
}

// Open creates a new change listener for an item. Do not close the channel
// yourself because then the watcher thread will crash. Call Close() instead
func (s *ChangeWatcher[T]) Open(id string) chan T {
	var c = make(chan T)
	s.OpenWithChan(id, c)
	return c
}

func (s *ChangeWatcher[T]) OpenWithChan(id string, c chan T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if the watcher already exists
	w, ok := s.watchers[id]

	if !ok {
		// Watcher does not exist yet. Create it
		w = &watcher[T]{
			listeners: 0,
			addrem:    make(chan listenerOp[T], 4),
		}
		s.watchers[id] = w
		go s.watch(id, w.addrem)
	}

	// Create channel and add it to the watcher. Then return the channel to the
	// listener
	w.listeners++
	s.totalListeners++
	w.addrem <- listenerOp[T]{true, c}
}

// Close closes a channel and removes it from the list of change listeners. If
// this is the last listener for that feed the feed will be removed
func (s *ChangeWatcher[T]) Close(id string, c chan T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	w, ok := s.watchers[id]
	if !ok {
		panic(fmt.Errorf(
			"tried to close channel %v for watcher %s, but watcher doesn't exist",
			c, id,
		))
	}

	w.listeners--
	s.totalListeners--
	w.addrem <- listenerOp[T]{false, c}

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
func (s *ChangeWatcher[T]) Stats() (watchers int, listeners int) {
	s.mu.Lock()
	watchers = len(s.watchers)
	listeners = s.totalListeners
	s.mu.Unlock()
	return watchers, listeners
}

func (s *ChangeWatcher[T]) watch(id string, addrem <-chan listenerOp[T]) {
	var (
		listeners []chan T
		changed   bool
		thing     T
		timeout   = s.intervalStep * 10
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
						"tried to remove channel %v from watcher %s but it doesn't exist",
						lop.listener, id,
					))
				}
			}
			continue
		case <-timer.C:
		}

		// Check if the thing has changed
		changed, thing = s.changeFunc(id, thing)

		// Reset the timer
		timer.Reset(timeout)

		// If it did not change we increase the timeout and sleep again
		if !changed {
			if timeout < s.maxInterval {
				timeout += s.intervalStep
			}
			continue
		}

		// If it did change we decrease the timeout and relay the new thing to
		// all our channels
		if timeout > s.intervalStep {
			timeout -= s.intervalStep
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
