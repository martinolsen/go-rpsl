package rpsl

import (
	"bufio"
	"io"
	"strings"
)

type Reader struct {
	l chan item
}

func NewReader(rd io.Reader) *Reader {
	return &Reader{l: lex(rd)}
}

func (r *Reader) Read() (*Object, error) {
	var object *Object
	var key string

	for item := range r.l {
		switch item.Type {
		case Value:
			if object == nil {
				object = &Object{
					Values: make(map[string][]string),
				}
			}
			object.Values[key] = append(object.Values[key], item.Text)
		case Key:
			key = strings.ToLower(item.Text)
			if object == nil {
				object = &Object{
					Values: make(map[string][]string),
				}
			}
			if object.Class == "" {
				object.Class = key
			}
		case EOR:
			if object != nil {
				return object, nil
			}
		}
	}

	if object != nil {
		return object, nil
	}

	return nil, io.EOF
}

type lexer struct {
	C    chan item
	buf  *bufio.Reader
	text []byte
}

type item struct {
	Type itemType
	Text string
}

type itemType int

const (
	Key itemType = iota
	Value
	EOR
)

type stateFn func(*lexer) stateFn

func lex(rd io.Reader) chan item {
	l := &lexer{
		buf: bufio.NewReader(rd),
		C:   make(chan item),
	}

	go func(l *lexer) {
		var state stateFn
		for state = lexStart; state != nil; state = state(l) {
		}
		close(l.C)
	}(l)

	return l.C
}

func (l *lexer) next() (byte, error) {
	c, err := l.buf.ReadByte()
	if err == nil {
		l.text = append(l.text, c)
	}
	return c, err
}

func (l *lexer) discard() {
	if ln := len(l.text); ln > 0 {
		l.text = l.text[:ln-1]
	}
}

func (l *lexer) emit(t itemType) {
	l.C <- item{
		Type: t,
		Text: string(l.text),
	}
	l.text = nil
}

func lexStart(l *lexer) stateFn {
	for {
		c, err := l.next()
		if err == io.EOF {
			return nil
		} else if err != nil {
			panic(err)
		}

		switch c {
		case ' ', '\t', '\n':
			l.discard()
			continue
		case '%', '#':
			l.discard()
			return lexComment
		default:
			return lexKey
		}
	}
}

func lexNewline(l *lexer) stateFn {
	for {
		c, err := l.next()
		if err == io.EOF {
			return nil
		} else if err != nil {
			panic(err)
		}

		switch c {
		case '\n':
			l.discard()
			l.emit(EOR)
			return lexStart
		case ' ', '\t':
			l.discard()
			return lexValue
		default:
			return lexKey
		}
	}
}

func lexComment(l *lexer) stateFn {
	for {
		c, err := l.next()
		if err == io.EOF {
			return nil
		} else if err != nil {
			panic(err)
		}

		l.discard()

		switch c {
		case '\n':
			return lexStart
		}
	}
}

func lexKey(l *lexer) stateFn {
	for {
		c, err := l.next()
		if err == io.EOF {
			return nil
		} else if err != nil {
			panic(err)
		}

		switch c {
		case ':':
			l.discard()
			l.emit(Key)
			return lexValue
		case '\n':
			l.discard()
			l.emit(Key)
			return lexNewline
		}
	}
}

func lexValue(l *lexer) stateFn {
	var started bool

	for {
		c, err := l.next()
		if err == io.EOF {
			return nil
		} else if err != nil {
			panic(err)
		}

		switch c {
		case ' ', '\t':
			if !started {
				l.discard()
			}
		case '\n':
			l.discard()
			l.emit(Value)
			return lexNewline
		default:
			started = true
		}
	}
}
