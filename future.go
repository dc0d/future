package future

import "sync"

type Interface[T any] interface {
	// Settle executes the future.
	Settle()

	// Settled is a channel that is closed when the future is settled.
	Settled() <-chan struct{}

	// Result blocks until the future is settled.
	Result() (T, error)
}

type settlable[T any] struct {
	once    sync.Once
	settled chan struct{}
	result  T
	err     error
}
