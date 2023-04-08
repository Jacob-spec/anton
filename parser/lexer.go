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

func (l LexemeType) String() string {
	switch l {
	case Text:
		return "Text"
	case Err:
		return "Err"
	case EOF:
		return "EOF"
	case Null:
		return "Null"
	case LParenthesis:
		return "("
	case RParenthesis:
		return ")"
	case LSquareBracket:
		return "["
	case RSquareBracket:
		return "]"
	case LCurlyBracket:
		return "{"
	case RCurlyBracket:
		return "}"
	case VertBar:
		return "|"
	case Equals:
		return "="
	case Dash:
		return "-"
	case Plus:
		return "+"
	case Colon:
		return ":"
	case Tilde:
		return "~"
	default:
		return "ERRRRR"
	}
}

func PrintLexemes(lexemes []Lexeme) {
	for _, lex := range lexemes {
		fmt.Println(lex.Typ.String())
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
