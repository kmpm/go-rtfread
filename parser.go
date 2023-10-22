package rtfread

import (
	"bufio"
	"bytes"
	"fmt"
	"log/slog"
)

type stateParser func(ch byte) error

type msg struct {
	Type  string
	Value string
	Param string
}

type parser struct {
	parserState  stateParser
	pos          int
	text         *bytes.Buffer
	groups       int
	hexChar      string
	keyWord      string
	keyWordParam string
}

func parse(r *bufio.Reader) (*parser, error) {
	p := &parser{
		text: &bytes.Buffer{},
	}
	p.parserState = p.parseText

	for {
		ch, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		p.pos++
		if err := p.parserState(ch); err != nil {
			return nil, err
		}
	}

	return p, nil
}

func (p *parser) String() string {

	return ""
}

func (p *parser) writeByte(ch byte) error {
	return p.text.WriteByte(ch)
}

func (p *parser) writeRune(r rune) error {
	_, err := p.text.WriteRune(r)
	return err
}

func (p *parser) writeString(s string) error {
	_, err := p.text.WriteString(s)
	return err
}

func (p *parser) push(m msg) {
	slog.Info("push", "pos", p.pos, "group", p.groups, "type", m.Type, "value", m.Value, "param", m.Param)
}

func (p *parser) parseText(ch byte) error {
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

func (p *parser) parseEscape(ch byte) error {
	if ch == '\\' || ch == '{' || ch == '}' {
		p.writeByte(ch)
		p.parserState = p.parseText
		return nil
	}
	p.parserState = p.parseSymbol
	return p.parseSymbol(ch)
}

func (p *parser) parseSymbol(ch byte) error {
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
	case '\r', '\n':
		p.emitEndParagraph()
		p.parserState = p.parseText
	default:
		p.parserState = p.parseKeyword
		return p.parseKeyword(ch)
	}
	return nil
}

func (p *parser) parseHex(ch byte) error {

	if !ishex(ch) {
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

func (p *parser) parseKeyword(ch byte) error {

	if ch == ' ' {
		p.emitKeyword()
		p.parserState = p.parseText
	} else if ch == '-' || isdigit(ch) {
		p.parserState = p.parseKeywordParam
		p.keyWordParam += string(ch)
	} else if isalpha(ch) {
		p.keyWord += string(ch)
	} else {
		p.emitKeyword()
		p.parserState = p.parseText
		return p.parseText(ch)
	}

	return nil
}

func (p *parser) parseKeywordParam(ch byte) error {
	if isdigit(ch) {
		p.keyWordParam += string(ch)
	} else if ch == ' ' {
		p.emitKeyword()
		p.parserState = p.parseText
	} else {
		p.emitKeyword()
		p.parserState = p.parseText
		return p.parserState(ch)
	}
	return nil
}

func (p *parser) emitText() {
	if p.text.Len() > 0 {
		p.push(msg{
			Type:  "text",
			Value: p.text.String(),
		})
		p.text.Reset()
	}
}

func (p *parser) emitIgnorable() {
	p.emitText()
	p.push(msg{Type: "ignorable"})
}

func (p *parser) emitHexChar() {
	p.emitText()
	// v, err := strconv.ParseInt(p.hexChar, 16, 32)
	// if err != nil {
	// 	return err
	// }
	p.push(msg{Type: "hex-char", Value: p.hexChar})
	p.hexChar = ""
}

func (p *parser) emitStartGroup() {
	p.emitText()
	p.push(msg{Type: "group-start"})
	p.groups++
}

func (p *parser) emitEndGroup() {
	p.emitText()
	p.push(msg{Type: "group-end"})
	p.groups--
}

func (p *parser) emitFormula() {
	panic("not implemented")
}

func (p *parser) emitIndexSubEntry() {
	panic("not implemented")
}

func (p *parser) emitEndParagraph() {
	p.emitText()
	p.push(msg{Type: "paragraph-end"})
}

func (p *parser) emitError(message string, err error) {
	p.emitText()
	p.push(msg{Type: "error", Value: message})
}

func (p *parser) emitKeyword() {
	p.emitText()
	if p.keyWord == "" {
		p.emitError("empty keyword", nil)
		p.keyWordParam = ""
		return
	}
	p.push(msg{"keyword", p.keyWord, p.keyWordParam})
	p.keyWord = ""
	p.keyWordParam = ""
}
