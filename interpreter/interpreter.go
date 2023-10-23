package interpreter

import (
	"fmt"
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
			ipr.stack.addRune(int(v))
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
		slog.Warn("read", "group", ipr.group, "stack", ipr.stack.size(), "type", msg.Type, "value", msg.Value, "param", msg.Param)
	}

	return nil
}

func (ipr *interpreter) handle(msg Message) error {
	kw, ok := keywords[msg.Value]
	if !ok {
		slog.Warn("ipr.handle", "keyword", msg.Value, "not_found", true, "stack", ipr.stack.size())
		return nil
	}
	switch kw.kwd {
	case kwdProp:
		ipr.stack.set(msg.Value, msg.Param)
	case kwdDest:
		ipr.setDestination(dest(kw.idx))
	case kwdChar:
		ipr.stack.addRune(kw.idx)
	case kwdSpec:
		return ipr.handleSpecial(msg, ipfn(kw.idx))
	}
	// switch msg.Value {
	// case "par":
	// 	ipr.stack.addString("\n")
	// case "u":
	// 	ipr.stack.addString(string(rune(msg.Param)))
	// case "fonttbl", "colortbl", "stylesheet", "expandedcolortbl", "panose", "xmlnstbl", "operator":
	// 	ipr.stack.current().destination = 1
	// }

	return nil
}

func (ipr *interpreter) setDestination(d dest) {
	slog.Debug("setDestination", "d", d, "stack", ipr.stack.size())
	ipr.stack.current().destination = d
}

func (ipr *interpreter) handleSpecial(msg Message, idx ipfn) error {
	switch idx {
	case ipfnBin:
		slog.Debug("ipfnBin", "binary", msg.Value)
	case ipfnUnicode:
		if msg.Param < 0 {
			msg.Param += 65536
		}
		ipr.stack.addRune(msg.Param)
	default:
		return fmt.Errorf("special function not implemented %v", idx)
	}
	return nil
}
