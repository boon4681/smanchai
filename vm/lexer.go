package smanchai

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
)

type Range struct {
	Line   int
	Column int
	Index  int
}

func (r *Range) String() string {
	return fmt.Sprintf("%d:%d", r.Line, r.Column)
}

type Lexer struct {
	line   int
	column int
	index  int
	reader *Reader
	Dg     bool
	buf    struct {
		r     Range
		token Token
		str   string
	}
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		line:   1,
		column: 0,
		index:  0,
		Dg:     false,
		reader: NewReader(bufio.NewReader(reader)),
	}
}

func (l *Lexer) save() Range {
	return Range{Line: l.line, Column: l.column, Index: l.index}
}

func (l *Lexer) load(r Range) {
	for i := 0; i < l.index-r.Index; i++ {
		l.back()
	}
	for i := 0; i < r.Index-l.index; i++ {
		l.next()
	}
	l.line = r.Line
	l.column = r.Column
	l.index = r.Index
}

func (l *Lexer) next() (r rune, size int, err error) {
	c, s, err := l.reader.ReadRune()
	if err == nil {
		l.column++
		l.index++
	}
	return c, s, err
}

func (l *Lexer) back() {
	if err := l.reader.UnreadRune(); err != nil {
		panic(err)
	}
	if l.index > 0 {
		l.index--
	}
	if l.column > 0 {
		l.column--
	}
}

func (l *Lexer) Lex() (Range, Token, string) {
	r, token, str := l.lex()
	l.buf.r = r
	l.buf.token = token
	l.buf.str = str
	if l.Dg {
		fmt.Printf("Lx: \t%s\t%s\t%s\n", r.String(), token.String(), str)
	}
	return r, token, str
}

func (l *Lexer) lex() (Range, Token, string) {
	defer l.reader.CleanUp()
	c, _, err := l.next()
	if err != nil {
		if err == io.EOF {
			return l.save(), EOF, ""
		}
		panic(err)
	}
	if c == '\n' {
		l.line++
		l.column = 0
	}
	if c == '"' {
		return l.lexString()
	}
	if c == '@' {
		return l.save(), AT, "@"
	}
	if c == '+' {
		return l.save(), ADD, "+"
	}
	if c == '-' {
		return l.save(), SUB, "-"
	}
	if l.isKeyword(c, "**") {
		return l.save(), POW, "**"
	}
	if c == '*' {
		return l.save(), MULT, "*"
	}
	if c == '/' {
		return l.save(), DIV, "/"
	}
	if c == '(' {
		return l.save(), LParent, "("
	}
	if c == ')' {
		return l.save(), RParent, ")"
	}
	if c == ',' {
		return l.save(), COMMA, ","
	}
	if c == '.' {
		return l.save(), DOT, "."
	}
	if l.isKeyword(c, "==") {
		return l.save(), EQUALITY_OPERATOR, "=="
	}
	if l.isKeyword(c, "!=") {
		return l.save(), EQUALITY_OPERATOR, "!="
	}
	if c == '>' {
		return l.save(), COMPARISON_OPERATOR, ">"
	}
	if c == '<' {
		return l.save(), COMPARISON_OPERATOR, "<"
	}
	if l.isKeyword(c, ">=") {
		return l.save(), COMPARISON_OPERATOR, ">="
	}
	if l.isKeyword(c, "<=") {
		return l.save(), COMPARISON_OPERATOR, "<="
	}
	if l.isKeyword(c, "and") {
		return l.save(), CONJUNCTION, "and"
	}
	if l.isKeyword(c, "or") {
		return l.save(), DISJUNCTION, "or"
	}
	if l.isKeyword(c, "true") {
		return l.save(), BOOL, "true"
	}
	if l.isKeyword(c, "false") {
		return l.save(), BOOL, "false"
	}
	if unicode.IsDigit(c) {
		l.back()
		return l.lexNumber()
	}
	if unicode.IsSpace(c) {
		l.back()
		return l.lexWhiteSpace()
	}
	if unicode.IsLetter(c) {
		l.back()
		return l.lexIdentifier()
	}
	return l.save(), ILLEGAL, string(c)
}

func (l *Lexer) lexIdentifier() (Range, Token, string) {
	result := ""
	r := l.save()
	for {
		c, _, err := l.next()
		if err != nil {
			if err == io.EOF {
				return r, IDENTIFIER, result
			}
			panic(err)
		}
		if unicode.IsLetter(c) {
			result += string(c)
		} else {
			l.back()
			return r, IDENTIFIER, result
		}
	}
}

func (l *Lexer) lexWhiteSpace() (Range, Token, string) {
	result := ""
	r := l.save()
	for {
		c, _, err := l.next()
		if err != nil {
			if err == io.EOF {
				return r, WS, result
			}
			panic(err)
		}
		if unicode.IsSpace(c) {
			result += string(c)
		} else {
			l.back()
			return r, WS, result
		}
	}
}

func (l *Lexer) lexNumber() (Range, Token, string) {
	result := ""
	float := false
	zc := 0
	r := l.save()
	for {
		c, _, err := l.next()
		if err != nil {
			if err == io.EOF {
				return r, NUMBER, result
			}
			panic(err)
		}
		if unicode.IsDigit(c) {
			if !float && c == '0' {
				zc++
			}
			if !float && zc > 1 {
				panic("WTF Number begin with double 0??")
			}
			result += string(c)
		} else if c == '.' {
			if float {
				panic("WTF Number with 2 dot what are you doing????????")
			}
			float = true
			result += string(c)
		} else {
			l.back()
			return r, NUMBER, result
		}
	}
}

func (l *Lexer) lexString() (Range, Token, string) {
	result := ""
	escape := false
	r := l.save()
	for {
		c, _, err := l.next()
		if err != nil {
			if err == io.EOF {
				panic("String must be closed with \"")
			}
			panic(err)
		}
		if c == '\\' {
			escape = !escape
		}
		if c == '"' {
			if len(result) >= 0 && escape {
				escape = false
				result += string(c)
				continue
			}
			return r, STRING, result
		}
		if c == '\n' {
			panic("String must be closed with \"")
		}
		result += string(c)
	}
}

func (l *Lexer) isKeyword(c rune, base string) bool {
	target := []rune(base)
	if c != target[0] {
		return false
	}
	r := l.save()
	l.back()
	for _, char := range target {
		c, _, err := l.next()
		if err != nil {
			if err == io.EOF {
				l.load(r)
				return false
			}
			panic(err)
		}
		if c != char {
			l.load(r)
			return false
		}
	}
	return true
}
