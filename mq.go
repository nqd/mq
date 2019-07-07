package mq

import (
	"errors"
	"reflect"
	"sync"

	"github.com/nqd/mq/matcher"
)

// Err
var (
	ErrBadSubscription   = errors.New("invalid subscription")
	ErrBadUnsubscription = errors.New("invalid unsubscription")
	ErrMQClosed          = errors.New("mq is closed")
)

// handler is a specific callback used for Subscribe
type handler struct {
	fn      reflect.Value // value of the cb
	argType reflect.Type  // type of the arg
}

// MQ is the structure that stores hdlr callback with interested topic
type MQ struct {
	sync.Mutex
	matcher matcher.Matcher
	closed  bool
}

// A Subscription represents interest in a given subject.
type Subscription struct {
	mq        *MQ
	operation *matcher.Operation
}

// NewMQ return new structure of MQ
func NewMQ() *MQ {
	matcher := matcher.NewTrieMatcher(matcher.GetDefaultOption())
	return &MQ{
		matcher: matcher,
		closed:  false,
	}
}

// Publish publishes the data argument to the given subject. The data
// argument is left untouched and needs to be correctly interpreted on
// the receiver.
func (m *MQ) Publish(topic string, data interface{}) error {
	m.Lock()

	if m.closed {
		m.Unlock()
		return ErrMQClosed
	}

	hdlrs := m.matcher.Lookup(topic)
	m.Unlock()

	if len(hdlrs) == 0 {
		return nil
	}

	dataType := reflect.TypeOf(data)
	dataValue := []reflect.Value{reflect.ValueOf(data)}

	for _, h := range hdlrs {
		hdlr := h.(handler)
		if hdlr.argType == dataType {
			go hdlr.fn.Call(dataValue)
		}
	}

	return nil
}

// Subscribe will express interest in the given subject. The subject
// can have wildcards (partial:*, full:#). Messages will be delivered
// to the associated cb.
func (m *MQ) Subscribe(topic string, cb interface{}) (*Subscription, error) {
	m.Lock()
	if m.closed {
		m.Unlock()
		return nil, ErrMQClosed
	}
	m.Unlock()

	if cb == nil {
		return nil, ErrBadSubscription
	}

	cbType := reflect.TypeOf(cb)
	if cbType.Kind() != reflect.Func {
		return nil, ErrBadSubscription
	}
	if cbType.NumIn() != 1 {
		return nil, ErrBadSubscription
	}
	cbValue := reflect.ValueOf(cb)

	hdlr := handler{
		fn:      cbValue,
		argType: cbType.In(0),
	}

	m.Lock()
	opt, err := m.matcher.Add(topic, hdlr)
	m.Unlock()

	if err != nil {
		return nil, err
	}

	sub := &Subscription{
		mq:        m,
		operation: opt,
	}

	return sub, nil
}

// Close the given queue.
// After, all publishes will return error
func (m *MQ) Close() {
	m.Lock()
	m.closed = true
	m.Unlock()
}

// Unsubscribe will remove interest in the given subject.
func (s *Subscription) Unsubscribe() error {
	s.mq.Lock()
	defer s.mq.Unlock()

	if err := s.mq.matcher.Remove(s.operation); err != nil {
		return ErrBadUnsubscription
	}

	return nil
}
