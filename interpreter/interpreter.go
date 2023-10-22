package interpreter

import (
	"log/slog"
	"strconv"
)

type Interpreter interface {
	Read(msg Message) error
	Value() string
}
type Message struct {
	Type  string
	Value string
	Param int
}

type stateHandler func(msg Message) error

type interpreter struct {
	group int
	state stateHandler
	text  string
	stack *stack
}

func New() (Interpreter, error) {
	ipr := &interpreter{
		stack: newStack(),
	}
	ipr.state = ipr.handle
	return ipr, nil
}

func (ipr *interpreter) Value() string {
	return ipr.stack.text
}

func (ipr *interpreter) Read(msg Message) error {
	switch msg.Type {
	case "group-start":
		ipr.group++
		ipr.stack.push()
	case "group-end":
		ipr.group--
		ipr.stack.pop()
	case "hex-char":
		if ipr.group == 1 {
			v, err := strconv.ParseInt(msg.Value, 16, 32)
			if err != nil {
				return err
			}
			ipr.stack.addString(string(rune(v)))
		}
	case "text":
		ipr.stack.addString(msg.Value)
	case "keyword":
		ipr.stack.set(msg.Value, msg.Param)
		switch msg.Value {
		case "par":
			ipr.stack.addString("\n")
		case "fcharset", "f":
			// ipr.stack.top.ignoreOutput = true
		}
	default:
		slog.Debug("read", "group", ipr.group, "stack", ipr.stack.size(), "type", msg.Type, "value", msg.Value, "param", msg.Param)
	}

	return nil
}

func (ipr *interpreter) handle(msg Message) error {

	return nil
}

type keyword struct {
	wantText bool
}

var keywords = map[string]keyword{
	"rtf":       keyword{wantText: false},
	"fccharset": keyword{wantText: false},
	"par":       keyword{wantText: false},
}
