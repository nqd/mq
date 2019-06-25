# MQ [![Build Status](https://secure.travis-ci.org/nqd/mq.png?branch=master)](http://travis-ci.org/nqd/mq)

An Title Queue with callback for Golang

## Installation

## Basic Usage

```{go}
import "github.com/nqd/mq"

m := mq.NewMQ()

// Async subscription
sub, err := m.Subscribe("foo", func(m []byte) {
    fmt.Printf("Received a message %s\n", string(m))
})
// ...
err := m.Publish("foo", []byte("hello world"))

// MQ with a Golang struct
type todo struct {
    Title string
    Finish    bool
}

sub, err := m.Subscribe("bar", func(t *todo) {
    log.Printf("received a todo %+v\n", t)
    done <- true
})
// ...
mtd := &todo{
    Title:  "get mq work",
    Finish: false,
}
err := m.Publish("bar", mtd)
// ...

// Unsubscribe
err := sub.Unsubscribe()
```

## Wildcard subscription

MQ supports the use of wildcard with amqp-like topics:

- topic is splitted according to delimiter, default `.`,
- wildcard one (default `*`) matches exactly one word,
- and wildcard some (default `#`) matches zero or many words.

Example of wildcard one:

```{go}
_, err := m.Subscribe("foo.*.baz", func(s string) {
    // will receive s = "hello world"
})
_, err := m.Subscribe("foo.*", func(s string) {
    // will not be called
})

err := m.Publish("foo.bar", "hello world")
```

Example of wildcard some:

```{go}
_, err := m.Subscribe("foo.#", func(s string) {
    // will receive s = "hello my world"
})
_, err := m.Subscribe("#", func(s string) {
    // will receive s = "hello my world"
})
})
_, err := m.Subscribe("foo.bar.baz.#", func(s string) {
    // will receive s = "hello my world"
})

err := m.Publish("foo.bar.baz", "hello my world")
```
