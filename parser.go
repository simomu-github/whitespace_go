package whitespace_go

import (
	"strings"
)

const (
	TAB   = "\t"
	LF    = "\n"
	SPACE = " "
)

var (
	Tokens = []string{TAB, LF, SPACE}
)

type Parser struct {
	filename      string
	rawSourceCode string
	sourceCode    []rune
	currentIndex  int
	currentLine   int
	currentColumn int
	newLine       bool
	instructions  []Instruction
}

func (parser *Parser) parseAll() error {
	parser.sourceCode = []rune(parser.rawSourceCode)
	parser.currentIndex = -1
	parser.currentLine = 1
	parser.currentColumn = 0
	parser.newLine = false

	instructions, err := parser.parse()
	parser.instructions = instructions

	return err
}

func (parser *Parser) parse() ([]Instruction, error) {
	token := parser.nextToken()
	switch string(token) {
	case SPACE:
		return parser.parseStackManipulation()
	case TAB:
		token = parser.nextToken()
		switch string(token) {
		case SPACE:
			return parser.parseArithmetic()
		case TAB:
			return parser.parseHeapAccess()
		case LF:
			return parser.parseIO()
		default:
			return parser.instructions, parseError(parser, "expected instruction modification parameters")
		}
	case LF:
		return parser.parseFlowControll()
	default:
		return parser.instructions, nil
	}
}

func (parser *Parser) parseStackManipulation() ([]Instruction, error) {
	token := parser.nextToken()
	switch string(token) {
	case SPACE:
		n, err := parser.parseNumber()
		if err != nil {
			return parser.instructions, err
		}

		return parser.addInstruction(Push{value: n})
	case LF:
		token = parser.nextToken()
		switch string(token) {
		case SPACE:
			return parser.addInstruction(Duplicate{})
		case TAB:
			return parser.addInstruction(Swap{})
		case LF:
			return parser.addInstruction(Discard{})
		default:
			return parser.instructions, parseError(parser, "expected stack manipulation command")
		}
	default:
		return parser.instructions, parseError(parser, "expected stack manipulation command")
	}
}

func (parser *Parser) parseHeapAccess() ([]Instruction, error) {
	token := parser.nextToken()
	switch string(token) {
	case SPACE:
		return parser.addInstruction(Store{})
	case TAB:
		return parser.addInstruction(Retrieve{})
	default:
		return parser.instructions, parseError(parser, "expected heap access command")
	}
}

func (parser *Parser) parseArithmetic() ([]Instruction, error) {
	token := parser.nextToken()
	switch string(token) {
	case SPACE:
		token = parser.nextToken()
		switch string(token) {
		case SPACE:
			return parser.addInstruction(Addition{})
		case TAB:
			return parser.addInstruction(Subtraction{})
		case LF:
			return parser.addInstruction(Multiplication{})
		default:
			return parser.instructions, parseError(parser, "expected artithemetic command")
		}
	case TAB:
		token = parser.nextToken()
		switch string(token) {
		case SPACE:
			return parser.addInstruction(Division{})
		case TAB:
			return parser.addInstruction(Modulo{})
		default:
			return parser.instructions, parseError(parser, "expected artithemetic command")
		}
	default:
		return parser.instructions, parseError(parser, "expected artithemetic command")
	}
}

func (parser *Parser) parseIO() ([]Instruction, error) {
	token := parser.nextToken()
	switch string(token) {
	case SPACE:
		token = parser.nextToken()
		switch string(token) {
		case SPACE:
			return parser.addInstruction(Putc{})
		case TAB:
			return parser.addInstruction(Putn{})
		default:
			return parser.instructions, parseError(parser, "expected I/O command")
		}
	case TAB:
		token = parser.nextToken()
		switch string(token) {
		case SPACE:
			return parser.addInstruction(Getc{})
		case TAB:
			return parser.addInstruction(Getn{})
		default:
			return parser.instructions, parseError(parser, "expected I/O command")
		}
	default:
		return parser.instructions, parseError(parser, "expected I/O command")
	}
}

