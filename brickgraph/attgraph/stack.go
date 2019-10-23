package main

type stack struct {
	items []string
}

func newStack() *stack {
	return &stack{}
}

func (s *stack) push(item string) {
	s.items = append(s.items, item)
}

func (s *stack) length() int {
	return len(s.items)
}

func (s *stack) pop() string {
	var ret string
	if s.length() > 0 {
		ret = s.items[len(s.items)-1]
		s.items = s.items[:len(s.items)-1]
	}
	return ret
}

type edgeStack struct {
	items []edge
}

func newEdgeStack() *edgeStack {
	return &edgeStack{}
}

func (s *edgeStack) push(item edge) {
	s.items = append(s.items, item)
}

func (s *edgeStack) length() int {
	return len(s.items)
}

func (s *edgeStack) pop() edge {
	var ret edge
	if s.length() > 0 {
		ret = s.items[len(s.items)-1]
		s.items = s.items[:len(s.items)-1]
	}
	return ret
}