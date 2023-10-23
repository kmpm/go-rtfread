package parser

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/kmpm/go-rtfread/internal"
	"github.com/kmpm/go-rtfread/interpreter"
)

var debug bool = strings.Contains(os.Getenv("DEBUG"), "rtfread")

type stateParser func(ch byte) error

type Parser struct {
	parserState  stateParser
	pos          int
	text         *bytes.Buffer
	groups       int
	hexChar      string
	unicodeChar  string
	keyWord      string
	keyWordParam string
	ipr          interpreter.Interpreter
	done         chan struct{}
	err          error
}

func New(ipr interpreter.Interpreter) (*Parser, error) {
	p := &Parser{
		text: &bytes.Buffer{},
		ipr:  ipr,
		done: make(chan struct{}),
	}
	p.parserState = p.parseText

	return p, nil
}

func (p *Parser) Done() <-chan struct{} {
	return p.done
}

func (p *Parser) Parse(ctx context.Context, r *bufio.Reader) error {

	defer func() {
		if debug {
			slog.Debug("done parsing")
		}
		close(p.done)
	}()

	for {
		ch, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			p.err = err
			return err
		}
		p.pos++
		if err := p.parserState(ch); err != nil {
			p.err = err
			return err
		}
	}
}

func (p *Parser) writeByte(ch byte) error {
	return p.text.WriteByte(ch)
}

// func (p *Parser) writeRune(r rune) error {
// 	_, err := p.text.WriteRune(r)
// 	return err
// }

func (p *Parser) writeString(s string) error {
	_, err := p.text.WriteString(s)
	return err
}

func (p *Parser) push(m interpreter.Message) {
	// slog.Debug("push", "pos", p.pos, "group", p.groups, "type", m.Type, "value", m.Value, "param", m.Param)
	err := p.ipr.Read(m)
	if err != nil {
		panic(err)
	}
}

func (p *Parser) parseText(ch byte) error {
	switch ch {
	case '\r', '\n':
		//noop
	case '\\':
		p.parserState = p.parseEscape
	case '{':
		p.emitStartGroup()
	case '}':
		p.emitEndGroup()
	default:
		err := p.writeByte(ch)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) parseEscape(ch byte) error {
	if ch == '\\' || ch == '{' || ch == '}' {
		p.writeByte(ch)
		p.parserState = p.parseText
		return nil
	}
	p.parserState = p.parseSymbol
	return p.parseSymbol(ch)
}

func (p *Parser) parseSymbol(ch byte) error {
	switch ch {
	case '*':
		p.emitIgnorable()
		p.parserState = p.parseText
	case ':': // sub index
		p.emitIndexSubEntry()
		p.parserState = p.parseText
	case '_': // non-breaking hyphen
		p.writeString("_")
	case '|': //formula
		p.emitFormula()
		p.parserState = p.parseText
	case '~': // non-breaking space
		p.writeByte(' ')
		p.parserState = p.parseText
	case '-': // soft hyphen
		p.writeString("-")
	case '\'':
		p.parserState = p.parseHex
	// case 'u':
	// 	p.parserState = p.parseUnicode
	case '\r', '\n':
		p.emitEndParagraph()
		p.parserState = p.parseText
	default:
		p.parserState = p.parseKeyword
		return p.parseKeyword(ch)
	}
	return nil
}

func (p *Parser) parseHex(ch byte) error {
	if !internal.IsHex(ch) {
		p.parserState = p.parseText
		return fmt.Errorf("invalid hex character: %c", ch)
	}
	p.hexChar += string(ch)
	if len(p.hexChar) >= 2 {
		p.emitHexChar()
		p.parserState = p.parseText
	}
	return nil
}

// func (p *Parser) parseUnicode(ch byte) error {
// 	if internal.IsDigit(ch) {
// 		p.unicodeChar += string(ch)
// 	} else {
// 		p.emitUnicode()
// 		p.parserState = p.parseText
// 	}
// }

func (p *Parser) parseKeyword(ch byte) error {
	if ch == ' ' {
		p.emitKeyword()
		p.parserState = p.parseText
	} else if ch == '-' || internal.IsDigit(ch) {
		p.parserState = p.parseKeywordParam
		p.keyWordParam += string(ch)
	} else if internal.IsAlpha(ch) {
		p.keyWord += string(ch)
	} else {
		p.emitKeyword()
		p.parserState = p.parseText
		return p.parseText(ch)
	}

	return nil
}

func (p *Parser) parseKeywordParam(ch byte) error {
	if internal.IsDigit(ch) {
		p.keyWordParam += string(ch)
	} else if ch == ' ' || ch == '?' {
		p.emitKeyword()
		p.parserState = p.parseText
	} else {
		p.emitKeyword()
		p.parserState = p.parseText
		return p.parserState(ch)
	}
	return nil
}

func (p *Parser) emitText() {
	if p.text.Len() > 0 {
		p.push(interpreter.Message{
			Type:  "text",
			Value: p.text.String(),
		})
		p.text.Reset()
	}
}

func (p *Parser) emitIgnorable() {
	p.emitText()
	p.push(interpreter.Message{Type: "ignorable"})
}

func (p *Parser) emitHexChar() {
	p.emitText()
	// v, err := strconv.ParseInt(p.hexChar, 16, 32)
	// if err != nil {
	// 	return err
	// }
	p.push(interpreter.Message{Type: "hex-char", Value: p.hexChar})
	p.hexChar = ""
}

// func (p *Parser) emitUnicode() {
// 	p.emitText()
// 	p.push(interpreter.Message{Type: "unicode", Value: p.unicodeChar})
// 	p.unicodeChar = ""
// }

func (p *Parser) emitStartGroup() {
	p.emitText()
	p.push(interpreter.Message{Type: "group-start"})
	p.groups++
}

func (p *Parser) emitEndGroup() {
	p.emitText()
	p.push(interpreter.Message{Type: "group-end"})
	p.groups--
}

func (p *Parser) emitFormula() {
	panic("not implemented")
}

func (p *Parser) emitIndexSubEntry() {
	panic("not implemented")
}

func (p *Parser) emitEndParagraph() {
	p.emitText()
	p.push(interpreter.Message{Type: "paragraph-end"})
}

func (p *Parser) emitError(message string, err error) {
	p.err = err
	p.emitText()
	p.push(interpreter.Message{Type: "error", Value: message})
}

func (p *Parser) emitKeyword() {
	p.emitText()
	if p.keyWord == "" {
		p.emitError("empty keyword", nil)
		p.keyWordParam = ""
		return
	}
	var value int
	if p.keyWordParam != "" {
		var err error
		value, err = strconv.Atoi(p.keyWordParam)
		if err != nil {
			p.emitError("invalid keyword param", err)
			p.keyWordParam = ""
			p.keyWord = ""
			return
		}
	}
	p.push(interpreter.Message{Type: "keyword", Value: p.keyWord, Param: value})
	p.keyWord = ""
	p.keyWordParam = ""
}
