package mq

import (
	"reflect"
	"sync"
)

// Handler is a specific callback used for Subscribe
type Handler struct {
	fn      reflect.Value // value of the cb
	argType reflect.Type  // type of the arg
}

type handlers []Handler

type MQ struct {
	sync.Mutex
	idCounter int
	emit      map[string]handlers // topic - handler cb
}

type Subscription struct{}

// NewMQ return new structure of MQ
func NewMQ() *MQ {
	return &MQ{}
}

// Publish publishes the data argument to the given subject. The data
// argument is left untouched and needs to be correctly interpreted on
// the receiver.
func (m *MQ) Publish(topic string, data interface{}) error {
	return nil
}

// Subscribe will create a subscription on the given subject and process incoming
// messages using the specified Handler. The Handler should be a func that matches
// a signature from the description of Handler from above.
func (m *MQ) Subscribe(topic string, cb interface{}) (sub *Subscription, err error) {
	m.Lock()
	defer m.Unlock()

	m.idCounter++

	cbType := reflect.TypeOf(cb)
	if cbType.Kind() != reflect.Func {
		panic("mq: Handler needs to be a function")
	}
	if cbType.NumIn() != 1 {
		panic("mq: Handler need to be a function with one arg")
	}
	cbValue := reflect.ValueOf(cb)

	handler := Handler{
		fn:      cbValue,
		argType: cbType.In(0),
	}

	if m.emit[topic] == nil {
		m.emit[topic] = handlers{handler}
	} else {
		m.emit[topic] = append(m.emit[topic], handler)
	}

	// fake return
	sub = &Subscription{}

	return
}
