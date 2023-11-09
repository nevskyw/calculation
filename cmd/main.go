package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type RomanNumeral struct {
	Value  int
	Symbol string
}

var romanNumerals = []RomanNumeral{
	{1000, "M"},
	{900, "CM"},
	{500, "D"},
	{400, "CD"},
	{100, "C"},
	{90, "XC"},
	{50, "L"},
	{40, "XL"},
	{10, "X"},
	{9, "IX"},
	{5, "V"},
	{4, "IV"},
	{1, "I"},
}

type Calculator struct {
	Operand1 string
	Operand2 string
	Operator string
	isRoman  bool
}

func (c *Calculator) Parse(input string) error {
	operands := strings.Split(input, " ")
	if len(operands) != 3 {
		return errors.New("неверный формат")
	}

	rArabic, _ := regexp.Compile("^(10|[1-9])$")
	rRoman, _ := regexp.Compile("^M{0,3}(CM|CD|D?C{0,3})(XC|XL|L?X{0,3})(IX|IV|V?I{0,3})$")

	for i := 0; i < 3; i++ {
		if i%2 == 0 {
			if !rArabic.MatchString(operands[i]) && !rRoman.MatchString(operands[i]) {
				return errors.New("неправильный формат чисел")
			}

			isRoman := rRoman.MatchString(operands[i])

			if i == 0 {
				c.isRoman = isRoman
			} else if c.isRoman != isRoman {
				return errors.New("используются одновременно разные системы счисления")
			}
		} else {
			matched, _ := regexp.MatchString("^[*/+-]$", operands[i])
			if !matched {
				return errors.New("неправильный оператор")
			}
		}
	}

	c.Operand1, c.Operand2, c.Operator = operands[0], operands[2], operands[1]
	return nil
}

func (c *Calculator) Calculate() (int, error) {
	op1, _ := strconv.Atoi(c.Operand1)
	op2, _ := strconv.Atoi(c.Operand2)
	if c.isRoman {
		var err error
		op1, err = parseRoman(c.Operand1)
		if err != nil {
			return 0, err
		}
		op2, err = parseRoman(c.Operand2)
		if err != nil {
			return 0, err
		}
	}

	switch c.Operator {
	case "+":
		return op1 + op2, nil
	case "-":
		if c.isRoman && op1-op2 < 1 {
			return 0, errors.New("результат не может быть меньше 1 для римских чисел")
		}
		return op1 - op2, nil
	case "*":
		return op1 * op2, nil
	case "/":
		if op2 == 0 {
			return 0, errors.New("Нельзя делить на 0")
		}
		return op1 / op2, nil
	default:
		return 0, errors.New("Неподдерживаемый оператор")
	}
}

func parseRoman(roman string) (int, error) {
	result := 0
	i := 0
	for _, numeral := range romanNumerals {
		for strings.HasPrefix(roman[i:], numeral.Symbol) {
			result += numeral.Value
			i += len(numeral.Symbol)
		}
	}
	return result, nil
}

func formatRoman(num int) string {
	var result strings.Builder
	for _, numeral := range romanNumerals {
		for num >= numeral.Value {
			result.WriteString(numeral.Symbol)
			num -= numeral.Value
		}
	}
	return result.String()
}

func main() {
	calc := Calculator{}
	fmt.Println("Пожалуйста, введите операцию:")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		err := calc.Parse(scanner.Text())
		if err != nil {
			fmt.Println("Ошибка:", err)
			return
		}

		result, err := calc.Calculate()
		if err != nil {
			fmt.Println("Ошибка:", err)
			return
		}

		if calc.isRoman {
			fmt.Println(formatRoman(result))
		} else {
			fmt.Println(result)
		}

		fmt.Println("Пожалуйста, введите другую операцию:")
	}
}
