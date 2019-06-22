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
		m = NewTrieMatcher()
	)
	m.Add("#.*.#", "t22")
	m.Add("#.one.#", "t26")

	assertEqual(assert, []Handler{"t22"}, m.Lookup("oneword"))
}

func TestRabbitMQBinding(t *testing.T) {
	assert := assert.New(t)

	m := NewTrieMatcher()

	rabbitmqBinding := []struct {
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

	for _, tt := range rabbitmqBinding {
		_, err := m.Add(tt.topic, tt.handler)
		assert.NoError(err)
	}

	matchings := []struct {
		in  string
		out []Handler
	}{
		{"a.b.c", []Handler{"t1", "t2", "t5", "t6", "t10", "t11", "t12", "t18", "t20", "t21", "t22", "t23", "t24", "t26"}},
		{"a.b", []Handler{"t3", "t5", "t6", "t7", "t8", "t9", "t11", "t12", "t15", "t21", "t22", "t23", "t24", "t26"}},
		{"a.b.b", []Handler{"t3", "t5", "t6", "t7", "t11", "t12", "t14", "t18", "t21", "t22", "t23", "t24", "t26"}},
		// empty string, todo {"", []Handler{"t5", "t6", "t17", "t24"}},
		{"b.c.c", []Handler{"t5", "t6", "t18", "t21", "t22", "t23", "t24", "t26"}},
		{"a.a.a.a.a", []Handler{"t5", "t6", "t11", "t12", "t21", "t22", "t23", "t24"}},
		{"vodka.gin", []Handler{"t5", "t6", "t8", "t21", "t22", "t23", "t24"}},
		{"vodka.martini", []Handler{"t5", "t6", "t8", "t19", "t21", "t22", "t23", "t24"}},
		{"b.b.c", []Handler{"t5", "t6", "t10", "t13", "t18", "t21", "t22", "t23", "t24", "t26"}},
		{"nothing.here.at.all", []Handler{"t5", "t6", "t21", "t22", "t23", "t24"}},
		{"oneword", []Handler{"t5", "t6", "t21", "t22", "t23", "t24", "t25"}},
	}
	for _, tt := range matchings {
		assertEqual(assert, tt.out, m.Lookup(tt.in))
	}
}

func assertEqual(assert *assert.Assertions, expected, actual []Handler) {
	assert.Len(actual, len(expected), actual)
	for _, sub := range expected {
		assert.Contains(actual, sub, actual)
	}
}
