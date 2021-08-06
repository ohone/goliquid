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
		t.Error("Raw text should not be templatable.")
	}
	if eme.Token != strVal {
		t.Error("Expected produced lexeme to be " + strVal + ", was " + eme.Token)
	}
}

func TestLexesTemplateLexemesHaveCorrectTemplatability(t *testing.T) {
	myLexer := lexer.Lex("name", "{{hello}}")
	eme := myLexer.NextLexeme()
	if eme.Templatable {
		t.Error("Opening delimeter shouldn't be templatable.")
	}
	eme2 := myLexer.NextLexeme()
	if !eme2.Templatable {
		t.Error("Contents of template delimeters should be templatable.")
	}
	eme3 := myLexer.NextLexeme()
	if eme3.Templatable {
		t.Error("Closing delimeter shouldn't be templatable.")
	}
}

func TestLexesTemplateLexemesHaveCorrectTokenValue(t *testing.T) {
	myLexer := lexer.Lex("name", "{{hello}}")
	eme := myLexer.NextLexeme()
	if eme.Token != "{{" {
		t.Error("Expected first emitted token to be opening delimeter")
	}
	eme2 := myLexer.NextLexeme()
	if eme2.Token != "hello" {
		t.Error("Expected second emitted token to be `hello`.")
	}
	eme3 := myLexer.NextLexeme()
	if eme3.Token != "}}" {
		t.Error("Expected third emitted token to be closing delimeter.")
	}
}
