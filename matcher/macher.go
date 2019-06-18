package matcher

const (
	delimiter = "."
	wcOne     = "*"
	wcSome    = "#"
	empty     = ""
)

// Handler is a value associated with a subscription.
type Handler interface{}

// Subscription represents a topic subscription.
type Subscription struct {
	id      uint32
	topic   string
	handler Handler
}

// Matcher contains topic subscriptions and performs matches on them.
type Matcher interface {
	// Subscribe adds the Subscriber to the topic and returns a Subscription.
	Add(topic string, hdl Handler) (*Subscription, error)

	// Unsubscribe removes the Subscription.
	Remove(*Subscription) error

	// Lookup returns the Subscribers for the given topic.
	Lookup(topic string) []Handler
}
