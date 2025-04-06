package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Token types
const (
	NUMBER = iota
	OPERATOR
	LPAREN
	RPAREN
)

type Token struct {
	Type  int
	Value string
}

// Tokenizer: converts input string to tokens
func tokenize(input string) ([]Token, error) {
	var tokens []Token
	var numBuilder strings.Builder

	for _, r := range input {
		switch {
		case unicode.IsDigit(r) || r == '.':
			numBuilder.WriteRune(r)
		case strings.ContainsRune("+-*/", r):
			if numBuilder.Len() > 0 {
				tokens = append(tokens, Token{NUMBER, numBuilder.String()})
				numBuilder.Reset()
			}
			tokens = append(tokens, Token{OPERATOR, string(r)})
		case r == '(':
			if numBuilder.Len() > 0 {
				tokens = append(tokens, Token{NUMBER, numBuilder.String()})
				numBuilder.Reset()
			}
			tokens = append(tokens, Token{LPAREN, string(r)})
		case r == ')':
			if numBuilder.Len() > 0 {
				tokens = append(tokens, Token{NUMBER, numBuilder.String()})
				numBuilder.Reset()
			}
			tokens = append(tokens, Token{RPAREN, string(r)})
		case unicode.IsSpace(r):
			if numBuilder.Len() > 0 {
				tokens = append(tokens, Token{NUMBER, numBuilder.String()})
				numBuilder.Reset()
			}
		default:
			return nil, fmt.Errorf("invalid character: %c", r)
		}
	}

	if numBuilder.Len() > 0 {
		tokens = append(tokens, Token{NUMBER, numBuilder.String()})
	}

	return tokens, nil
}

// Shunting Yard Algorithm to convert infix to postfix
func toPostfix(tokens []Token) ([]Token, error) {
	var output []Token
	var stack []Token

	precedence := map[string]int{
		"+": 1,
		"-": 1,
		"*": 2,
		"/": 2,
	}

	for _, token := range tokens {
		switch token.Type {
		case NUMBER:
			output = append(output, token)
		case OPERATOR:
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				if top.Type == OPERATOR && precedence[top.Value] >= precedence[token.Value] {
					output = append(output, top)
					stack = stack[:len(stack)-1]
				} else {
					break
				}
			}
			stack = append(stack, token)
		case LPAREN:
			stack = append(stack, token)
		case RPAREN:
			for len(stack) > 0 && stack[len(stack)-1].Type != LPAREN {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 || stack[len(stack)-1].Type != LPAREN {
				return nil, fmt.Errorf("mismatched parentheses")
			}
			stack = stack[:len(stack)-1]
		}
	}

	for len(stack) > 0 {
		if stack[len(stack)-1].Type == LPAREN {
			return nil, fmt.Errorf("mismatched parentheses")
		}
		output = append(output, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return output, nil
}

// Evaluator for postfix expression
func evaluatePostfix(tokens []Token) (float64, error) {
	var stack []float64

	for _, token := range tokens {
		switch token.Type {
		case NUMBER:
			num, err := strconv.ParseFloat(token.Value, 64)
			if err != nil {
				return 0, err
			}
			stack = append(stack, num)
		case OPERATOR:
			if len(stack) < 2 {
				return 0, fmt.Errorf("not enough operands for operator %s", token.Value)
			}
			b, a := stack[len(stack)-1], stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			switch token.Value {
			case "+":
				stack = append(stack, a+b)
			case "-":
				stack = append(stack, a-b)
			case "*":
				stack = append(stack, a*b)
			case "/":
				if b == 0 {
					return 0, fmt.Errorf("division by zero")
				}
				stack = append(stack, a/b)
			}
		}
	}

	if len(stack) != 1 {
		return 0, fmt.Errorf("invalid expression")
	}

	return stack[0], nil
}

func calculate(input string) (float64, error) {
	tokens, err := tokenize(input)
	if err != nil {
		return 0, err
	}

	postfix, err := toPostfix(tokens)
	if err != nil {
		return 0, err
	}

	return evaluatePostfix(postfix)
}

func main() {
	fmt.Println("Enter a math expression:")
	var input string
	fmt.Scanln(&input)

	result, err := calculate(input)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Result = %v\n", result)
	}
}

