package parser

import (
	"fmt"
	"os"
)

func LexAll(filename string) {
	dat, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Read file")
	}
	lexemes := lex(string(dat))
	for _, lex := range lexemes {
		PrintLexeme(lex)
	}

}
