package rtfread

import (
	"bufio"
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/kmpm/go-rtfread/interpreter"
	"github.com/kmpm/go-rtfread/parser"
)

var debug bool = strings.Contains(os.Getenv("DEBUG"), "rtfread")

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
	if debug {
		slog.Info("parse")
	}
	go p.Parse(ctx, r)
	if debug {
		slog.Info("wait for done")
	}
	<-p.Done()
	return ipr.Value(), nil
}
