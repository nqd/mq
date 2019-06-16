package goglob

import "sync"

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

func (t *trieMatcher) Add(topic string, hdl Handler) error {
	return nil
}
func (t *trieMatcher) Remove(topic string, hdl Handler) error {
	return nil
}
func (t *trieMatcher) Lookup(topic string) []Handler {
	return nil
}