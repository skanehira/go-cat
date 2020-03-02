package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

var printNumber = flag.Bool("n", false, "print number")
var enableColor = flag.Bool("c", false, "enable ansi color")

func init() {
	flag.Parse()
}

func main() {
	for _, f := range flag.Args() {
		var str string
		var err error

		str, err = readFile(f)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if *enableColor {
			ext := filepath.Ext(f)
			str = highlight(str, ext)
		}

		println(str)
	}
}

func readFile(file string) (string, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func highlight(data string, ext string) string {
	l := lexers.Get(ext)
	if l == nil {
		l = lexers.Analyse(data)
	}
	if l == nil {
		l = lexers.Fallback
	}
	l = chroma.Coalesce(l)

	f := formatters.Get("terminal256")
	if f == nil {
		f = formatters.Fallback
	}

	// Determine style.
	s := styles.Get("monokai")
	if s == nil {
		s = styles.Fallback
	}

	it, err := l.Tokenise(nil, data)
	if err != nil {
		log.Println(err)
		return err.Error()
	}

	var buf bytes.Buffer

	if err := f.Format(&buf, s, it); err != nil {
		log.Println(err)
		return err.Error()
	}

	return buf.String()
}
