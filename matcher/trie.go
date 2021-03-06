// a trie implementation of searching a handler with matching topic
// much code is inspired from https://github.com/tylertreat/fast-topic-matching

package matcher

import (
	"errors"
	"strings"
	"sync"
)

type trieMatcher struct {
	sync.Mutex
	root   *node
	option Option
}

type node struct {
	word     string
	subs     map[Handler]struct{}
	parent   *node
	children map[string]*node
}

func (n *node) orphan() {
	if n.parent == nil {
		// Root
		return
	}
	delete(n.parent.children, n.word)
	if len(n.parent.subs) == 0 && len(n.parent.children) == 0 {
		n.parent.orphan()
	}
}

// NewTrieMatcher returns a default trie structure for matching
func NewTrieMatcher(opt Option) Matcher {
	return &trieMatcher{
		root: &node{
			subs:     make(map[Handler]struct{}),
			children: make(map[string]*node),
		},
		option: opt,
	}
}

func (t *trieMatcher) Add(topic string, hdl Handler) (*Operation, error) {
	t.Lock()
	curr := t.root
	for _, word := range strings.Split(topic, t.option.Delimiter) {
		child, ok := curr.children[word]
		if !ok {
			child = &node{
				word:     word,
				subs:     make(map[Handler]struct{}),
				parent:   curr,
				children: make(map[string]*node),
			}
			// with wildcast some, the child is children itself
			if word == t.option.WildcardSome {
				child.children[word] = child
			}

			curr.children[word] = child
		}

		curr = child
	}
	curr.subs[hdl] = struct{}{}
	t.Unlock()

	return &Operation{topic: topic, handler: hdl}, nil
}

func (t *trieMatcher) Remove(sub *Operation) error {
	t.Lock()
	curr := t.root
	for _, word := range strings.Split(sub.topic, t.option.Delimiter) {
		child, ok := curr.children[word]
		if !ok {
			// Operation doesn't exist.
			t.Unlock()
			return errors.New("invalid unsubscription")
		}
		curr = child
	}
	delete(curr.subs, sub.handler)
	if len(curr.subs) == 0 && len(curr.children) == 0 {
		curr.orphan()
	}
	t.Unlock()

	return nil
}

func (t *trieMatcher) Lookup(topic string) []Handler {
	t.Lock()
	var (
		subMap = t.lookup(strings.Split(topic, t.option.Delimiter), t.root)
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
		hlds := node.subs

		// match("#.b.#", "a.b") == true
		if child, ok := node.children[t.option.WildcardSome]; ok {
			// if the child has only a child itself
			if len(child.children) == 1 {
				for k, v := range child.subs {
					hlds[k] = v
				}
			}
		}
		return hlds
	}

	subs := make(map[Handler]struct{})

	if n, ok := node.children[words[0]]; ok {
		for k, v := range t.lookup(words[1:], n) {
			subs[k] = v
		}
	}
	if words[0] != empty {
		if n, ok := node.children[t.option.WildcardOne]; ok {
			for k, v := range t.lookup(words[1:], n) {
				subs[k] = v
			}
		}
	}

	if n, ok := node.children[t.option.WildcardSome]; ok {
		// check the child of child with words[0]
		// if yes, looking to use grandchild, wcSome count = 0
		// match("a.#.b", "a.b") == true
		if nn, ok := n.children[words[0]]; ok {
			for k, v := range t.lookup(words[1:], nn) {
				subs[k] = v
			}
		}
		// match("a.#.*", "a.b") == true
		if words[0] != empty {
			if nn, ok := n.children[t.option.WildcardOne]; ok {
				for k, v := range t.lookup(words[1:], nn) {
					subs[k] = v
				}
			}
		}

		// always use child
		for k, v := range t.lookup(words[1:], n) {
			subs[k] = v
		}
	}
	return subs
}
