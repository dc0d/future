package future

// All is a settlable future that is settled when
// all of the settlable futures are settled - concurrently.
func All[T any](futures ...Interface[T]) Interface[[]Interface[T]] {
	return New(func() ([]Interface[T], error) {
		set := make(map[<-chan struct{}]Interface[T])
		for _, p := range futures {
			go p.Settle()
			set[p.Settled()] = p
		}

		for {
			if len(set) == 0 {
				return futures, nil
			}

			for settled := range set {
				select {
				case <-settled:
					delete(set, settled)
				default:
				}
			}
		}
	})
}
