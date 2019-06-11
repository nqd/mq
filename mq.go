package mq

import "sync"

// Handler is a specific callback used for Subscribe
type Handler interface{}

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
func (m *MQ) Subscribe(topic string, cb Handler) (sub *Subscription, err error) {
	m.Lock()
	defer m.Unlock()

	m.idCounter++

	if m.emit[topic] == nil {
		m.emit[topic] = handlers{cb}
	} else {
		m.emit[topic] = append(m.emit[topic], cb)
	}

	sub = &Subscription{}

	return
}
