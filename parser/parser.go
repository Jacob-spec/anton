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

func (i IntExt) String() string {
	switch i {
	case INT:
		return "INT"
	case EXT:
		return "EXT"
	case INTEXT:
		return "INT/EXT"
	default:
		return "ERR"
	}
}

type SceneHeading struct {
	IntOrExt IntExt
	Location string
	HasTime  bool
	Time     string
}

func (s *SceneHeading) String() string {
	return fmt.Sprintf("%s\n%s\n%t\n%s\n", s.IntOrExt.String(), s.Location, s.HasTime, s.Time)
}

type Scene struct {
	Heading  SceneHeading
	Contents []Lexeme
}

func (s *Scene) PrintScene() {
	fmt.Printf("Scene:\n%s\n\n", s.Heading.String())
	PrintLexemes(s.Contents)
}

// type Screenplay struct {
// 	Scenes []Scene
// 	Author string
// 	Year   int
// 	Genre  string
// }

// func ParseScreenplay(lexemes []Lexeme) {}

func parseSceneHeading(lexemes []Lexeme) (SceneHeading, []Lexeme) {
	scene := SceneHeading{}
	var remainingLexemes []Lexeme
	var sceneLexemes []Lexeme
	// consume opening pipe symbol
	lexemes = lexemes[1:]
	for i, lex := range lexemes {
		if lex.Typ == Dash {
			scene.HasTime = true
		} else if lex.Typ == VertBar {
			sceneLexemes = lexemes[0:(i + 1)]
			remainingLexemes = lexemes[(i + 1):]
			break
		}
	}
	scene.IntOrExt, sceneLexemes[0] = parseIntExtKeyword(sceneLexemes[0])
	scene.Location = sceneLexemes[0].Value
	// if the header has a time, we know there's a dash in between the text [0] and the time
	if scene.HasTime {
		scene.Time = strings.ToUpper(sceneLexemes[2].Value)
		return scene, remainingLexemes
	} else {
		scene.Time = ""
		return scene, remainingLexemes
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

func ParseScene(lexemes []Lexeme) (Scene, []Lexeme) {
	var scene Scene
	var remainingLexemes []Lexeme
	scene.Heading, remainingLexemes = parseSceneHeading(lexemes)
	for i, lexeme := range remainingLexemes {
		if lexeme.Typ == VertBar {
			scene.Contents = remainingLexemes[0:i]
			remainingLexemes = remainingLexemes[i:]
			break
		} else if lexeme.Typ == EOF {
			scene.Contents = remainingLexemes[0:i]
			remainingLexemes = nil
		}
	}

	return scene, remainingLexemes
}

func LexAll(filename string) []Lexeme {
	dat, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Read file")
	}
	lexemes := lex(string(dat))
	// for _, lex := range lexemes {
	// 	PrintLexeme(lex)
	// }
	scene, _ := ParseScene(lexemes)
	scene.PrintScene()

	return lexemes

}
