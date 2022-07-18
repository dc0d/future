package future

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var _ Interface[int] = &SettlableFuture[int]{}

func TestSettlableFuture(t *testing.T) {
	t.Parallel()

	t.Run(`should result expected value when settled`, func(t *testing.T) {
		t.Parallel()

		expectedResult := 19
		sut := New(func() (int, error) { return expectedResult, nil })

		sut.Settle()
		actualResult, _ := sut.Result()

		assert.Equal(t, expectedResult, actualResult)
	})

	t.Run(`should result expected error when settled`, func(t *testing.T) {
		t.Parallel()

		expectedError := errors.New(`expectedError`)
		sut := New(func() (int, error) { return 0, expectedError })

		sut.Settle()
		_, actualError := sut.Result()

		assert.Equal(t, expectedError, actualError)
	})

	t.Run(`should be fine to get the result multiple times`, func(t *testing.T) {
		t.Parallel()

		expectedResult := 19
		sut := New(func() (int, error) { return expectedResult, nil })

		sut.Settle()
		var results []int
		for i := 0; i < 1000; i++ {
			actualResult, _ := sut.Result()
			results = append(results, actualResult)
		}

		for _, actualResult := range results {
			assert.Equal(t, expectedResult, actualResult)
		}
	})

	t.Run(`should be fine to settle multiple times`, func(t *testing.T) {
		t.Parallel()

		expectedResult := 19
		sut := New(func() (int, error) { return expectedResult, nil })

		for i := 0; i < 1000; i++ {
			sut.Settle()
		}
		actualResult, _ := sut.Result()

		assert.Equal(t, expectedResult, actualResult)
	})

	t.Run(`Result() should block as long as SettlableFuture is not settled`, func(t *testing.T) {
		t.Parallel()

		resultIsCalled := make(chan bool)
		callingResultStarted := make(chan struct{})
		go func() {
			for {
				select {
				case resultIsCalled <- false:
				case <-callingResultStarted:
					return
				}
			}
		}()
		fn := func() (int, error) {
			close(callingResultStarted)
			defer func() {
				go func() {
					for {
						resultIsCalled <- true
					}
				}()
			}()
			return 19, nil
		}

		sut := New(fn)

		assert.False(t, <-resultIsCalled)

		sut.Settle()

		assert.Eventually(t, func() bool {
			select {
			case res := <-resultIsCalled:
				return res
			default:
				return false
			}
		}, time.Millisecond*350, time.Millisecond^5)
	})

	t.Run(`Settled() should block as long as SettlableFuture is not settled`, func(t *testing.T) {
		t.Parallel()

		sut := New(func() (int, error) { return 19, nil })

		assert.False(t, func() bool {
			select {
			case <-sut.Settled():
				return true
			default:
				return false
			}
		}())

		sut.Settle()

		assert.True(t, func() bool {
			select {
			case <-sut.Settled():
				return true
			default:
				return false
			}
		}())
	})
}
