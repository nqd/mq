# MQ [![Build Status](https://secure.travis-ci.org/nqd/mq.png?branch=master)](http://travis-ci.org/nqd/mq)

A Golang message queue that does struct encoding and matches topics like  RabbitMQ.

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

// Close the MQ
m.Close()
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

More concrete example of RabbitMQ could be found at [RabbitMQ tutorial - Topics](https://www.rabbitmq.com/tutorials/tutorial-five-go.html). The topic matching github.com/nqd/mq/matcher implement trie topic routing, inherited much from [fast-topic-matching](https://github.com/tylertreat/fast-topic-matching).

## Benchmark

This is the result running in an Ubuntu 16.04 64bit, i7-5600U CPU @ 2.60GHz × 4.

Subscribe all 1000 topics:

```{bash}
BenchmarkMQSubscribe-4              2000            693580 ns/op          144124 B/op       4001 allocs/op
```

Publish all 100000 messages to a MQ with 1000 subscribers

```{bash}
BenchmarkMQPublish-4                   3         453835682 ns/op        104704792 B/op   1442775 allocs/op
```

## See also

- [fast-topic-matching](https://github.com/tylertreat/fast-topic-matching)
- [qlobber](https://github.com/davedoesdev/qlobber)
- Very fast and scalable topic routing [1](https://www.rabbitmq.com/blog/2010/09/14/very-fast-and-scalable-topic-routing-part-1/), [2](https://www.rabbitmq.com/blog/2011/03/28/very-fast-and-scalable-topic-routing-part-2/)
