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

	if err := wait(done); err != nil {
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

	if err := wait(done); err != nil {
		t.Fatal("Did not get message")
	}

	// subscribe a number, 2 times
	done1 := make(chan bool)
	done2 := make(chan bool)
	if _, err := m.Subscribe("topic3", func(subMsg int) {
		if subMsg != 123 {
			t.Fatal("Received wrong message")
		}
		done1 <- true
	}); err != nil {
		t.Fatal("Failed to subscribe: ", err)
	}
	if _, err := m.Subscribe("topic3", func(subMsg int) {
		if subMsg != 123 {
			t.Fatal("Received wrong message")
		}
		done2 <- true
	}); err != nil {
		t.Fatal("Failed to subscribe: ", err)
	}

	if err := m.Publish("topic3", 123); err != nil {
		t.Fatal("Failed to publish: ", err)
	}

	err1 := wait(done1)
	err2 := wait(done2)
	if err1 != nil || err2 != nil {
		t.Fatal("Did not get message")
	}
}

func TestUnsub(t *testing.T) {
	done := make(chan bool)

	m := NewMQ()

	// subscribe string
	sub, err := m.Subscribe("topic", func(subMsg string) {
		done <- true
	})
	if err != nil {
		t.Fatal("Failed to subscribe: ", err)
	}

	if err = sub.Unsubscribe(); err != nil {
		t.Fatal("Failed to unsubscribe: ", err)
	}

	if err := m.Publish("topic", "hello"); err != nil {
		t.Fatal("Failed to publish: ", err)
	}

	if err1 := wait(done); err1 == nil {
		t.Fatal("Unsubscribe actually does not remove subscribe")
	}
}

func TestUnsubErr(t *testing.T) {
	m := NewMQ()

	sub, err := m.Subscribe("topic", func(subMsg string) {})
	if err != nil {
		t.Fatal("Failed to subscribe: ", err)
	}
	if err = sub.Unsubscribe(); err != nil {
		t.Fatal("Failed to unsubscribe: ", err)
	}
	if err = sub.Unsubscribe(); err == nil {
		t.Fatal("Cannot unsubscribe two times")
	}
}

func TestClose(t *testing.T) {
	m := NewMQ()

	m.Close()

	if err := m.Publish("topic", "hello"); err != ErrMQClosed {
		t.Fatal("Publish after close MQ: ", err)
	}

	_, err := m.Subscribe("topic", func(subMsg string) {})
	if err != ErrMQClosed {
		t.Fatal("Subscribe after close MQ: ", err)
	}
}

func wait(ch chan bool) (err error) {
	timeout := 1 * time.Second
	select {
	case <-ch:
		return
	case <-time.After(timeout):
		err = errors.New("timeout")
	}
	return
}
