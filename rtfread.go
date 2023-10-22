package rtfread

import (
	"bufio"
	"io"
	"os"
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
	p, err := parse(r)
	if err != nil && err != io.EOF {
		return "", err
	}
	return p.String(), nil
}
