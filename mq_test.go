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
func TestSubS(t *testing.T) {
	done := make(chan bool)

	m := NewMQ()

	// subscribe string
	if _, err := m.Subscribe("topic", func(subMsg string) {
		if subMsg != "hello world" {
			t.Fatal("Received wrong message")
		}
		done <- true
	}); err != nil {
		t.Fatal("Failed to subscribe: ", err)
	}

	if err := m.Publish("topic", "hello world"); err != nil {
		t.Fatal("Failed to publish: ", err)
	}

	if err := Wait(done); err != nil {
		t.Fatal("Did not get message")
	}

	// subscribe a struct
	type msg struct {
		x int
		y string
	}
	pubMsg := msg{
		x: 1,
		y: "hello",
	}
	if _, err := m.Subscribe("topic2", func(subMsg msg) {
		if subMsg.x != 1 || subMsg.y != "hello" {
			t.Fatal("Received wrong message")
		}
		done <- true
	}); err != nil {
		t.Fatal("Failed to subscribe: ", err)
	}

	if err := m.Publish("topic2", pubMsg); err != nil {
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
