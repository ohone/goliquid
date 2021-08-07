package main

import (
	"testing"

	"github.com/ohone/goliquid/lexer"
)

func TestLexesString(t *testing.T) {
	const strVal = "stringlol"
	b := lexer.Lex("name", strVal)
	eme, _ := b.NextLexeme()
	if eme.Type == lexer.ItemTemplatable {
		t.Error("Raw text should not be templatable.")
	}
	if eme.Token != strVal {
		t.Error("Expected produced lexeme to be " + strVal + ", was " + eme.Token)
	}
}

func TestLexesTemplateLexemesHaveCorrectTemplatability(t *testing.T) {
	myLexer := lexer.Lex("name", "{{hello}}")
	eme, _ := myLexer.NextLexeme()
	if eme.Type == lexer.ItemTemplatable {
		t.Error("Opening delimeter shouldn't be templatable.")
	}
	eme2, _ := myLexer.NextLexeme()
	if eme2.Type != lexer.ItemTemplatable {
		t.Error("Contents of template delimeters should be templatable.")
	}
	eme3, _ := myLexer.NextLexeme()
	if eme3.Type == lexer.ItemTemplatable {
		t.Error("Closing delimeter shouldn't be templatable.")
	}
}

func TestLexesTemplateLexemesHaveCorrectTokenValue(t *testing.T) {
	myLexer := lexer.Lex("name", "{{hello}}")
	eme, _ := myLexer.NextLexeme()
	if eme.Token != "{{" {
		t.Error("Expected first emitted token to be opening delimeter")
	}
	eme2, _ := myLexer.NextLexeme()
	if eme2.Token != "hello" {
		t.Error("Expected second emitted token to be `hello`.")
	}
	eme3, _ := myLexer.NextLexeme()
	if eme3.Token != "}}" {
		t.Error("Expected third emitted token to be closing delimeter.")
	}
}

func TestUnclosedTemplateEmitsErrorToken(t *testing.T) {
	myLexer := lexer.Lex("name", "{{hello")
	eme, _ := myLexer.NextLexeme()
	if eme.Token != "{{" {
		t.Error("Expected first emitted token to be opening delimeter")
	}
	eme2, _ := myLexer.NextLexeme()
	if eme2.Token != "hello" {
		t.Error("Expected second emitted token to be `hello`.")
	}
	eme3, _ := myLexer.NextLexeme()
	if eme3.Type != lexer.ItemError {
		t.Error("Expected third emitted token to be erroneous.")
	}
}

func TestErrEofReturnedAfterLexicalError(t *testing.T) {
	myLexer := lexer.Lex("name", "{{hello")
	myLexer.NextLexeme()
	myLexer.NextLexeme()
	eme3, _ := myLexer.NextLexeme()
	if eme3.Type != lexer.ItemError {
		t.Error("Expected third emitted token to be erroneous.")
	}
}

func TestErrEofReturnedOnEof(t *testing.T) {
	myLexer := lexer.Lex("name", "string")
	myLexer.NextLexeme()

	_, err := myLexer.NextLexeme()
	if err != lexer.ErrEof {
		t.Error("Expected ErrEof to be thrown when eof reached.")
	}
}

func TestNilPointerReturnedOnEof(t *testing.T) {
	myLexer := lexer.Lex("name", "string")
	myLexer.NextLexeme()

	ptr, _ := myLexer.NextLexeme()
	if ptr != nil {
		t.Error("Expected nil pointer for lexeme to be returned when eof.")
	}
}
