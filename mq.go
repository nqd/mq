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
	emit := make(map[string]handlers)
	return &MQ{
		emit: emit,
	}
}

// Publish publishes the data argument to the given subject. The data
// argument is left untouched and needs to be correctly interpreted on
// the receiver.
func (m *MQ) Publish(topic string, data interface{}) (err error) {
	m.Lock()
	hs := m.emit[topic]
	m.Unlock()

	if hs == nil {
		return
	}

	dataType := reflect.TypeOf(data)
	dataValue := []reflect.Value{reflect.ValueOf(data)}

	for _, h := range hs {
		if h.argType == dataType {
			go h.fn.Call(dataValue)
		}
	}

	return nil
}

// Subscribe will create a subscription on the given subject and process incoming
// messages using the specified Handler. The Handler should be a func that matches
// a signature from the description of Handler from above.
func (m *MQ) Subscribe(topic string, cb interface{}) (sub *Subscription, err error) {
	cbType := reflect.TypeOf(cb)
	if cbType.Kind() != reflect.Func {
		panic("mq: Handler needs to be a function")
	}
	if cbType.NumIn() != 1 {
		panic("mq: Handler needs to be a function with one arg")
	}
	cbValue := reflect.ValueOf(cb)

	handler := Handler{
		fn:      cbValue,
		argType: cbType.In(0),
	}

	m.Lock()
	defer m.Unlock()
	m.idCounter++

	if m.emit[topic] == nil {
		m.emit[topic] = handlers{handler}
	} else {
		m.emit[topic] = append(m.emit[topic], handler)
	}

	// fake return
	sub = &Subscription{}

	return
}
