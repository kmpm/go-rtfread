package rtfread

import (
	"bufio"
	"context"
	"log/slog"
	"os"

	"github.com/kmpm/go-rtfread/interpreter"
	"github.com/kmpm/go-rtfread/parser"
)

func ParseFile(filename string) (string, error) {
	fp, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer fp.Close()
	return Parse(bufio.NewReader(fp))
}

func Parse(r *bufio.Reader) (string, error) {
	ipr, err := interpreter.New()
	if err != nil {
		return "", err
	}
	p, err := parser.New(ipr)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	slog.Info("parse")
	go p.Parse(ctx, r)
	slog.Info("wait for done")
	<-p.Done()
	return ipr.Value(), nil
}
