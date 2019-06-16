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

func NewGlob() Glob {
	return &trieMatcher{
		root: &node{
			subs:     make(map[Handler]struct{}),
			children: make(map[string]*node),
		},
	}
}
