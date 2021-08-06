package main

import (
	"testing"

	"github.com/ohone/goliquid/lexer"
)

func TestLexesString(t *testing.T) {
	const strVal = "stringlol"
	b := lexer.Lex("name", strVal)
	eme := b.NextLexeme()
	if eme.Templatable {
		t.Error()
	}
	if eme.Token != strVal {
		t.Error()
	}
}

func TestLexesTemplateLexemesHaveCorrectTemplatability(t *testing.T) {
	b := lexer.Lex("name", "{{hello}}")
	eme := b.NextLexeme()
	if eme.Templatable {
		t.Error()
	}
	eme2 := b.NextLexeme()
	if !eme2.Templatable {
		t.Error()
	}
	eme3 := b.NextLexeme()
	if eme3.Templatable {
		t.Error()
	}
}

func TestLexesTemplateLexemesHaveCorrectTokenValue(t *testing.T) {
	b := lexer.Lex("name", "{{hello}}")
	eme := b.NextLexeme()
	if eme.Token != "{{" {
		t.Error("Expected first emitted token to be opening delimeter")
	}
	eme2 := b.NextLexeme()
	if eme2.Token != "hello" {
		t.Error("Expected second emitted token to be `hello`.")
	}
	eme3 := b.NextLexeme()
	if eme3.Token != "}}" {
		t.Error("Expected third emitted token to be closing delimeter.")
	}
}
