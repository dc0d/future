package future

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

func ExampleSettlableFuture() {
	f := New(func() (string, error) {
		return "Hello Future!", nil
	})

	f.Settle()

	message, _ := f.Result()
	fmt.Println(message)

	// Output:
	// Hello Future!
}

func ExampleSettlableFuture_Settle_async() {
	f := New(func() (string, error) {
		return "Hello Future!", nil
	})

	// settle asynchronously
	go f.Settle()

	message, _ := f.Result() // blocks until SettlableFuture is settled
	fmt.Println(message)

	// Output:
	// Hello Future!
}

func ExampleSettlableFuture_timeout() {
	f := New(func() (string, error) {
		time.Sleep(time.Minute)
		return "Hello Future!", nil
	})

	go f.Settle()

	select {
	case <-f.Settled():
		message, _ := f.Result()
		fmt.Println(message)
	case <-time.After(time.Millisecond * 50):
		fmt.Println("Timed out!")
	}

	// Output:
	// Timed out!
}

func ExampleAny_timeout() {
	f1 := New(func() (string, error) {
		time.Sleep(time.Minute)
		return "Hello Future! 1", nil
	})
	f2 := New(func() (string, error) {
		time.Sleep(time.Minute)
		return "Hello Future! 2", nil
	})
	timeout := New(func() (string, error) {
		time.Sleep(time.Millisecond * 50)
		return "", errors.New("Timed out!")
	})

	anyFuture := Any[string](f1, f2, timeout)
	anyFuture.Settle()

	_, err := anyFuture.Result()

	fmt.Println(err)

	// Output:
	// Timed out!
}

func ExampleAll() {
	f1 := New(func() (string, error) {
		return "1", nil
	})
	f2 := New(func() (string, error) {
		return "2", nil
	})
	f3 := New(func() (string, error) {
		return "3", nil
	})

	allFeatures := All[string](f1, f2, f3)
	allFeatures.Settle()

	items, _ := allFeatures.Result()
	var parts []string
	for _, item := range items {
		part, _ := item.Result()
		parts = append(parts, part)
	}

	fmt.Println(strings.Join(parts, " - "))

	// Output:
	// 1 - 2 - 3
}

func ExampleThen() {
	step1 := New(fetchData)
	step2 := Then(step1, getItems)
	step3 := Then(step2, itemsToText)

	go step3.Settle()

	text, _ := step3.Result()
	fmt.Println(text)

	// Output:
	// PS5 - D100
}

func fetchData() (data, error) {
	return data{ID: "ID1", Items: []string{"PS5", "D100"}, CreatedAt: time.Now().UTC()}, nil
}

func getItems(d data) ([]string, error) { return d.Items, nil }

func itemsToText(items []string) (string, error) { return strings.Join(items, " - "), nil }

type data struct {
	ID        string    `json:"id"`
	Items     []string  `json:"items,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
