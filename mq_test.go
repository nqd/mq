package mq

import (
	"errors"
	"testing"
	"time"
)

func TestSubErr(t *testing.T) {
	m := NewMQ()

	// todo: test err when cb is not fnc
	// todo: test err when cb arg number is not 1
	if _, err := m.Subscribe("topic", nil); err == nil {
		t.Fatal("Expect an error")
	}

	var s string
	if _, err := m.Subscribe("topic", s); err == nil {
		t.Fatal("Expect an error")
	}
	if _, err := m.Subscribe("topic", func(arg1 string, arg2 string) {}); err == nil {
		t.Fatal("Expect an error")
	}
}
func TestSub(t *testing.T) {
	done := make(chan bool)

	m := NewMQ()

	// todo: test err when cb is not fnc
	// todo: test err when cb arg number is not 1
	_, err := m.Subscribe("topic", func(t string) {
		done <- true
	})
	if err != nil {
		t.Fatal("Failed to subscribe: ", err)
	}

	err = m.Publish("topic", "hello world")
	if err != nil {
		t.Fatal("Failed to publish: ", err)
	}

	if err := Wait(done); err != nil {
		t.Fatal("Did not get message")
	}
}

func Wait(ch chan bool) (err error) {
	timeout := 1 * time.Second
	select {
	case <-ch:
		return
	case <-time.After(timeout):
		err = errors.New("timeout")
	}
	return
}
