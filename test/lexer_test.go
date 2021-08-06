package main

import (
	"fmt"
	"testing"

	"github.com/ohone/goliquid/lexer"
)

func TestLex(t *testing.T) {
	b := lexer.Lex("name", "stringlol")
	eme := b.NextLexeme()
	fmt.Printf(eme.Token)
	if eme.Templatable {
		t.Error()
	}
}
func TestLexTemplateable(t *testing.T) {
	b := lexer.Lex("name", "{{hello}}")
	eme := b.NextLexeme()
	fmt.Printf(eme.Token)
	eme2 := b.NextLexeme()
	fmt.Printf(eme2.Token)
	eme3 := b.NextLexeme()
	fmt.Printf(eme3.Token)

	if eme.Templatable {
		t.Error()
	}
}
