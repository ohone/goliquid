package lexer

import (
	"strings"
	"unicode/utf8"
)

const openDelimeter = "{{"
const closeDelimeter = "}}"
const eof = '_'

type Lexeme struct {
	Templatable bool
	Token       string
	Error       bool
}

type stateFn func(*lexer) stateFn
type lexer struct {
	name  string
	input string      // string being scanned
	start int         // start position of this item
	pos   int         // current position in the input
	width int         // width of the last run read
	items chan Lexeme // channel of scanned items
	state stateFn
}

func lexLeftMeta(l *lexer) stateFn {
	l.pos += len(openDelimeter)
	l.emit(false)
	return lexInsideTemplate // {{}}
}

func lexRightMeta(l *lexer) stateFn {
	l.pos += len(closeDelimeter)
	l.emit(false)
	return lexText
}

func lexText(l *lexer) stateFn {
	for { // loop
		if strings.HasPrefix(l.input[l.pos:], openDelimeter) { // if we're starting a token
			if l.pos > l.start { // emit previous tokens as plain text
				l.emit(false)
			}
			return lexLeftMeta // return next state
		}
		if l.next() == eof {
			l.emit(false)
			break
		}
	}
	return nil
}

func lexInsideTemplate(l *lexer) stateFn {
	if strings.HasPrefix(l.input[l.pos:], closeDelimeter) {
		return lexRightMeta
	}

	// if first character in template isn't alphanumeric
	if !l.accept("abcdefhijklmnopqrstuvwxyzABCDEFHIJKLMNOPQRSTUVWXYZ") {
		return l.errorf("object to template must be alphanumeric")
	}
	// move cursor to the end of alphanumeric string
	l.acceptRun("abcdefhijklmnopqrstuvwxyzABCDEFHIJKLMNOPQRSTUVWXYZ")

	// if we haven't hit a close delimeter
	if !strings.HasPrefix(l.input[l.pos:], closeDelimeter) {
		l.emit(true)
		return l.errorf("object template must finish with closing delimeter `}}`")
	}

	// emit templateable token
	l.emit(true)

	// lex close delimeter
	return lexRightMeta
}

func (l *lexer) next() (rune rune) {
	if l.pos == len(l.input) {
		l.width = 0
		return eof
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

// Error token on the channel, nil function.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- Lexeme{
		Error: true,
	}
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
	l.items <- Lexeme{template, l.input[l.start:l.pos], false} // send token to parser
	l.start = l.pos                                            // update position
}

// Get the next lexeme from the text.
func (l *lexer) NextLexeme() Lexeme {
	for {
		select {
		case lexeme := <-l.items: // if item can be recieved from channel (will halt here if nothing to recieve)
			return lexeme // return item from channel, deliver to caller
		default: // if item can't be recieved from channel, do lex iteration (may generate token)
			l.state = l.state(l)
		}
	}
}

// Create a lexer for the input text.
// TODO: take reader instead of string
func Lex(name, input string) *lexer {
	l := &lexer{
		name:  name,
		input: input,
		items: make(chan Lexeme, 2), // refactor - no need to use channel, by-hand ring buffer would be better. this is no longer goroutine-y so channel is overhead when a linear datastructure would do the same job
		state: lexText,
	}

	return l
}
