# MQ [![Build Status](https://secure.travis-ci.org/nqd/mq.png?branch=master)](http://travis-ci.org/nqd/mq)

An Message Queue with callback for Golang

## Installation

## Basic Usage

```{go}
import "github.com/nqd/mq"

m := mq.NewMQ()

m.Subscribe("foo", func(m []byte) {
    fmt.Printf("Received a message %s\n", string(m))
})
```
