package main

import (
	"assessment-tool-cli/parser"
	"assessment-tool-cli/tui"
)

func main() {
	tui.StartTea(parser.DecodeTOML())
}
