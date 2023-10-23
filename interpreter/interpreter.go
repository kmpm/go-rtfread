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
	stack *stack
}

func New() (Interpreter, error) {
	ipr := &interpreter{
		stack: newStack(),
	}
	// ipr.stack.setHandler(ipr.handle)
	return ipr, nil
}

func (ipr *interpreter) Value() string {
	return ipr.stack.text
}

func (ipr *interpreter) Read(msg Message) error {
	slog.Debug("read", "msg_type", msg.Type, "msg_value", msg.Value, "msg_param", msg.Param)
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
		// ipr.stack.current().destination = destNormal

	case "keyword":
		err := ipr.handle(msg)
		if err != nil {
			return err
		}
	default:
		slog.Debug("read", "group", ipr.group, "stack", ipr.stack.size(), "type", msg.Type, "value", msg.Value, "param", msg.Param)
	}

	return nil
}

func (ipr *interpreter) handle(msg Message) error {
	slog.Debug("ipr.handle")
	ipr.stack.set(msg.Value, msg.Param)
	switch msg.Value {
	case "par":
		ipr.stack.addString("\n")
	case "fonttbl", "colortbl", "stylesheet":
		ipr.stack.current().destination = 1
	}

	return nil
}

type dest int

const (
	destNormal dest = iota
	destSkip
)

type keyword struct {
	destination dest
}

var keywords = map[string]keyword{
	"fonttbl":  {destSkip},
	"colortbl": {destSkip},
	"ltrpar":   {destSkip},
}
