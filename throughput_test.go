package mq

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"
)

const (
	numSubs = 1000
	numMsgs = 100000
)

var (
	subs = make([]string, numSubs)
	msgs = make([]string, numMsgs)
)

func init() {
	for i := 0; i < numSubs; i++ {
		if i%10 == 0 {
			subs[i] = fmt.Sprintf("*.%d.%d", rand.Intn(10), rand.Intn(10))
		} else if i%25 == 0 {
			subs[i] = fmt.Sprintf("%d.*.%d", rand.Intn(10), rand.Intn(10))
		} else if i%45 == 0 {
			subs[i] = fmt.Sprintf("%d.%d.*", rand.Intn(10), rand.Intn(10))
		} else {
			subs[i] = fmt.Sprintf("%d.%d.%d", rand.Intn(10), rand.Intn(10), rand.Intn(10))
		}
	}
	for i := 0; i < numMsgs; i++ {
		topic := subs[i%numSubs]
		msgs[i] = strings.Replace(topic, "*", strconv.Itoa(rand.Intn(10)), -1)
	}
}

func BenchmarkMQPopulate(b *testing.B) {
	m := NewMQ()

	b.ReportAllocs()
	b.ResetTimer()

	var err error
	for j := 0; j < b.N; j++ {
		for _, sub := range subs {
			if _, err = m.Subscribe(sub, func(subMsg string) {}); err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkMQThroughput(b *testing.B) {
	m := NewMQ()

	for _, sub := range subs {
		if _, err := m.Subscribe(sub, func(subMsg string) {}); err != nil {
			b.Fatal(err)
		}
	}

	b.ReportAllocs()
	b.ResetTimer()

	var err error
	for j := 0; j < b.N; j++ {
		for _, msg := range msgs {
			if err = m.Publish(msg, msg); err != nil {
				b.Fatal(err)
			}
		}
	}
}
