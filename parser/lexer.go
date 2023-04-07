package parser

import (
	"fmt"
	"strings"
)

type LexemeType int

const (
	Text    LexemeType = iota
	VertBar            // vertical bar character
	LParenthesis
	RParenthesis
	LSquareBracket
	RSquareBracket
	LCurlyBracket
	RCurlyBracket
	Equals
	Dash
	Plus
	Colon
	Tilde
	Null
	Err
	EOF
)

type lexError struct {
	message string
}

type Lexeme struct {
	Typ          LexemeType
	Value        string
	LineNumber   int
	ColumnNumber int
}
type lexer struct {
	input        string // the text being lexed
	start        int    // start position of current Lexeme
	pos          int    // current position in the input
	lineNumber   int
	columnNumber int
	lexemes      []Lexeme
}

func PrintLexeme(l Lexeme) {
	switch l.Typ {
	case Text:
		fmt.Printf("Text(%s, %d, %d)\n", l.Value, l.LineNumber, l.ColumnNumber)
	case Err:
		fmt.Printf("Err(%s)\n", l.Value)
	case EOF:
		fmt.Println("EOF")
	case Null:
		fmt.Println("Null")
	case LParenthesis:
		fmt.Printf("LParenthesis, %d, %d\n", l.LineNumber, l.ColumnNumber)
	case RParenthesis:
		fmt.Printf("RParenthesis, %d, %d\n", l.LineNumber, l.ColumnNumber)
	case LSquareBracket:
		fmt.Printf("LSquareBracket, %d, %d\n", l.LineNumber, l.ColumnNumber)
	case RSquareBracket:
		fmt.Printf("RSquareBracket, %d, %d\n", l.LineNumber, l.ColumnNumber)
	case LCurlyBracket:
		fmt.Printf("LCurlyBracket, %d, %d\n", l.LineNumber, l.ColumnNumber)
	case RCurlyBracket:
		fmt.Printf("RCurlyBracket, %d, %d\n", l.LineNumber, l.ColumnNumber)
	case VertBar:
		fmt.Printf("VertBar, %d, %d\n", l.LineNumber, l.ColumnNumber)
	case Equals:
		fmt.Printf("Equals, %d, %d\n", l.LineNumber, l.ColumnNumber)
	case Dash:
		fmt.Printf("Dash, %d, %d\n", l.LineNumber, l.ColumnNumber)
	case Plus:
		fmt.Printf("Plus, %d, %d\n", l.LineNumber, l.ColumnNumber)
	case Colon:
		fmt.Printf("Colon, %d, %d\n", l.LineNumber, l.ColumnNumber)
	case Tilde:
		fmt.Printf("Tilde, %d, %d\n", l.LineNumber, l.ColumnNumber)
	default:
		fmt.Printf("%s\n", l.Value)
	}
}

func PrintLexemes(lexemes []Lexeme) {
	for _, lex := range lexemes {
		PrintLexeme(lex)
	}
}

func (l *lexer) emit(t LexemeType) {
	value := l.input[l.start:l.pos]
	l.columnNumber += len(value)

	for _, char := range value {
		if char == '\n' {
			l.lineNumber += 1
			l.columnNumber = 1
		}
	}

	l.lexemes = append(l.lexemes, Lexeme{t, value, l.lineNumber, l.columnNumber})
	l.start = l.pos
	l.pos += 1
}

func (l *lexer) current() (byte, *lexError) {
	if l.pos >= len(l.input) {
		return 0, &lexError{message: "Index out of range"}
	} else {
		return l.input[l.pos], nil
	}
}

func lex(input string) []Lexeme {
	l := lexer{
		input:        input,
		lexemes:      make([]Lexeme, 0),
		lineNumber:   1,
		columnNumber: 1,
		pos:          0,
		start:        0,
	}

	l.lexText()
	cleanedLex := cleanLexemes(l.lexemes)
	PrintLexemes(cleanedLex)

	return cleanedLex
}

func cleanLexemes(l []Lexeme) []Lexeme {
	clean := make([]Lexeme, 0)
	for _, lex := range l {
		if lex.Typ != Text {
			clean = append(clean, lex)
		} else if len(strings.TrimSpace(lex.Value)) > 1 {
			lex.Value = strings.TrimSpace(lex.Value)
			clean = append(clean, lex)
		}
	}
	clean = append(clean, Lexeme{Typ: EOF, Value: "EOF"})
	return clean
}

func (l *lexer) lexText() {

	for l.pos < len(l.input) {
		current, err := l.current()
		if err != nil {
			l.emit(EOF)
			break
		}

		char := isSpecialCharacter(current)
		if char != Err {
			if l.pos+1 > l.start {
				l.emit(Text)
				l.emit(char)
			} else {
				l.emit(char)
			}

		} else {
			l.pos += 1
		}
	}
}

func isSpecialCharacter(character byte) LexemeType {
	switch character {
	case '|':
		return VertBar
	case '(':
		return LParenthesis
	case ')':
		return RParenthesis
	case '[':
		return LSquareBracket
	case ']':
		return RSquareBracket
	case '{':
		return LCurlyBracket
	case '}':
		return RCurlyBracket
	case '=':
		return Equals
	case '-':
		return Dash
	case '+':
		return Plus
	case ':':
		return Colon
	case '~':
		return Tilde
	default:
		return Err
	}
}
