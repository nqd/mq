package mq

import (
	"errors"
	"testing"
	"time"
)

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

	err = m.Publish("topic", []byte("hello world"))
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
		return nil
	case <-time.After(timeout):
		return errors.New("timeout")
	}
}
