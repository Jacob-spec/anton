package parser

import (
	"fmt"
	"strings"
)

type lexemeType int

const (
	Text    lexemeType = iota
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
	Null
	Err
	EOF
)

type lexError struct {
	message string
}

type lexeme struct {
	typ   lexemeType
	value string
}
type lexer struct {
	input   string // the text being lexed
	start   int    // start position of current lexeme
	pos     int    // current position in the input
	lexemes []lexeme
}

func PrintLexeme(l lexeme) {
	switch l.typ {
	case Text:
		fmt.Printf("Text(%s)\n", l.value)
	case Err:
		fmt.Printf("Err(%s)\n", l.value)
	case EOF:
		fmt.Println("EOF")
	case Null:
		fmt.Println("Null")
	case LParenthesis:
		fmt.Printf("LParenthesis(%s)\n", l.value)
	case RParenthesis:
		fmt.Printf("RParenthesis(%s)\n", l.value)
	case LSquareBracket:
		fmt.Printf("LSquareBracket(%s)\n", l.value)
	case RSquareBracket:
		fmt.Printf("RSquareBracket(%s)\n", l.value)
	case LCurlyBracket:
		fmt.Printf("LCurlyBracket(%s)\n", l.value)
	case RCurlyBracket:
		fmt.Printf("RCurlyBracket(%s)\n", l.value)
	case VertBar:
		fmt.Printf("VertBar(%s)\n", l.value)
	case Equals:
		fmt.Printf("Equals(%s)\n", l.value)
	case Dash:
		fmt.Printf("Dash(%s)\n", l.value)
	case Plus:
		fmt.Printf("Plus(%s)\n", l.value)
	default:
		fmt.Printf("%s\n", l.value)
	}
}

func (l *lexer) emit(t lexemeType) {
	value := l.input[l.start:l.pos]
	l.lexemes = append(l.lexemes, lexeme{t, value})
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

// currently lexText is doing the job of lex. Fix this.
func lex(input string) []lexeme {
	l := lexer{
		input:   input,
		lexemes: make([]lexeme, 0),
		pos:     0,
		start:   0,
	}

	l.lexText()
	cleanedLex := cleanLexemes(l.lexemes)

	return cleanedLex
}

func cleanLexemes(l []lexeme) []lexeme {
	clean := make([]lexeme, 0)
	for _, lex := range l {
		if lex.typ != Text {
			clean = append(clean, lex)
		} else if len(strings.TrimSpace(lex.value)) > 1 {
			lex.value = strings.TrimSpace(lex.value)
			clean = append(clean, lex)
		}
	}
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

func isSpecialCharacter(character byte) lexemeType {
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
	default:
		return Err
	}
}
