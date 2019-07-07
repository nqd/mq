package matcher

const (
	delimiter = "."
	wcOne     = "*"
	wcSome    = "#"
	empty     = ""
)

type Option struct {
	WildcardOne  string
	WildcardSome string
	Delimiter    string
}

// Handler is a value associated with a subscription.
type Handler interface{}

// Operation represents a topic subscription.
type Operation struct {
	topic   string
	handler Handler
}

// Matcher contains topic subscriptions and performs matches on them.
type Matcher interface {
	// Subscribe adds the Subscriber to the topic and returns an Operation.
	Add(topic string, hdl Handler) (*Operation, error)

	// Unsubscribe removes the Operation.
	Remove(*Operation) error

	// Lookup returns the Subscribers for the given topic.
	Lookup(topic string) []Handler
}

func GetDefaultOption() Option {
	return Option{
		WildcardOne:  wcOne,
		WildcardSome: wcSome,
		Delimiter:    delimiter,
	}
}
