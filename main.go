package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func isOperator(r rune) bool {
	return r == '+' || r == '-' || r == '/' || r == '*'
}

func IsValidFormula(expression string) bool {
	prevWasOperator := true
	stack := make([]rune, 0)

	for _, r := range expression {
		switch {
		case unicode.IsDigit(r) || r == '.':
			prevWasOperator = false
		case r == '(':
			stack = append(stack, r)
			prevWasOperator = true
		case r == ')':
			if len(stack) == 0 {
				return false
			}
			stack = stack[:len(stack)-1]
			prevWasOperator = false
		case isOperator(r):
			if prevWasOperator {
				return false
			}
			prevWasOperator = true
		default:
			return false
		}
	}

	if len(stack) > 0 || prevWasOperator {
		return false
	}

	return true
}

func applyOperation(numbers_stack *[]float64, operator rune) error {
	if len(*numbers_stack) < 2 {
		return errors.New("недостаточно чисел для операции")
	}

	b := (*numbers_stack)[len(*numbers_stack)-1]
	a := (*numbers_stack)[len(*numbers_stack)-2]
	*numbers_stack = (*numbers_stack)[:len(*numbers_stack)-2]

	var result float64
	switch operator {
	case '+':
		result = a + b
	case '-':
		result = a - b
	case '*':
		result = a * b
	case '/':
		if b == 0 {
			return errors.New("деление на ноль")
		}
		result = a / b
	}
	*numbers_stack = append(*numbers_stack, result)
	return nil
}

func precedence(op rune) int {
	switch op {
	case '+', '-':
		return 1
	case '*', '/':
		return 2
	}
	return 0
}

func Calc(expression string) (float64, error) {
	expression = strings.ReplaceAll(expression, "\r", "")
	expression = strings.ReplaceAll(expression, "\n", "")
	trimmed := strings.ReplaceAll(expression, " ", "")

	if !IsValidFormula(trimmed) {
		return 0.0, errors.New("кекоректаня формула")
	}

	numbers_stack := make([]float64, 0)
	operators_stack := make([]rune, 0)

	var buffer []rune
	for _, r := range trimmed {
		switch {
		case unicode.IsDigit(r) || r == '.':
			buffer = append(buffer, r)

		case r == '+' || r == '-' || r == '*' || r == '/':
			if len(buffer) > 0 {
				num, err := strconv.ParseFloat(string(buffer), 64)
				if err != nil {
					return 0.0, errors.New("ошибка преобразования числа")
				}
				numbers_stack = append(numbers_stack, num)
				buffer = buffer[:0]
			}

			for len(operators_stack) > 0 && precedence(operators_stack[len(operators_stack)-1]) >= precedence(r) {
				if operators_stack[len(operators_stack)-1] == '(' {
					break
				}
				if err := applyOperation(&numbers_stack, operators_stack[len(operators_stack)-1]); err != nil {
					return 0.0, err
				}
				operators_stack = operators_stack[:len(operators_stack)-1]
			}
			operators_stack = append(operators_stack, r)

		case r == '(':
			operators_stack = append(operators_stack, r)

		case r == ')':
			if len(buffer) > 0 {
				num, err := strconv.ParseFloat(string(buffer), 64)
				if err != nil {
					return 0.0, errors.New("ошибка преобразования числа")
				}
				numbers_stack = append(numbers_stack, num)
				buffer = buffer[:0]
			}

			for len(operators_stack) > 0 && operators_stack[len(operators_stack)-1] != '(' {
				if err := applyOperation(&numbers_stack, operators_stack[len(operators_stack)-1]); err != nil {
					return 0.0, err
				}
				operators_stack = operators_stack[:len(operators_stack)-1]
			}

			if len(operators_stack) == 0 || operators_stack[len(operators_stack)-1] != '(' {
				return 0.0, errors.New("несбалансированные скобки")
			}
			operators_stack = operators_stack[:len(operators_stack)-1]
		}
	}

	if len(buffer) > 0 {
		num, err := strconv.ParseFloat(string(buffer), 64)
		if err != nil {
			return 0.0, errors.New("ошибка преобразования числа")
		}
		numbers_stack = append(numbers_stack, num)
	}

	for len(operators_stack) > 0 {
		if operators_stack[len(operators_stack)-1] == '(' {
			return 0.0, errors.New("несбалансированные скобки")
		}
		if err := applyOperation(&numbers_stack, operators_stack[len(operators_stack)-1]); err != nil {
			return 0.0, err
		}
		operators_stack = operators_stack[:len(operators_stack)-1]
	}

	if len(numbers_stack) != 1 {
		return 0.0, errors.New("ошибка вычислений")
	}

	return numbers_stack[0], nil
}

func main() {
	expression := "(1   +   2   ) /3"
	result, err := Calc(expression)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Printf("Результат: %.2f\n", result)
	}
}
