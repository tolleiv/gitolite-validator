package gitolite_validator

import (
	"io"
)

type parser struct {
	root   node
	input  chan token
	backup []token
}

func Read(r io.Reader) (error) {
	_, err := parse(r)
	return err
}

func (p *parser) parse() error {
	if p.root == nil {
		p.root = newRootNode()
	}
	return p.root.parse(p)
}

func parse(r io.Reader) (node, error) {
	p := &parser{
		root:   newRootNode(),
		input:  lex(r),
		backup: make([]token, 0, 8),
	}
	if err := p.parse(); err != nil {
		return nil, err
	}
	return p.root, nil
}
func (p *parser) next() token {
	if len(p.backup) > 0 {
		t := p.backup[len(p.backup) - 1]
		p.backup = p.backup[:len(p.backup) - 1]
		return t
	}
	SKIP_COMMENTS:
	t, ok := <-p.input
	if !ok {
		return token{t_eof, "eof"}
	}
	if t.t == t_comment {
		goto SKIP_COMMENTS
	}
	return t
}

func (p *parser) unread(t token) {
	if p.backup == nil {
		p.backup = make([]token, 0, 8)
	}
	p.backup = append(p.backup, t)
}

func (p *parser) peek() token {
	t := p.next()
	p.unread(t)
	return t
}
