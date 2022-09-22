package whitespace_go

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Instruction interface {
	Execute(executor *Executor) error
}

type Push struct {
	value int
}

func (p Push) Execute(executor *Executor) error {
	executor.stack = append(executor.stack, p.value)

	return nil
}

type Swap struct{}

func (s Swap) Execute(executor *Executor) error {
	a, errA := executor.Pop()
	if errA != nil {
		return errA
	}

	b, errB := executor.Pop()
	if errB != nil {
		return errB
	}

	executor.Push(a)
	executor.Push(b)
	return nil
}

type Duplicate struct{}

func (d Duplicate) Execute(executor *Executor) error {
	a, err := executor.Pop()
	if err != nil {
		return err
	}

	executor.Push(a)
	executor.Push(a)
	return nil
}

type Discard struct{}

func (d Discard) Execute(executor *Executor) error {
	executor.Pop()

	return nil
}

type Addition struct{}

func (a Addition) Execute(executor *Executor) error {
	rhs, errRhs := executor.Pop()
	if errRhs != nil {
		return errRhs
	}

	lhs, errLhs := executor.Pop()
	if errLhs != nil {
		return errLhs
	}

	executor.Push(lhs + rhs)

	return nil
}

type Subtraction struct{}

func (s Subtraction) Execute(executor *Executor) error {
	rhs, errRhs := executor.Pop()
	if errRhs != nil {
		return errRhs
	}

	lhs, errLhs := executor.Pop()
	if errLhs != nil {
		return errLhs
	}

	executor.Push(lhs - rhs)

	return nil
}

type Multiplication struct{}

func (m Multiplication) Execute(executor *Executor) error {
	rhs, errRhs := executor.Pop()
	if errRhs != nil {
		return errRhs
	}

	lhs, errLhs := executor.Pop()
	if errLhs != nil {
		return errLhs
	}

	executor.Push(lhs * rhs)

	return nil
}

type Division struct{}

func (d Division) Execute(executor *Executor) error {
	rhs, errRhs := executor.Pop()
	if errRhs != nil {
		return errRhs
	}

	lhs, errLhs := executor.Pop()
	if errLhs != nil {
		return errLhs
	}

	executor.Push(lhs / rhs)

	return nil
}

type Modulo struct{}

func (m Modulo) Execute(executor *Executor) error {
	rhs, errRhs := executor.Pop()
	if errRhs != nil {
		return errRhs
	}

	lhs, errLhs := executor.Pop()
	if errLhs != nil {
		return errLhs
	}

	executor.Push(lhs % rhs)

	return nil
}

type Getc struct{}

func (g Getc) Execute(executor *Executor) error {

	stdin := bufio.NewScanner(os.Stdin)
	stdin.Scan()
	text := stdin.Text()

	if len(text) == 0 {
		return runtimeError(executor, "input is empty")
	}

	address, err := executor.Pop()
	if err != nil {
		return err
	}

	executor.heap[address] = int([]rune(text)[0])
	return nil
}

type Getn struct{}

func (g Getn) Execute(executor *Executor) error {
	stdin := bufio.NewScanner(os.Stdin)
	stdin.Scan()
	text := stdin.Text()
	n, err := strconv.Atoi(text)
	if err != nil {
		return runtimeError(executor, "input character is not numeric")
	}

	address, err := executor.Pop()
	if err != nil {
		return err
	}

	executor.heap[address] = n
	return nil
}

type Putc struct{}

func (p Putc) Execute(executor *Executor) error {
	n, err := executor.Pop()
	if err != nil {
		return err
	}

	fmt.Printf("%c", n)

	return nil
}

type Putn struct{}

func (p Putn) Execute(executor *Executor) error {
	n, err := executor.Pop()
	if err != nil {
		return err
	}

	fmt.Printf("%d", n)

	return nil
}

type Store struct{}

func (s Store) Execute(executor *Executor) error {
	value, errValue := executor.Pop()
	if errValue != nil {
		return errValue
	}

	address, errAddress := executor.Pop()
	if errAddress != nil {
		return errAddress
	}

	executor.heap[address] = value
	return nil
}

type Retrieve struct{}

func (r Retrieve) Execute(executor *Executor) error {
	address, err := executor.Pop()
	if err != nil {
		return err
	}

	value, ok := executor.heap[address]
	if !ok {
		return runtimeError(executor, "invalid heap access")
	}

	executor.Push(value)
	return nil
}

type MarkLabel struct {
	label string
}

func (m MarkLabel) Execute(executor *Executor) error {
	return nil
}

type CallSubroutine struct {
	label string
}

func (c CallSubroutine) Execute(executor *Executor) error {
	counter := executor.programCounter
	executor.PushCallStack(counter)

	for i, instruction := range executor.instructions {
		m, ok := instruction.(MarkLabel)
		if ok && m.label == c.label {
			executor.programCounter = i
			return nil
		}
	}

	return runtimeError(executor, "label not found")
}

type EndSubroutine struct{}

func (e EndSubroutine) Execute(executor *Executor) error {
	counter, err := executor.PopCallStack()
	if err != nil {
		return err
	}

	executor.programCounter = counter
	return nil
}

type JumpLabel struct {
	label string
}

func (j JumpLabel) Execute(executor *Executor) error {
	for i, instruction := range executor.instructions {
		m, ok := instruction.(MarkLabel)
		if ok && m.label == j.label {
			executor.programCounter = i
			return nil
		}
	}

	return runtimeError(executor, "label not found")
}

type JumpLabelWhenZero struct {
	label string
}

func (j JumpLabelWhenZero) Execute(executor *Executor) error {
	value, err := executor.Pop()
	if err != nil {
		return err
	}

	if value != 0 {
		return nil
	}

	for i, instruction := range executor.instructions {
		m, ok := instruction.(MarkLabel)
		if ok && m.label == j.label {
			executor.programCounter = i
			return nil
		}
	}

	return runtimeError(executor, "label not found")
}

type JumpLabelWhenNegative struct {
	label string
}

func (j JumpLabelWhenNegative) Execute(executor *Executor) error {
	value, err := executor.Pop()
	if err != nil {
		return err
	}

	if value >= 0 {
		return nil
	}

	for i, instruction := range executor.instructions {
		m, ok := instruction.(MarkLabel)
		if ok && m.label == j.label {
			executor.programCounter = i
			return nil
		}
	}

	return runtimeError(executor, "label not found")
}

type EndProgram struct{}

func (e EndProgram) Execute(executor *Executor) error {
	executor.programCounter = len(executor.instructions)
	return nil
}
