package future

// SettlableFuture is a future that can be settled by the caller.
type SettlableFuture[T any] struct {
	settlable[T]

	fn func() (T, error)
}

func New[T any](fn func() (T, error)) *SettlableFuture[T] {
	result := &SettlableFuture[T]{
		fn: fn,
	}
	result.settled = make(chan struct{})
	return result
}

func (sf *SettlableFuture[T]) Settle() {
	sf.once.Do(func() {
		defer close(sf.settled)
		sf.result, sf.err = sf.fn()
		// NOTE: could be a gc optimization: sf.fn = nil
	})
}

func (sf *SettlableFuture[T]) Settled() <-chan struct{} { return sf.settled }

func (sf *SettlableFuture[T]) Result() (T, error) {
	<-sf.settled
	return sf.result, sf.err
}
