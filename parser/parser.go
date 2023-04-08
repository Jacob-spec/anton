package parser

import (
	"strings"
)

func ParseScreenplay(lexemes []Lexeme) Screenplay {
	var screenplay Screenplay
	var scene Scene
	screenplay.Metadata, lexemes = parseScreenplayMetadata(lexemes)

	for {
		scene, lexemes = ParseScene(lexemes)
		scene.PrintScene()
		screenplay.Scenes = append(screenplay.Scenes, scene)
		if lexemes == nil {
			break
		}
	}

	return screenplay
}

func parseScreenplayMetadata(lexemes []Lexeme) ([]Metadata, []Lexeme) {
	var tildeIndices []int
	metadata := make([]Metadata, 0)

	for i, lexeme := range lexemes {
		if lexeme.Typ == Tilde {
			tildeIndices = append(tildeIndices, i)
		}
	}

	for i := 0; (i + 1) < len(tildeIndices); i += 1 {
		// eliminates problem of metadata entries only having a newline between them
		if len(lexemes[tildeIndices[i]:tildeIndices[i+1]]) > 1 {
			metadata = append(metadata, parseMetadataPair(lexemes[tildeIndices[i]:tildeIndices[i+1]]))
		} else {
			continue
		}
	}
	// + 1 is to consume the closing tilde and return only the lexemes of the actual screenplay
	return metadata, lexemes[tildeIndices[len(tildeIndices)-1]+1:]
}

// parseScreenplayMetadata takes care of closing tilde bc this func
// only returns the Metadata struct, not that plus remaining lexemes
func parseMetadataPair(lexemes []Lexeme) Metadata {
	var data Metadata
	// consume opening tilde
	lexemes = lexemes[1:]

	assertLexemeTypes(lexemes, []LexemeType{Text, Colon, Text})
	data.Key = lexemes[0].Value
	data.Value = lexemes[2].Value
	return data

}

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

// returns the int/ext attribute and the rest of the text lexeme
func parseIntExtKeyword(lex Lexeme) (IntExt, Lexeme) {
	s := strings.ToUpper(lex.Value)
	// Seperates the int/ext portion from the rest of scene header
	strWords := strings.Split(s, " ")
	// checks if it's an int/ext scene
	intExt := strings.Split(strWords[0], "/")
	s = strings.Join(strWords[1:], " ")

	if !strings.Contains(intExt[0], "INT") && !strings.Contains(intExt[0], "EXT") {
		column := len(lex.Value) - len(s)
		Throw(ExpectingFoundErr("INT., EXT., or INT/EXT", intExt[0], lex.LineNumber, column))
	}

	if len(intExt) > 1 {
		return INTEXT, Lexeme{Text, s, lex.LineNumber, lex.ColumnNumber}
	} else if strings.Contains(intExt[0], "INT") {
		return INT, Lexeme{Text, s, lex.LineNumber, lex.ColumnNumber}
	} else {
		return EXT, Lexeme{Text, s, lex.LineNumber, lex.ColumnNumber}
	}
}

func ParseScene(lexemes []Lexeme) (Scene, []Lexeme) {
	var scene Scene
	var remainingLexemes []Lexeme
	scene.Heading, remainingLexemes = parseSceneHeading(lexemes)
	for i, lexeme := range remainingLexemes {
		if lexeme.Typ == VertBar {
			scene.Contents = parseSceneContents(remainingLexemes[0:i])
			remainingLexemes = remainingLexemes[i:]
			break
		} else if lexeme.Typ == EOF {
			scene.Contents = parseSceneContents(remainingLexemes[0:i])
			remainingLexemes = nil
		}
	}

	return scene, remainingLexemes
}

func parseSceneContents(lexemes []Lexeme) []SceneItem {
	items := make([]SceneItem, 0)

	for {
		if len(lexemes) <= 0 {
			break
		}
		switch lexemes[0].Typ {
		case Equals:
			var dUnit DialogueUnit
			dUnit, lexemes = parseDialogueUnit(lexemes)
			items = append(items, dUnit)
		case LSquareBracket:
			var action Action
			action, lexemes = parseAction(lexemes)
			items = append(items, action)
		case Dash:
			var shot Shot
			shot, lexemes = parseShot(lexemes)
			items = append(items, shot)
		case Plus:
			var transition Transition
			transition, lexemes = parseTransition(lexemes)
			items = append(items, transition)
		default:
			Throw(SyntaxErr("Character, Action, Shot, Transition, etc.", lexemes[0].LineNumber, lexemes[0].ColumnNumber))
		}
	}

	return items
}

func parseDialogueUnit(lexemes []Lexeme) (DialogueUnit, []Lexeme) {
	var dUnit DialogueUnit
	// assert then consume equals then character name
	assertLexemeTypes(lexemes[:2], []LexemeType{Equals, Text})
	dUnit.Character.Name = lexemes[1].Value
	lexemes = lexemes[2:]
	if lexemes[0].Typ == LParenthesis {
		dUnit.HasParenthetical = true
		assertLexemeTypes(lexemes[1:], []LexemeType{Text, RParenthesis, LCurlyBracket, Text, RCurlyBracket})
		dUnit.Parenthetical.Value = lexemes[1].Value
		dUnit.Dialogue.Value = lexemes[4].Value
		return dUnit, lexemes[6:]
	} else if lexemes[0].Typ == LCurlyBracket {
		dUnit.HasParenthetical = false
		dUnit.Parenthetical.Value = ""
		assertLexemeTypes(lexemes[1:], []LexemeType{Text, RCurlyBracket})
		dUnit.Dialogue.Value = lexemes[1].Value
		return dUnit, lexemes[3:]
	}

	return dUnit, nil
}

func parseAction(lexemes []Lexeme) (Action, []Lexeme) {
	assertLexemeTypes(lexemes, []LexemeType{LSquareBracket, Text, RSquareBracket})
	action := Action{Value: lexemes[1].Value}
	return action, lexemes[3:]
}

func parseShot(lexemes []Lexeme) (Shot, []Lexeme) {
	assertLexemeTypes(lexemes, []LexemeType{Dash, Text, Dash})
	shot := Shot{Value: lexemes[1].Value}
	return shot, lexemes[3:]
}

func parseTransition(lexemes []Lexeme) (Transition, []Lexeme) {
	assertLexemeTypes(lexemes, []LexemeType{Plus, Text, Plus})
	transition := Transition{Value: lexemes[1].Value}
	return transition, lexemes[3:]
}

func Parse(contents string) []Lexeme {
	lexemes := lex(contents)
	s := ParseScreenplay(lexemes)
	PrintMetadata(s.Metadata)

	return lexemes
}
