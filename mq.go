package mq

type MQ struct{}
type Subscription struct{}

// Handler is a specific callback used for Subscribe
type Handler interface{}

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
func (m *MQ) Subscribe(topic string, cb Handler) (*Subscription, error) {
	return nil, nil
}
