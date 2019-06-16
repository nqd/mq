package goglob

const (
	delimiter = "."
	wcOne     = "*"
	wcSome    = "#"
	empty     = ""
)

// Handler is a value associated with a subscription.
type Handler interface{}

// Macher contains topic subscriptions and performs matches on them.
type Macher interface {
	// Subscribe adds the Subscriber to the topic and returns a Subscription.
	Add(topic string, hdl Handler) error

	// Unsubscribe removes the Subscription.
	Remove(topic string, hdl Handler)

	// Lookup returns the Subscribers for the given topic.
	Lookup(topic string) []Handler
}
