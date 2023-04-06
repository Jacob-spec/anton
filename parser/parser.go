package parser

import (
	"fmt"
	"os"
	"strings"
)

type IntExt int

const (
	INT = iota
	EXT
	INTEXT
)

type SceneHeading struct {
	IntOrExt IntExt
	Location string
	HasTime  bool
	Time     string
}

func (s *SceneHeading) String() string {
	return fmt.Sprintf("%d\n%s\n%t\n%s\n", s.IntOrExt, s.Location, s.HasTime, s.Time)
}

// type Scene struct {
// 	Heading  SceneHeading
// 	Contents []Lexeme
// }

// type Screenplay struct {
// 	Scenes []Scene
// 	Author string
// 	Year   int
// 	Genre  string
// }

// func ParseScreenplay(lexemes []Lexeme) {}

func ParseSceneHeading(lexemes []Lexeme) (SceneHeading, []Lexeme) {
	scene := SceneHeading{}
	var remainingLexemes []Lexeme
	// consume opening pipe symbol
	lexemes = lexemes[1:]
	for i, lex := range lexemes {
		if lex.Typ == Dash {
			scene.HasTime = true
		} else if lex.Typ == VertBar {
			remainingLexemes = lexemes[0:(i + 1)]
			break
		}
	}
	scene.IntOrExt, remainingLexemes[0] = parseIntExtKeyword(remainingLexemes[0])
	scene.Location = remainingLexemes[0].Value
	// if the header has a time, we know there's a dash in between the text [0] and the time
	if scene.HasTime {
		scene.Time = remainingLexemes[2].Value
		return scene, remainingLexemes[4:]
	} else {
		scene.Time = ""
		return scene, remainingLexemes[2:]
	}

}

// returns the int/ext attributre and the rest of the text lexeme
func parseIntExtKeyword(lex Lexeme) (IntExt, Lexeme) {
	s := strings.ToUpper(lex.Value)
	// Seperates the int/ext portion from the rest of scene header
	strWords := strings.Split(s, " ")
	// checks if it's an int/ext scene
	intExt := strings.Split(strWords[0], "/")
	s = strings.Join(strWords[1:], " ")

	if len(intExt) > 1 {
		return INTEXT, Lexeme{Text, s}
	} else if strings.Contains(intExt[0], "INT") {
		return INT, Lexeme{Text, s}
	} else {
		return EXT, Lexeme{Text, s}
	}
}

func LexAll(filename string) []Lexeme {
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
	heading, _ := ParseSceneHeading(lexemes)
	fmt.Printf("%s", heading)

	return lexemes

}
