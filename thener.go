package future

// Then returns a settlable future - a *Thener - that is settled when
// the settlable future passed in is settled.
func Then[In, Out any, F Interface[In]](f F, fn func(In) (Out, error)) *Thener[In, Out] {
	return newThener(f, fn)
}

type Thener[In, Out any] struct {
	fn func(In) (Out, error)
	f  Interface[In]
	settlable[Out]
}

func newThener[In, Out any, F Interface[In]](f F, fn func(In) (Out, error)) *Thener[In, Out] {
	result := &Thener[In, Out]{
		fn: fn,
		f:  f,
	}
	result.settled = make(chan struct{})
	return result
}

func (tn *Thener[In, Out]) Settle() {
	tn.f.Settle()
	tn.once.Do(func() {
		defer close(tn.settled)
		previousResult, previousError := tn.f.Result()
		if previousError != nil {
			tn.err = previousError
			return
		}

		tn.result, tn.err = tn.fn(previousResult)
		// NOTE: could be a gc optimization: tn.fn = nil
	})
}

func (tn *Thener[In, Out]) Settled() <-chan struct{} { return tn.settled }

func (tn *Thener[In, Out]) Result() (Out, error) {
	<-tn.settled
	return tn.result, tn.err
}
