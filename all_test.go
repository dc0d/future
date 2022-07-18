package future

import (
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	t.Parallel()

	t.Run(`should return when all settled`, func(t *testing.T) {
		t.Parallel()

		expectedResults := []struct {
			result int
			err    error
		}{
			{19, nil},
			{0, errors.New(`some error`)},
			{119, nil},
			{0, errors.New(`some other error`)},
		}
		var items []Interface[int]
		for _, res := range expectedResults {
			res := res
			items = append(items, New(func() (int, error) {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(250)))
				return res.result, res.err
			}))
		}

		sut := All(items...)

		sut.Settle()

		settledItems, actualError := sut.Result()

		assert.NoError(t, actualError)
		for k, v := range expectedResults {
			res, err := settledItems[k].Result()
			assert.Equal(t, v.result, res)
			assert.Equal(t, v.err, err)
		}
	})
}
