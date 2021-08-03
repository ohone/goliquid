package main

import (
	"testing"

	"github.com/ohone/goliquid/templater"
)

func TestLex(t *testing.T) {
	b := templater.Lex("name", "string")
	eme := b.NextLexeme()
	if eme.Templatable {
		t.Error()
	}
}
