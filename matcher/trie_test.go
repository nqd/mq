// a trie implementation of searching a handler with matching topic
// much code is inspired from https://github.com/tylertreat/fast-topic-matching

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

func TestRabbitMQBinding(t *testing.T) {
	assert := assert.New(t)

	var m = NewTrieMatcher()

	var rabbitmqBinding = []struct {
		topic   string
		handler string
	}{
		{"a.b.c", "t1"},
		{"a.*.c", "t2"},
		{"a.#.b", "t3"},
		{"a.b.b.c", "t4"},
		{"#", "t5"},
		{"#.#", "t6"},
		{"#.b", "t7"},
		{"*.*", "t8"},
		{"a.*", "t9"},
		{"*.b.c", "t10"},
		{"a.#", "t11"},
		{"a.#.#", "t12"},
		{"b.b.c", "t13"},
		{"a.b.b", "t14"},
		{"a.b", "t15"},
		{"b.c", "t16"},
		{"", "t17"},
		{"*.*.*", "t18"},
		{"vodka.martini", "t19"},
		{"a.b.c", "t20"},
		{"*.#", "t21"},
		{"#.*.#", "t22"},
		{"*.#.#", "t23"},
		{"#.#.#", "t24"},
		{"*", "t25"},
		{"#.b.#", "t26"},
	}

	var flagtests = []struct {
		in  string
		out []string
	}{
		{"a.b.c", []string{"t1", "t2", "t5", "t6", "t10", "t11", "t12", "t18", "t20", "t21", "t22", "t23", "t24", "t26"}},
		{"a.b", []string{"t3", "t5", "t6", "t7", "t8", "t9", "t11", "t12", "t15", "t21", "t22", "t23", "t24", "t26"}},
		{"a.b.b", []string{"t3", "t5", "t6", "t7", "t11", "t12", "t14", "t18", "t21", "t22", "t23", "t24", "t26"}},
		{"", []string{"t5", "t6", "t17", "t24"}},
		{"b.c.c", []string{"t5", "t6", "t18", "t21", "t22", "t23", "t24", "t26"}},
		{"a.a.a.a.a", []string{"t5", "t6", "t11", "t12", "t21", "t22", "t23", "t24"}},
		{"vodka.gin", []string{"t5", "t6", "t8", "t21", "t22", "t23", "t24"}},
		{"vodka.martini", []string{"t5", "t6", "t8", "t19", "t21", "t22", "t23", "t24"}},
		{"b.b.c", []string{"t5", "t6", "t10", "t13", "t18", "t21", "t22", "t23", "t24", "t26"}},
		{"nothing.here.at.all", []string{"t5", "t6", "t21", "t22", "t23", "t24"}},
		{"oneword", []string{"t5", "t6", "t21", "t22", "t23", "t24", "t25"}},
	}

}

func assertEqual(assert *assert.Assertions, expected, actual []Handler) {
	assert.Len(actual, len(expected))
	for _, sub := range expected {
		assert.Contains(actual, sub)
	}
}
