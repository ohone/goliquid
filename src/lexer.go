package templater

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const openDelimiter = "{{"
const closeDelimiter = "}}"

type lexeme struct {
	templatable bool
	token       string
}

type stateFn func(*lexer) stateFn
type lexer struct {
	name  string
	input string      // string being scanned
	start int         // start position of this item
	pos   int         // current position in the input
	width int         // width of the last run read
	items chan lexeme // channel of scanned items
	state stateFn
}

func lexLeftMeta(l *lexer) stateFn {
	l.pos += len(openDelimiter)
	l.emit(false)
	return lexInsideTemplate // {{}}
}

func lexRightMeta(l *lexer) stateFn {
	l.pos += len(closeDelimiter)
	l.emit(false)
	return lexText
}

func lexText(l *lexer) stateFn {
	for { // loop
		if strings.HasPrefix(l.input[l.pos:], openDelimiter) { // if we're starting a token
			if l.pos > l.start { // emit previous tokens as plain text
				l.emit(false)
			}
			return lexLeftMeta // return next state
		}
		l.next()
	}
}

func lexInsideTemplate(l *lexer) stateFn {
	for { // loop
		if strings.HasPrefix(l.input[l.pos:], closeDelimiter) {
			if l.pos > l.start {
				l.emit(true)
			}

			if !l.accept("abcdefhijklmnopqrstuvwxyzABCDEFHIJKLMNOPQRSTUVWXYZ") {
				return l.errorf("object to template must be alphanumeric")
			}
			l.acceptRun("abcdefhijklmnopqrstuvwxyzABCDEFHIJKLMNOPQRSTUVWXYZ")

			if !strings.HasPrefix(l.input[l.pos:], closeDelimiter) {
				return l.errorf("object template must finish with closing delimiter `}}`")
			}

			return lexRightMeta
		}
	}
}

func (l *lexer) next() (rune rune) {
	if l.pos == len(l.input) {
		l.width = 0
		return '_'
	}
	rune, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return rune
}

// returns whether the next character is in the charset
// TODO: modify to take a boolean function instead of charset
//		 let the caller decide acceptance criteria
func (l *lexer) accept(charset string) bool {
	if strings.IndexRune(charset, l.next()) > 0 {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRun(charset string) {
	for {
		if !l.accept(charset) {
			return
		}
	}
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	// TODO: return an error token rather than printing to console
	// requires that emit takes a type rather than a bool
	println(fmt.Sprintf(format, args...))
	return nil
}

// skip current char set
func (l *lexer) ignore() {
	l.start = l.pos
}

// go back one rune
// should revert a next
func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) peek() rune {
	rune := l.next()
	l.backup()
	return rune
}

func (l *lexer) emit(template bool) {
	l.items <- lexeme{template, l.input[l.start:l.pos]} // send token to parser
	l.start = l.pos                                     // update position
}

func (l *lexer) NextLexeme() lexeme {
	for {
		select {
		case lexeme := <-l.items: // if item can be recieved from channel (will halt here if nothing to recieve)
			return lexeme // return item from channel, deliver to caller
		default: // if item can't be recieved from channel, do lex iteration (may generate token)
			l.state = l.state(l)
		}
	}
}

// TODO: take reader instead of string
func Lex(name, input string) *lexer {
	l := &lexer{
		name:  name,
		input: input,
		items: make(chan lexeme, 2), // refactor - no need to use channel, by-hand ring buffer would be better. this is no longer goroutine-y so channel is overhead when a linear datastructure would do the same job
		state: lexText,
	}

	return l
}
