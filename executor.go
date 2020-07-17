package whitespace_go

type Executor struct {
	instructions   []Instruction
	stack          []int
	heap           map[int]int
	programCounter int
	callStack      []int
}

func (executor *Executor) Run() error {
	executor.heap = map[int]int{}
	executor.programCounter = 0

	for executor.programCounter = 0; executor.programCounter < len(executor.instructions); executor.programCounter++ {
		executor.instructions[executor.programCounter].Execute(executor)
	}

	return nil
}

func (executor *Executor) Push(value int) {
	executor.stack = append(executor.stack, value)
}

func (executor *Executor) Pop() (int, error) {
	if len(executor.stack) == 0 {
		return 0, runtimeError(executor, "stack is epmty")
	}

	value := executor.stack[len(executor.stack)-1]
	executor.stack = executor.stack[:len(executor.stack)-1]
	return value, nil
}

func (executor *Executor) PushCallStack(counter int) {
	executor.callStack = append(executor.callStack, counter)
}

func (executor *Executor) PopCallStack() (int, error) {
	if len(executor.callStack) == 0 {
		return 0, runtimeError(executor, "call stack is empry")
	}

	counter := executor.callStack[len(executor.callStack)-1]
	executor.callStack = executor.callStack[:len(executor.callStack)-1]
	return counter, nil
}
