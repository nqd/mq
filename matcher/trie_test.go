package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatcher(t *testing.T) {
	assert := assert.New(t)
	var (
		m  = NewTrieMatcher()
		s0 = 0
		s1 = 1
		s2 = 2
	)
	sub0, err := m.Add("forex.*", s0)
	assert.NoError(err)
	sub1, err := m.Add("*.usd", s0)
	assert.NoError(err)
	sub2, err := m.Add("forex.eur", s0)
	assert.NoError(err)
	sub3, err := m.Add("*.eur", s1)
	assert.NoError(err)
	sub4, err := m.Add("forex.*", s1)
	assert.NoError(err)
	sub5, err := m.Add("trade", s1)
	assert.NoError(err)
	sub6, err := m.Add("*", s2)
	assert.NoError(err)

	assertEqual(assert, []Handler{s0, s1}, m.Lookup("forex.eur"))
	assertEqual(assert, []Handler{s2}, m.Lookup("forex"))
	assertEqual(assert, []Handler{}, m.Lookup("trade.jpy"))
	assertEqual(assert, []Handler{s0, s1}, m.Lookup("forex.jpy"))
	assertEqual(assert, []Handler{s1, s2}, m.Lookup("trade"))

	m.Remove(sub0)
	m.Remove(sub1)
	m.Remove(sub2)
	m.Remove(sub3)
	m.Remove(sub4)
	m.Remove(sub5)
	m.Remove(sub6)

	assertEqual(assert, []Handler{}, m.Lookup("forex.eur"))
	assertEqual(assert, []Handler{}, m.Lookup("forex"))
	assertEqual(assert, []Handler{}, m.Lookup("trade.jpy"))
	assertEqual(assert, []Handler{}, m.Lookup("forex.jpy"))
	assertEqual(assert, []Handler{}, m.Lookup("trade"))
}

func assertEqual(assert *assert.Assertions, expected, actual []Handler) {
	assert.Len(actual, len(expected))
	for _, sub := range expected {
		assert.Contains(actual, sub)
	}
}
