package whitespace_go

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func newParser(code string) Parser {
	return Parser{
		sourceCode:   []rune(code),
		currentIndex: -1,
	}
}

func TestParsePositiveNumber(t *testing.T) {
	positiveNumberCode := SPACE + TAB + SPACE + LF
	parser := newParser(positiveNumberCode)

	n, err := parser.parseNumber()
	if err != nil {
		t.Errorf("expected parseNumber to return 2, but can not parse number %s", err.Error())
	}

	if n != 2 {
		t.Errorf("expected parseNumber to return 2, but got %d", n)
	}
}

func TestParseNegativeNumberCode(t *testing.T) {
	negativeNumberCode := TAB + TAB + SPACE + LF
	parser := newParser(negativeNumberCode)

	n, err := parser.parseNumber()
	if err != nil {
		t.Errorf("expected parseNumber to return -2, but can not parse number %s", err.Error())
	}

	if n != -2 {
		t.Errorf("expected parseNumber to return -2, but got %d", n)
	}
}

func TestParseStackManipulationInstruction(t *testing.T) {
	sources := []string{
		SPACE + SPACE + SPACE + SPACE + SPACE + LF,
		SPACE + LF + SPACE,
		SPACE + LF + TAB,
		SPACE + LF + LF,
	}

	expectedInstructions := []Instruction{
		Push{value: 0},
		Duplicate{},
		Swap{},
		Discard{},
	}

	for i, s := range sources {
		expected := expectedInstructions[i]
		parser := newParser(s)

		instructions, err := parser.parse()
		if err != nil {
			t.Errorf("expected parse stack manipulate instruction, but raise error %s", err.Error())
		}

		assert.Equal(t, instructions[0], expected)
	}
}

func TestParseArithmeticInstrucition(t *testing.T) {
	sources := []string{
		TAB + SPACE + SPACE + SPACE,
		TAB + SPACE + SPACE + TAB,
		TAB + SPACE + SPACE + LF,
		TAB + SPACE + TAB + SPACE,
		TAB + SPACE + TAB + TAB,
	}

	expectedInstructions := []Instruction{
		Addition{},
		Subtraction{},
		Multiplication{},
		Division{},
		Modulo{},
	}

	for i, s := range sources {
		expected := expectedInstructions[i]
		parser := newParser(s)

		instructions, err := parser.parse()
		if err != nil {
			t.Errorf("expected parse artithemetic instructions, but raise error %s", err.Error())
		}

		assert.Equal(t, instructions[0], expected)
	}
}

func TestParseIOInstrucition(t *testing.T) {
	sources := []string{
		TAB + LF + TAB + SPACE,
		TAB + LF + TAB + TAB,
		TAB + LF + SPACE + SPACE,
		TAB + LF + SPACE + TAB,
	}

	expectedInstructions := []Instruction{
		Getc{},
		Getn{},
		Putc{},
		Putn{},
	}

	for i, s := range sources {
		expected := expectedInstructions[i]
		parser := newParser(s)

		instructions, err := parser.parse()
		if err != nil {
			t.Errorf("expected parse I/O instructions, but raise error %s", err.Error())
		}

		assert.Equal(t, instructions[0], expected)
	}
}

func TestParseFlowControllInstruction(t *testing.T) {
	sources := []string{
		LF + SPACE + SPACE + TAB + LF,
		LF + SPACE + TAB + TAB + LF,
		LF + SPACE + LF + TAB + LF,
		LF + TAB + SPACE + TAB + LF,
		LF + TAB + TAB + TAB + LF,
		LF + TAB + LF,
		LF + LF + LF,
	}

	expectedInstructions := []Instruction{
		MarkLabel{label: TAB},
		CallSubroutine{label: TAB},
		JumpLabel{label: TAB},
		JumpLabelWhenZero{label: TAB},
		JumpLabelWhenNegative{label: TAB},
		EndSubroutine{},
		EndProgram{},
	}

	for i, s := range sources {
		expected := expectedInstructions[i]
		parser := newParser(s)

		instructions, err := parser.parse()
		if err != nil {
			t.Errorf("expected parse flow controll instructions, but raise error %s", err.Error())
		}

		assert.Equal(t, instructions[0], expected)
	}
}

func TestParseHeapAccessInstruction(t *testing.T) {
	sources := []string{
		TAB + TAB + SPACE,
		TAB + TAB + TAB,
	}

	expectedInstructions := []Instruction{
		Store{},
		Retrieve{},
	}

	for i, s := range sources {
		expected := expectedInstructions[i]
		parser := newParser(s)

		instructions, err := parser.parse()
		if err != nil {
			t.Errorf("expected parse flow controll instructions, but raise error %s", err.Error())
		}

		assert.Equal(t, instructions[0], expected)
	}
}
