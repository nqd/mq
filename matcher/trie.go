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
	return nil
}
