package parser

import "fmt"

// Types go from highest order to lowest order
type Screenplay struct {
	Scenes   []Scene
	Metadata []Metadata
}

type Metadata struct {
	Key   string
	Value string
}

type Scene struct {
	Heading  SceneHeading
	Contents []SceneItem
}

type SceneHeading struct {
	IntOrExt IntExt
	Location string
	HasTime  bool
	Time     string
}

type IntExt int

const (
	INT = iota
	EXT
	INTEXT
)

type SceneItem interface {
	Print()
}

type DialogueUnit struct {
	Character        Character
	Dialogue         Dialogue
	HasParenthetical bool
	Parenthetical    Parenthetical
}

type Character struct {
	Name     string
	Metadata []Metadata
}

type Parenthetical struct {
	Value string
}

type Dialogue struct {
	Value string
}

type Action struct {
	Value string
}

type Shot struct {
	Value string
}

type Transition struct {
	Value string
}

type ParseErr struct {
	Message    string
	LineNumber int
}

// debugging
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

// debugging
func (s *SceneHeading) String() string {
	return fmt.Sprintf("%s\n%s\n%t\n%s\n", s.IntOrExt.String(), s.Location, s.HasTime, s.Time)
}

// debugging
func (s *Scene) PrintScene() {
	fmt.Printf("Scene:\n%s\n\n", s.Heading.String())
	for _, sceneItem := range s.Contents {
		sceneItem.Print()
	}
}

// debugging
func (s Metadata) String() string {
	return fmt.Sprintf("Key(%s) | Value(%s)", s.Key, s.Value)
}

// debugging
func PrintMetadata(metadata []Metadata) {
	for _, pair := range metadata {
		fmt.Println(pair.String())
	}
}

func (d DialogueUnit) Print() {
	fmt.Println("DU {")
	d.Character.Print()
	if d.HasParenthetical {
		d.Parenthetical.Print()
	}
	d.Dialogue.Print()
	fmt.Printf("}\n\n")
}

func (c Character) Print() {
	fmt.Printf("Character(%s)\n", c.Name)
}

func (p Parenthetical) Print() {
	fmt.Printf("Parenthetical(%s)\n", p.Value)
}

func (d Dialogue) Print() {
	fmt.Printf("Dialogue(%s)\n", d.Value)
}

func (a Action) Print() {
	fmt.Printf("Action(%s)\n", a.Value)
}

func (t Transition) Print() {
	fmt.Printf("Transition(%s)\n", t.Value)
}

func (s Shot) Print() {
	fmt.Printf("Shot(%s)\n", s.Value)
}
