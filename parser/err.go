package parser

import (
	"fmt"
	"os"
)

type ParseErr struct {
	Message      string
	LineNumber   int
	ColumnNumber int
}

func Throw(err ParseErr) {
	fmt.Printf("Err!: %s  (%d,%d)\n", err.Message, err.LineNumber, err.ColumnNumber)
	os.Exit(1)
}

func assertLexemeType(lexeme Lexeme, assertedType LexemeType) {
	var expected, found string
	if lexeme.Typ != assertedType {
		if lexeme.Typ != Text {
			expected = fmt.Sprintf("'%s'", assertedType.String())
			found = fmt.Sprintf("'%s'", lexeme.Value)
		} else {
			expected = fmt.Sprintf("Value of type '%s'", assertedType.String())
			found = fmt.Sprintf("'%s'", lexeme.Value)
		}

		Throw(ExpectingFoundErr(expected, found, lexeme.LineNumber, lexeme.ColumnNumber))
	}
}

func assertLexemeTypes(lexemes []Lexeme, assertedTypes []LexemeType) {
	for i := 0; i < len(assertedTypes); i += 1 {
		assertLexemeType(lexemes[i], assertedTypes[i])
	}
}

func ExpectingFoundErr(expecting string, found string, line int, column int) ParseErr {
	var err ParseErr
	err.Message = fmt.Sprintf("Expecting %s; Found %s", expecting, found)
	err.LineNumber = line
	err.ColumnNumber = column
	return err
}

func SyntaxErr(expecting string, line int, column int) ParseErr {
	var err ParseErr
	err.Message = fmt.Sprintf("Expecting %s", expecting)
	err.LineNumber = line
	err.ColumnNumber = column
	return err
}
