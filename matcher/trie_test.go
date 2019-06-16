package goglob

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatcher(t *testing.T) {
	assert := assert.New(t)
	var (
		m  = NewMatcher()
		s0 = 0
		s1 = 1
		s2 = 2
	)
	sub0, err := m.Subscribe("forex.*", s0)
	assert.NoError(err)
	sub1, err := m.Subscribe("*.usd", s0)
	assert.NoError(err)
	sub2, err := m.Subscribe("forex.eur", s0)
	assert.NoError(err)
	sub3, err := m.Subscribe("*.eur", s1)
	assert.NoError(err)
	sub4, err := m.Subscribe("forex.*", s1)
	assert.NoError(err)
	sub5, err := m.Subscribe("trade", s1)
	assert.NoError(err)
	sub6, err := m.Subscribe("*", s2)
	assert.NoError(err)

	assertEqual(assert, []Subscriber{s0, s1}, m.Lookup("forex.eur"))
	assertEqual(assert, []Subscriber{s2}, m.Lookup("forex"))
	assertEqual(assert, []Subscriber{}, m.Lookup("trade.jpy"))
	assertEqual(assert, []Subscriber{s0, s1}, m.Lookup("forex.jpy"))
	assertEqual(assert, []Subscriber{s1, s2}, m.Lookup("trade"))

	m.Unsubscribe(sub0)
	m.Unsubscribe(sub1)
	m.Unsubscribe(sub2)
	m.Unsubscribe(sub3)
	m.Unsubscribe(sub4)
	m.Unsubscribe(sub5)
	m.Unsubscribe(sub6)

	assertEqual(assert, []Subscriber{}, m.Lookup("forex.eur"))
	assertEqual(assert, []Subscriber{}, m.Lookup("forex"))
	assertEqual(assert, []Subscriber{}, m.Lookup("trade.jpy"))
	assertEqual(assert, []Subscriber{}, m.Lookup("forex.jpy"))
	assertEqual(assert, []Subscriber{}, m.Lookup("trade"))
}
