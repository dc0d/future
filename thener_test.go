package future

import (
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var _ Interface[string] = &Thener[int, string]{}

func TestThener(t *testing.T) {
	t.Parallel()

	t.Run(`should result expected value when settled`, func(t *testing.T) {
		t.Parallel()

		f := New(func() (int, error) { return 19, nil })
		sut := Then(f, func(n int) (string, error) { return strconv.Itoa(n), nil })

		sut.Settle()
		actualResult, _ := sut.Result()

		expectedResult := "19"
		assert.Equal(t, expectedResult, actualResult)
	})

	t.Run(`should result expected error when settled`, func(t *testing.T) {
		t.Parallel()

		expectedError := errors.New(`expectedError`)
		f := New(func() (int, error) { return 0, expectedError })
		sut := Then(f, func(n int) (string, error) { return strconv.Itoa(n), nil })

		sut.Settle()
		_, actualError := sut.Result()

		assert.Equal(t, expectedError, actualError)
	})

	t.Run(`should be fine to get the result multiple times`, func(t *testing.T) {
		t.Parallel()

		f := New(func() (int, error) { return 19, nil })
		sut := Then(f, func(n int) (string, error) { return strconv.Itoa(n), nil })

		sut.Settle()
		var results []string
		for i := 0; i < 1000; i++ {
			actualResult, _ := sut.Result()
			results = append(results, actualResult)
		}

		expectedResult := "19"
		for _, actualResult := range results {
			assert.Equal(t, expectedResult, actualResult)
		}
	})

	t.Run(`should be fine to settle multiple times`, func(t *testing.T) {
		t.Parallel()

		f := New(func() (int, error) { return 19, nil })
		sut := Then(f, func(n int) (string, error) { return strconv.Itoa(n), nil })

		for i := 0; i < 1000; i++ {
			sut.Settle()
		}
		actualResult, _ := sut.Result()

		expectedResult := "19"
		assert.Equal(t, expectedResult, actualResult)
	})

	t.Run(`Result() should block as long as it is not settled`, func(t *testing.T) {
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
		fn := func(s string) (int, error) {
			close(callingResultStarted)
			defer func() {
				go func() {
					for {
						resultIsCalled <- true
					}
				}()
			}()
			res, err := strconv.Atoi(s)
			if err != nil {
				panic(err)
			}
			return res, nil
		}

		f := New(func() (string, error) { return "19", nil })
		sut := Then(f, fn)

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

	t.Run(`Settled() should block as long as it is not settled`, func(t *testing.T) {
		t.Parallel()

		f := New(func() (int, error) { return 19, nil })
		sut := Then(f, func(n int) (string, error) { return strconv.Itoa(n), nil })

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

	t.Run(`should short circuit on error in case previous item settled with error`, func(t *testing.T) {
		t.Parallel()

		expectedError := errors.New("expectedError")
		prev := New(func() (int, error) { return 0, expectedError })

		notExpectedError := errors.New("notExpectedError")
		sut := Then(prev, func(n int) (int64, error) {
			return int64(n), notExpectedError
		})

		sut.Settle()

		actualValue, actualError := sut.Result()

		assert.Equal(t, expectedError, actualError)
		assert.NotEqual(t, notExpectedError, actualError)
		assert.EqualValues(t, 0, actualValue)
	})

	t.Run(`Value should return the expected value from a chain of items`, func(t *testing.T) {
		t.Parallel()

		step1 := New(func() (int, error) { return 19, nil })
		step2 := Then(step1, func(n int) (int64, error) {
			return int64(n), nil
		})
		step3 := Then(step2, func(n int64) (uint64, error) {
			return uint64(n), nil
		})
		sut := Then(step3, func(n uint64) (string, error) {
			return fmt.Sprint(n), nil
		})

		sut.Settle()

		actualValue, actualError := sut.Result()

		expectedValue := "19"
		assert.EqualValues(t, expectedValue, actualValue)
		assert.NoError(t, actualError)
	})
}
