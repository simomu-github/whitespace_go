package whitespace_go

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

var (
	versionOpt = flag.Bool("v", false, "display version information")
)

const version = "v0.0.1"

type Interpreter struct {
	args     []string
	stderr   io.Writer
	parser   Parser
	executor Executor
}

func New() *Interpreter {
	return &Interpreter{
		args:   os.Args,
		stderr: os.Stderr,
	}
}

func (i *Interpreter) Run() int {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n  ws [FILE]\n", os.Args[0])
		flag.PrintDefaults()
	}

	if len(i.args) < 2 {
		flag.Usage()
		return 1
	}

	flag.Parse()
	if *versionOpt {
		fmt.Printf("ws version %s\n", version)
		return 1
	}

	filename := i.args[1]

	bytes, errReadFile := ioutil.ReadFile(filename)
	if errReadFile != nil {
		fmt.Fprintf(i.stderr, "%s can not read\n", filename)
		return 1
	}

	i.parser = Parser{filename: filename, rawSourceCode: string(bytes)}
	errParse := i.parser.parseAll()
	if errParse != nil {
		fmt.Fprintln(i.stderr, errParse.Error())
		return 1
	}

	i.executor = Executor{instructions: i.parser.instructions}
	errRuntime := i.executor.Run()
	if errRuntime != nil {
		fmt.Fprintln(i.stderr, errRuntime.Error())
		return 1
	}

	return 0
}
