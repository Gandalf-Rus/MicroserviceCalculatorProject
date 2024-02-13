package expression

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	c "MicroserviceCalculatorProject/orchestrator/pkg/collection"
)

// Function to return precedence of operators
func setPrecedence(s string) int {
	if (s == "/") || (s == "*") {
		return 2
	} else if (s == "+") || (s == "-") {
		return 1
	} else {
		return -1
	}
}

// If symbol is operator return true
func isOperator(s string) bool {
	return s == "+" || s == "-" || s == "*" || s == "/"
}

// Get subexpressions from AST (tree)
func splitIntoSubexpressions(root *c.Node, subexpressionMap map[int]string) (string, error) {
	if root == nil {
		return "", nil
	}

	leftExp, err1 := splitIntoSubexpressions(root.Left, subexpressionMap)
	rightExp, err2 := splitIntoSubexpressions(root.Right, subexpressionMap)
	if err1 != nil || err2 != nil {
		return "", errors.New("err")
	}

	if leftExp != "" && rightExp != "" {
		exp := fmt.Sprintf("%s%s%s", leftExp, root.Value, rightExp)
		for i := 1; i <= len(subexpressionMap); i++ {
			if subexpressionMap[i] == exp {
				return strconv.Itoa(i), nil
			}
		}
		subexpressionMap[len(subexpressionMap)+1] = exp
		return fmt.Sprintf("{%d}", len(subexpressionMap)), nil
	}

	if _, err := strconv.Atoi(root.Value); err == nil {
		return root.Value, nil
	}

	return "", errors.New("err")
}

// Build subexpression's tree
func buildAST(postfix []string) *c.Node {
	stack := make([]*c.Node, 0)

	for _, token := range postfix {
		strToken := string(token)
		node := &c.Node{Value: strToken}

		if isOperator(strToken) {
			node.Right = stack[len(stack)-1]
			node.Left = stack[len(stack)-2]
			stack = stack[:len(stack)-2]
		}

		stack = append(stack, node)
	}

	return stack[0]
}

// Convert classics infix expressions to postfix expression
func infixExpToPostfixExp(infixExp []string) ([]string, error) {
	var stack c.Stack
	var postfix []string

	for _, char := range infixExp {
		// if scanned character is operand, add it to output string
		if _, err := strconv.ParseFloat(char, 64); err == nil {
			postfix = append(postfix, char)
		} else if char == "(" {
			stack.Push(char)
		} else if char == ")" {
			for stack.Top() != "(" {
				postfix = append(postfix, stack.Top())
				stack.Pop()
			}
			stack.Pop()
		} else if isOperator(char) {
			for !stack.IsEmpty() && setPrecedence(char) <= setPrecedence(stack.Top()) {
				postfix = append(postfix, stack.Top())
				stack.Pop()
			}
			stack.Push(char)
		} else {
			return []string{}, errors.New("error")
		}
	}
	// Pop all the remaining elements from the stack
	for !stack.IsEmpty() {
		postfix = append(postfix, stack.Top())
		stack.Pop()
	}
	return postfix, nil
}

// Main function for get subexpressions
func ProcessExpression(infixExpression []string, subexpressionMap map[int]string) error {
	postfix, err := infixExpToPostfixExp(infixExpression)
	if err != nil {
		return errors.New("error in convertation from infix to postfix notation")
	}

	tree := buildAST(postfix)

	_, err = splitIntoSubexpressions(tree, subexpressionMap)
	if err != nil {
		return errors.New("error in spliting expression")
	}

	return nil
}

func CreateIdempotentKey(expression []string) string {
	key := strings.Join(expression, "")
	key = strings.ReplaceAll(strings.ReplaceAll(key, "+", "p"), "-", "m")
	key = strings.ReplaceAll(strings.ReplaceAll(key, "*", "u"), "/", "d")
	key = strings.ReplaceAll(strings.ReplaceAll(key, "(", "o"), ")", "c")
	return key
}

func FormatExpression(expression string) []string {
	expression = strings.ReplaceAll(expression, " ", "")
	expression = strings.ReplaceAll(strings.ReplaceAll(expression, "(", "( "), ")", " )")
	expression = strings.ReplaceAll(strings.ReplaceAll(expression, "-", " - "), "+", " + ")
	expression = strings.ReplaceAll(strings.ReplaceAll(expression, "/", " / "), "*", " * ")
	return strings.Fields(expression)
}

func IsValid(expression []string) bool {

	var braces c.Stack

	waitOperator, waitOperand := false, true

	for i, token := range expression {

		if _, err := strconv.ParseFloat(token, 64); err == nil && waitOperand {
			waitOperand = false
			waitOperator = true
		} else if isOperator(token) && waitOperator {
			waitOperator = false
			waitOperand = true
		} else if token == "(" {
			braces.Push(token)
		} else if token == ")" && braces.Top() == "(" && i >= 4 && (expression[i-1] != "(" && expression[i-2] != "(" && expression[i-3] != "(") {
			braces.Pop()
		} else {
			return false
		}
	}

	return len(braces) == 0

}
