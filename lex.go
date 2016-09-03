package gitolite_validator

import (
	"io"
	"bufio"
	"fmt"
	"unicode"
)

const eof = -1

type tokenType int

func (t tokenType) String() string {
	switch t {
	case t_error:
		return "t_error"
	case t_eof:
		return "t_eof"
	case t_eol:
		return "t_eol"
	case t_comment:
		return "t_comment"
	case t_group:
		return "t_group"
	case t_assign:
		return "t_assign"
	case t_string:
		return "t_string"
	case t_dash:
		return "t_dash"
	default:
		panic(fmt.Sprintf("unknown token type: %v", t))
	}
}

const (
	t_error tokenType = iota
	t_eof
	t_eol
	t_comment                                // a comment
	t_group
	t_assign
	t_string
	t_dash
)

type stateFn func(*lexer) stateFn

type token struct {
	t tokenType
	s string
}

func (t token) String() string {
	return fmt.Sprintf("{%s %s}", t.t, t.s)
}

type lexer struct {
	in     io.RuneReader
	out    chan token
	buf    []rune // running buffer for current lexeme
	backup []rune
	err    error
}

func lex(r io.Reader) chan token {
	l := lexer{
		in:     bufio.NewReader(r),
		out:    make(chan token),
		backup: make([]rune, 0, 4),
	}
	go l.lex()
	return l.out
}

func (l *lexer) lex() {
	defer close(l.out)
	for fn := lexRoot; fn != nil; {
		fn = fn(l)
		if l.err != nil {
			fn = lexErrorf("read error: %s", l.err)
		}
	}
}

func (l *lexer) next() rune {
	if len(l.backup) > 0 {
		r := l.backup[len(l.backup) - 1]
		l.backup = l.backup[:len(l.backup) - 1]
		return r
	}
	r, _, err := l.in.ReadRune()
	switch err {
	case io.EOF:
		return eof
	case nil:
		return r
	default:
		l.err = err
		return eof
	}
	return r
}

func (l *lexer) keep(r rune) {
	if l.buf == nil {
		l.buf = make([]rune, 0, 18)
	}
	l.buf = append(l.buf, r)
}

func (l *lexer) emit(t tokenType) {
	l.out <- token{t, string(l.buf)}
	l.buf = l.buf[0:0]
}

func lexRoot(l *lexer) stateFn {
	r := l.next()
	switch {
	case r == eof:
		return nil
	case r == '#':
		return lexComment
	case r == '@':
		return lexGroup
	case r == '-':
		l.emit(t_dash)
		return lexRoot
	case r == '=':
		l.emit(t_assign)
		return lexRoot
	case r == '\n':
		l.emit(t_eol)
		return lexRoot
	case unicode.IsSpace(r) && r != '\n':
		return lexRoot
	case unicode.IsPrint(r):
		l.keep(r)
		return lexNameOrString
	default:
		return lexErrorf("unexpected rune in lexRoot: %c", r)
	}
}

func (l *lexer) bufHasSpaces() bool {
	for _, r := range l.buf {
		if unicode.IsSpace(r) {
			return true
		}
	}
	return false
}

func lexErrorf(t string, args ...interface{}) stateFn {
	return func(l *lexer) stateFn {
		l.out <- token{t_error, fmt.Sprintf(t, args...)}
		return nil
	}
}

func lexNameOrString(l *lexer) stateFn {
	r := l.next()
	switch {
	case r == '\n':
		l.emit(t_string)
		l.emit(t_eol)
		return lexRoot
	case unicode.IsSpace(r), r == eof:
		l.emit(t_string)
		return lexRoot
	case unicode.IsPrint(r):
		l.keep(r)
		return lexNameOrString
	default:
		return lexErrorf("unhandled character type in lexDuration: %c", r)
	}
}

func lexGroup(l *lexer) stateFn {
	r := l.next();
	switch {
	case r == '\n':
		l.emit(t_group)
		l.emit(t_eol)
		return lexRoot
	case unicode.IsSpace(r), r == eof:
		l.emit(t_group)
		return lexRoot
	case unicode.IsPrint(r):
		l.keep(r)
		return lexGroup
	default:
		return lexErrorf("unhandled character type in lexDuration: %c", r)
	}
}

func lexComment(l *lexer) stateFn {
	switch r := l.next(); r {
	case '\n':
		l.emit(t_comment)
		l.emit(t_eol)
		return lexRoot
	case eof:
		l.emit(t_comment)
		return nil
	default:
		l.keep(r)
		return lexComment
	}
}