func (parser *Parser) parseFlowControll() ([]Instruction, error) {
	token := parser.nextToken()
	switch string(token) {
	case SPACE:
		token = parser.nextToken()
		switch string(token) {
		case SPACE:
			label, err := parser.parseLabel()
			if err != nil {
				return parser.instructions, err
			}

			return parser.addInstruction(MarkLabel{label: label})
		case TAB:
			label, err := parser.parseLabel()
			if err != nil {
				return parser.instructions, err
			}

			return parser.addInstruction(CallSubroutine{label: label})
		case LF:
			label, err := parser.parseLabel()
			if err != nil {
				return parser.instructions, err
			}

			return parser.addInstruction(JumpLabel{label: label})
		default:
			return parser.instructions, parseError(parser, "expected flow controll command")
		}
	case TAB:
		token = parser.nextToken()
		switch string(token) {
		case SPACE:
			label, err := parser.parseLabel()
			if err != nil {
				return parser.instructions, err
			}

			return parser.addInstruction(JumpLabelWhenZero{label: label})
		case TAB:
			label, err := parser.parseLabel()
			if err != nil {
				return parser.instructions, err
			}

			return parser.addInstruction(JumpLabelWhenNegative{label: label})
		case LF:
			return parser.addInstruction(EndSubroutine{})
		default:
			return parser.instructions, parseError(parser, "expected flow controll command")
		}
	case LF:
		token = parser.nextToken()
		if string(token) == LF {
			return parser.addInstruction(EndProgram{})
		} else {
			return parser.instructions, parseError(parser, "expected flow controll command")
		}
	default:
		return parser.instructions, parseError(parser, "expected flow controll command")
	}
}

func (parser *Parser) parseNumber() (int, error) {
	token := parser.nextToken()
	var sign int

	switch string(token) {
	case SPACE:
		sign = 1
	case TAB:
		sign = -1
	default:
		return 0, parseError(parser, "expected sign")
	}

	n, err := parser.parseBinaryNumber(0, 0, parser.nextToken())

	return sign * n, err
}

func (parser *Parser) parseBinaryNumber(n int, digits int, token rune) (int, error) {
	switch string(token) {
	case SPACE:
		return parser.parseBinaryNumber(2*n, digits+1, parser.nextToken())
	case TAB:
		return parser.parseBinaryNumber(2*n+1, digits+1, parser.nextToken())
	case LF:
		if digits > 0 {
			return n, nil
		} else {
			return 0, parseError(parser, "expected number")
		}
	default:
		return 0, parseError(parser, "expected numeric parameters end with a linefeed")
	}
}

func (parser *Parser) parseLabel() (string, error) {
	var label string
	token := parser.nextToken()
	for string(token) == SPACE || string(token) == TAB {
		label = strings.Join([]string{label, string(token)}, "")
		token = parser.nextToken()
	}

	if string(token) != LF {
		return label, parseError(parser, "expected label parameters end with a linefeed")
	}

	if len([]rune(label)) == 0 {
		return label, parseError(parser, "invalid label")
	}

	return label, nil
}

func (parser *Parser) addInstruction(instruction Instruction) ([]Instruction, error) {
	parser.instructions = append(parser.instructions, instruction)
	return parser.parse()
}

func (parser *Parser) nextToken() rune {
	parser.currentIndex++

	if parser.newLine {
		parser.newLine = false
		parser.currentColumn = 0
	}

	parser.currentColumn++

	if parser.currentIndex >= len(parser.sourceCode) {
		return 0
	}

	for !contains(parser.currentToken()) {
		parser.currentIndex++
		if parser.currentIndex >= len(parser.sourceCode) {
			return 0
		}
	}

	if string(parser.currentToken()) == LF {
		parser.currentLine++
		parser.newLine = true
	}

	return parser.currentToken()
}

func (parser *Parser) currentToken() rune {
	return parser.sourceCode[parser.currentIndex]
}

func contains(r rune) bool {
	for _, token := range Tokens {
		if token == string(r) {
			return true
		}
	}
	return false
}
