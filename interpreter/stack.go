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
	next         *element
}

type stack struct {
	top     *element
	params  map[string]int
	text    string
	keyword string
	count   int
}

func newStack() *stack {
	return &stack{
		params: make(map[string]int),
	}
}

func (s *stack) size() int {
	return s.count
}

func (s *stack) addString(v string) error {
	if s.top != nil {
		s.top.text += v
		return nil
	}
	s.text += v
	return nil
}
func (s *stack) set(p string, v int) error {
	s.setKeyword(p)
	if s.top != nil {
		s.top.params[p] = v
		return nil
	}
	s.params[p] = v
	return nil
}

func (s *stack) setKeyword(v string) error {
	if s.top != nil {
		s.top.keyword = v
		return nil
	}
	s.keyword = v
	return nil
}

func (s *stack) getKeyword() string {
	if s.top != nil {
		return s.top.keyword
	}
	return s.keyword
}

func (s *stack) push() error {
	e := &element{
		params: make(map[string]int),
	}
	if s.top != nil {
		e.params = s.top.params
	}
	e.next = s.top
	s.top = e
	s.count++
	return nil
}

func (s *stack) pop() (*element, error) {
	if s.top == nil {
		return nil, errors.New("stack is empty")
	}
	e := s.top
	slog.Debug("pop", "count", s.count, "kw", s.getKeyword(), "ignore", e.ignoreOutput, "params", e.params, "text", e.text)
	s.text += e.getText()
	slog.Debug("stack", "text", s.text)
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
