package future

// Any is a settlable future that is settled when
// any of the settlable futures passed in is settled - concurrently.
func Any[T any](futures ...Interface[T]) Interface[T] {
	return New(func() (T, error) {
		set := make(map[<-chan struct{}]Interface[T])
		for _, p := range futures {
			go p.Settle()
			set[p.Settled()] = p
		}

		for {
			for settled, p := range set {
				select {
				case <-settled:
					return p.Result()
				default:
				}
			}
		}
	})
}
