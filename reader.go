package smanchai

import (
	"bufio"
	"fmt"
	"unicode/utf8"
)

type Reader struct {
	mark   int
	index  int
	pad    int
	unread int
	buffer []rune
	reader *bufio.Reader
}

func NewReader(reader *bufio.Reader) *Reader {
	return &Reader{
		mark:   0,
		index:  0,
		pad:    0,
		unread: 0,
		buffer: make([]rune, 0, 256),
		reader: reader,
	}
}

func (r *Reader) ReadRune() (rune, int, error) {
	defer func() {
		r.unread = 0
	}()
	if r.index+r.pad < r.mark {
		c, s := r.buffer[r.index], utf8.RuneLen(r.buffer[r.index])
		r.index++
		return c, s, nil
	}
	c, s, err := r.reader.ReadRune()
	if err == nil {
		r.buffer = append(r.buffer, c)
		r.index++
		r.mark++
	}
	return c, s, err
}

func (r *Reader) UnreadRune() error {
	if r.index == 0 {
		return bufio.ErrInvalidUnreadRune
	}
	r.index--
	r.unread++
	return nil
}

func (r *Reader) CleanUp() {
	if r.unread > 0 {
		r.buffer = r.buffer[len(r.buffer)-r.unread:]
		r.pad += r.index
		r.index = 0
	}
}

func (r *Reader) Debug() {
	fmt.Printf("Reader {\n    buf: %v,\n    i: %d\n}\n", r.buffer, r.index)
}
