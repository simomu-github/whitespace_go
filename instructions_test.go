package whitespace_go

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newExecutor() *Executor {
	return &Executor{
		heap:           map[int]int{},
		programCounter: 0,
	}
}

func TestPush(t *testing.T) {
	executor := newExecutor()

	push := Push{value: 1}
	push.Execute(executor)

	assert.Equal(t, executor.stack, []int{1})
}

func TestDuplicate(t *testing.T) {
	executor := newExecutor()
	executor.stack = []int{1}

	duplicate := Duplicate{}
	duplicate.Execute(executor)

	assert.Equal(t, executor.stack, []int{1, 1})
}

func TestSwap(t *testing.T) {
	executor := newExecutor()
	executor.stack = []int{1, 2}

	swap := Swap{}
	swap.Execute(executor)

	assert.Equal(t, executor.stack, []int{2, 1})
}

func TestDiscard(t *testing.T) {
	executor := newExecutor()
	executor.stack = []int{1}

	discard := Discard{}
	discard.Execute(executor)

	assert.Equal(t, executor.stack, []int{})
}

func TestAddition(t *testing.T) {
	executor := newExecutor()
	executor.stack = []int{1, 2}

	addition := Addition{}
	addition.Execute(executor)

	assert.Equal(t, executor.stack, []int{3})
}

func TestSubtraction(t *testing.T) {
	executor := newExecutor()
	executor.stack = []int{2, 1}

	subtraction := Subtraction{}
	subtraction.Execute(executor)

	assert.Equal(t, executor.stack, []int{1})
}

func TestMultiplication(t *testing.T) {
	executor := newExecutor()
	executor.stack = []int{2, 2}

	multiplication := Multiplication{}
	multiplication.Execute(executor)

	assert.Equal(t, executor.stack, []int{4})
}

func TestDivision(t *testing.T) {
	executor := newExecutor()
	executor.stack = []int{2, 2}

	divition := Division{}
	divition.Execute(executor)

	assert.Equal(t, executor.stack, []int{1})
}

func TestModulo(t *testing.T) {
	executor := newExecutor()
	executor.stack = []int{5, 3}

	modulo := Modulo{}
	modulo.Execute(executor)

	assert.Equal(t, executor.stack, []int{2})
}

func TestStore(t *testing.T) {
	executor := newExecutor()
	executor.stack = []int{1, 2}

	store := Store{}
	store.Execute(executor)

	assert.Equal(t, executor.heap, map[int]int{1: 2})
}

func TestRetrieve(t *testing.T) {
	executor := newExecutor()
	executor.stack = []int{1}
	executor.heap = map[int]int{1: 2}

	retrieve := Retrieve{}
	retrieve.Execute(executor)

	assert.Equal(t, executor.stack, []int{2})
}

func TestCallSubroutine(t *testing.T) {
	executor := newExecutor()
	executor.instructions = append(executor.instructions, MarkLabel{label: TAB})
	executor.programCounter = 1

	callSubroutine := CallSubroutine{label: TAB}
	callSubroutine.Execute(executor)

	assert.Equal(t, executor.programCounter, 0)
}

func TestEndSubroutine(t *testing.T) {
	executor := newExecutor()
	executor.programCounter = 10
	executor.callStack = append(executor.callStack, 0)

	endSubroutine := EndSubroutine{}
	endSubroutine.Execute(executor)

	assert.Equal(t, executor.programCounter, 0)
}

func TestJumpLabel(t *testing.T) {
	executor := newExecutor()
	executor.instructions = append(executor.instructions, MarkLabel{label: TAB})
	executor.programCounter = 1

	jumpLabel := JumpLabel{label: TAB}
	jumpLabel.Execute(executor)

	assert.Equal(t, executor.programCounter, 0)
}

func TestJumplabelWhenZero(t *testing.T) {
	executor := newExecutor()
	executor.instructions = append(executor.instructions, MarkLabel{label: TAB})
	executor.programCounter = 1
	executor.stack = []int{0}

	jumpLabelWhenZero := JumpLabelWhenZero{label: TAB}
	jumpLabelWhenZero.Execute(executor)

	assert.Equal(t, executor.programCounter, 0)

	executor.programCounter = 1
	executor.stack = []int{1}
	jumpLabelWhenZero.Execute(executor)

	assert.Equal(t, executor.programCounter, 1)
}

func TestJumplabelWhenNegative(t *testing.T) {
	executor := newExecutor()
	executor.instructions = append(executor.instructions, MarkLabel{label: TAB})
	executor.programCounter = 1
	executor.stack = []int{-1}

	jumpLabelWhenNegative := JumpLabelWhenNegative{label: TAB}
	jumpLabelWhenNegative.Execute(executor)

	assert.Equal(t, executor.programCounter, 0)

	executor.programCounter = 1
	executor.stack = []int{0}
	jumpLabelWhenNegative.Execute(executor)

	assert.Equal(t, executor.programCounter, 1)

	executor.programCounter = 1
	executor.stack = []int{1}
	jumpLabelWhenNegative.Execute(executor)

	assert.Equal(t, executor.programCounter, 1)
}

func TestEndProgram(t *testing.T) {
	executor := newExecutor()
	executor.instructions = append(executor.instructions, MarkLabel{label: TAB})
	executor.programCounter = 0

	endProgram := EndProgram{}
	endProgram.Execute(executor)

	assert.Equal(t, executor.programCounter, 1)
}
