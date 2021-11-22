package lang

import (
	"errors"
	"io"
	"os"
	"strings"
)

const (
	LINE_END            byte = '\n'
	COMMENT_MULTI_BEGIN      = "#:"
	COMMENT_MULTI_END        = "##"
	COMMENT_SINGLE           = "#"
)

type Runner struct {
	interpreter *interpreter
	parser      *parser
	name        string
}

func NewRunner(name string) *Runner {
	p := &parser{}
	return &Runner{
		parser: p,
		interpreter: &interpreter{
			parser:   p,
			ctx:      ctxGlobal,
			labelMap: make(map[string][]*genStmt),
		},
		name: name,
	}
}

func (r *Runner) DoFile(fn string) error {
	fd, err := os.Open(fn)
	if err != nil {
		return r.errOf(err)
	}
	src, err := io.ReadAll(fd)
	if err != nil {
		return r.errOf(err)
	}
	ctx := ""
	line := ""
	isComment := false
	for _, b := range src {
		if b == LINE_END {
			line = strings.TrimSpace(ctx)
			ctx = ""
			if !isComment {
				if line == "" || strings.HasPrefix(line, COMMENT_SINGLE) {
					// do nothing lol
				} else if strings.HasPrefix(line, COMMENT_MULTI_BEGIN) {
					isComment = true
				} else {
					s, err := r.parser.ParseStatement(line)
					if err != nil {
						return r.errOf(err)
					}
					err = r.interpreter.RunStatement(s)
					if err != nil {
						return r.errOf(err)
					}
				}
			} else {
				if strings.HasSuffix(line, COMMENT_MULTI_END) {
					isComment = false
				}
			}
		} else {
			ctx += string(rune(b))
		}
	}
	return nil
}

func (r *Runner) errOf(err error) error {
	return errors.New("error in " + r.name + ":\n\t" + err.Error())
}
