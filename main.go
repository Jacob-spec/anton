package main

import (
	"github.com/Jacob-spec/anton/cli"
	"github.com/Jacob-spec/anton/parser"
)

func main() {
	fileContents := cli.GetScreenplayContents("/Users/jacobstoner/Code/Go/anton/misc/script.an")
	parser.Parse(fileContents)
}
