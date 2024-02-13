package collection

type Stack []string

// IsEmpty: check if stack is empty
func (stack *Stack) IsEmpty() bool {
	return len(*stack) == 0
}

// Push a new value onto the stack
func (stack *Stack) Push(str string) {
	*stack = append(*stack, str) //Simply append the new value to the end of the stack
}

// Remove top element of stack. Return false if stack is empty.
func (stack *Stack) Pop() bool {
	if stack.IsEmpty() {
		return false
	} else {
		index := len(*stack) - 1  // Get the index of top most element.
		*stack = (*stack)[:index] // Remove it from the stack by slicing it off.
		return true
	}
}

// Return top element of stack. Return false if stack is empty.
func (stack *Stack) Top() string {
	if stack.IsEmpty() {
		return ""
	} else {
		index := len(*stack) - 1   // Get the index of top most element.
		element := (*stack)[index] // Index onto the slice and obtain the element.
		return element
	}
}
