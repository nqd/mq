package matcher

import (
	"strings"
	"sync"
)

type trieMatcher struct {
	sync.Mutex
	root *node
}

type node struct {
	word     string
	subs     map[Handler]struct{}
	parent   *node
	children map[string]*node
}

func NewTrieMatcher() Matcher {
	return &trieMatcher{
		root: &node{
			subs:     make(map[Handler]struct{}),
			children: make(map[string]*node),
		},
	}
}

func (t *trieMatcher) Add(topic string, hdl Handler) (*Subscription, error) {
	t.Lock()
	curr := t.root
	for _, word := range strings.Split(topic, delimiter) {
		child, ok := curr.children[word]
		if !ok {
			child = &node{
				word:     word,
				subs:     make(map[Handler]struct{}),
				parent:   curr,
				children: make(map[string]*node),
			}
			curr.children[word] = child
		}
		curr = child
	}
	curr.subs[hdl] = struct{}{}
	t.Unlock()

	return &Subscription{topic: topic, handler: hdl}, nil
}

func (t *trieMatcher) Remove(sub *Subscription) error {
	return nil
}

func (t *trieMatcher) Lookup(topic string) []Handler {
	t.Lock()
	var (
		subMap = t.lookup(strings.Split(topic, delimiter), t.root)
		subs   = make([]Handler, len(subMap))
		i      = 0
	)
	t.Unlock()
	for sub := range subMap {
		subs[i] = sub
		i++
	}
	return subs
}

func (t *trieMatcher) lookup(words []string, node *node) map[Handler]struct{} {
	if len(words) == 0 {
		return node.subs
	}
	subs := make(map[Handler]struct{})
	if n, ok := node.children[words[0]]; ok {
		for k, v := range t.lookup(words[1:], n) {
			subs[k] = v
		}
	}
	if n, ok := node.children[wcOne]; ok {
		for k, v := range t.lookup(words[1:], n) {
			subs[k] = v
		}
	}
	return subs
}
