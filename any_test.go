package future

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAny(t *testing.T) {
	t.Run(`should return when any settled with a result (resolved)`, func(t *testing.T) {
		t.Parallel()

		expectedResult := 19
		items := []Interface[int]{
			New(func() (int, error) { time.Sleep(time.Millisecond * 100); return 119, nil }),
			New(func() (int, error) { time.Sleep(time.Millisecond * 100); return 1119, nil }),
			New(func() (int, error) { return expectedResult, nil }),
		}
		sut := Any(items...)

		sut.Settle()

		actualValue, actualError := sut.Result()

		assert.NoError(t, actualError)
		assert.Equal(t, expectedResult, actualValue)
	})

	t.Run(`should return when any settled with an error (rejected)`, func(t *testing.T) {
		t.Parallel()

		expectedError := errors.New("expectedError")
		items := []Interface[int]{
			New(func() (int, error) { time.Sleep(time.Millisecond * 100); return 119, nil }),
			New(func() (int, error) { time.Sleep(time.Millisecond * 100); return 1119, nil }),
			New(func() (int, error) { return 0, expectedError }),
		}
		sut := Any(items...)

		sut.Settle()

		_, actualError := sut.Result()

		assert.Equal(t, expectedError, actualError)
	})
}
