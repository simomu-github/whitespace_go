package main

import (
	"github.com/simomu-github/whitespace_go"
	"os"
)

func main() {
	interpreter := whitespace_go.New()
	os.Exit(interpreter.Run())
}
