package interpreter

import (
	"errors"
	"log/slog"
)

type element struct {
	params       map[string]int
	text         string
	keyword      string
	ignoreOutput bool
	// handler      stateHandler
	next        *element
	destination dest
}

type stack struct {
	top *element

	text    string
	keyword string
	count   int
}

func newStack() *stack {
	return &stack{
		top: &element{
			params: make(map[string]int),
		},
	}
}

func (s *stack) size() int {
	return s.count
}

func (s *stack) current() *element {
	return s.top
}

func (s *stack) addString(v string) error {
	switch s.current().destination {
	case destNormal:
		slog.Debug("stack.addstring", "text", v)
		s.current().text += v
	default:
		slog.Debug("stack.addString.ignoring", "text", v)
	}

	return nil
}
func (s *stack) set(p string, v int) error {
	s.current().keyword = p
	s.current().params[p] = v
	return nil
}

func (s *stack) getKeyword() string {
	return s.current().keyword
}

// func (s *stack) setHandler(h stateHandler) {
// 	slog.Debug("setHandler", "h", h)
// 	s.current().handler = h
// }

// func (s *stack) getHandler() stateHandler {
// 	return s.current().handler
// }

// func (s *stack) handle(msg Message) error {
// 	h := s.getHandler()
// 	slog.Debug("handle", "h", h)
// 	return h(msg)
// }

func (s *stack) push() error {
	e := &element{
		params: make(map[string]int),
	}
	// if s.top != nil {
	// 	e.params = s.top.params
	// }
	e.next = s.top
	// e.handler = s.getHandler()
	e.destination = s.top.destination

	s.top = e
	s.count++
	slog.Debug("push", "count", s.count)
	return nil
}

func (s *stack) pop() (*element, error) {
	if s.top == nil {
		return nil, errors.New("stack is empty")
	}
	e := s.top
	slog.Debug("stack.pop", "count", s.count, "kw", s.getKeyword(), "ignore", e.ignoreOutput, "params", e.params, "text", e.text)
	s.text += e.getText()
	// slog.Debug("stack", "text", s.text)
	s.top = e.next
	s.count--
	return e, nil
}

func (e *element) getText() string {
	if e.ignoreOutput {
		return ""
	}
	return e.text
}
