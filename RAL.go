package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	//"strconv"
	"strings"
)

//Error messages 
const LexicalErrMsg = "Quid dicis? You offend Caesar with your sloppy lexical habits!"
const SyntaxErrMsg = "Quid dicis? True Romans would not understand your syntax!"
const DeclErrMsg = "Quid dicis? Failure to declare allegiance to Caesar!"
const NegErrMsg = "Quid dicis? Caesar demands positive thoughts!"
const ZeroErrMsg = "Quid dicis? Arab merchants haven't left for India yet!"

//Error codes
const NOERROR = 0
const LEXERROR = -10
const SYNERROR = -100
const DECLERROR = -200
const NEGERROR = -300
const ZEROERROR = -500

//Regex matches function 
func matches(exp string, arg string) bool {
	if arg == "" {
		return false
	}
	c := []byte(arg)
	check, _ := regexp.Match(exp, c)
	return check
}

//Checks if number is valid roman numeral from [1, 3999]
func isRoman(arg string) bool {
	return matches("^M{0,4}(CM|CD|D?C{0,3})(XC|XL|L?X{0,3})(IX|IV|V?I{0,3})$", arg)
}

//Checks if string is a valid operator 
func isOperator(arg string) bool {
	return matches("((plus)|(minus)|(times)|(divide)|(power)|(modulo)|(est))", arg)
}

//Checks if string is a valid string of dictionary
func isValid(arg string) bool {
	return isOperator(arg) || isRoman(arg)
}


func isDigit(c byte) bool { return c >= '0' && c <= '9' }
func isLowerCase(c byte) bool { return c >= 'a' && c <= 'z' }
func isUpperCase(c byte) bool { return c >= 'A' && c <= 'Z' }
func isLetter(c byte) bool {
	return isLowerCase(c) || isUpperCase(c)
}

func isLetterOrDigit(c byte) bool {
	return isDigit(c) || isLetter(c)
}

//Checks if character is apart of alphabet 
func isAlphabet(c byte) bool {
	return isLetter(c) || isParen(c) || c == ' '
}

func isParen(c byte) bool {
	return c == '(' || c == ')'
}

// Returns true if the next token of the expression matches
func isNextToken(token string) bool {
	start := pos 

	if start < len(currentString) && isParen(currentString[start]) {
		return currentString[pos: start + 1 ] == token
	}

	for start < len(currentString) && isLetterOrDigit(currentString[start]) {
		start++ 
	}

	return currentString[pos: start] == token
}

// Returns the next token of the expression consuming the token. 
func getNextToken() string {
	start := pos

	if start < len(currentString) && isParen(currentString[start]) {
		return currentString[start : pos + 1]
	}

	for start < len(currentString) && isLetterOrDigit(currentString[start]){
		start++ 
	}

	return currentString[pos: start]
}

//Checks if the next token is valid without consuming it. 
func nextToken() bool {

	start := pos 
	if currentString[start] == ')' {
		start++ 

		if start < len(currentString) && !checkLex(start) {
			return false
		}

		for start < len(currentString) && currentString[start] == ' ' {
			start++

			if start < len(currentString) && !checkLex(start) {
				return false
			}
		}

		pos = start
		return true 
	}

	for start < len(currentString) && (isLetterOrDigit(currentString[start]) || isParen(currentString[start])) {

		if currentString[start] == ')' {
			pos = start
			return true 
		} else if currentString[start] == '(' {
			start++ 

			if start < len(currentString) && !checkLex(start) {
				return false
			}

			pos = start
			return true 
		}

		start++

		if start < len(currentString) && !checkLex(start) {
			return false
		}
	}

	pos = start
	for start < len(currentString) && currentString[start] == ' ' {
		start++
		
		if start < len(currentString) && !checkLex(start) {
			pos = start
			return false
		}
	}

	pos = start
	return true
}

func checkLex(i int) bool {

	return isAlphabet(currentString[i])
}

func checkOperator() bool {

	return isOperator(getNextToken())
}

func isZero(num int) bool {
	return num == 0
}

func isNeg(num int) bool {
	return num < 0 
}

var variableHolder = make(map[string]int)

var romanNums = make(map[byte]int)
var subRule = make(map[string]int)

var pos = 0
var currentString = ""

var syntaxSwitch = true

func calculateAll(romanExps []string) string {

	//Roman number values
	romanNums['I'] = 1
	romanNums['V'] = 5
	romanNums['X'] = 10
	romanNums['L'] = 50
	romanNums['C'] = 100
	romanNums['D'] = 500
	romanNums['M'] = 1000
 
	subRule["IV"] = 4
	subRule["IX"] = 9
	subRule["XL"] = 40
	subRule["XC"] = 90
	subRule["CD"] = 400
	subRule["CM"] = 900
	

	var lastEle = ""
	for _, element := range romanExps {

		currentString = element

		vars, value, valid := calculateExpression()
		
		variableHolder[vars] = value
		lastEle = vars 

		if valid != 0 {
			x := currentString + "\n"
			x += strings.Repeat(" ", pos)
			x += "^\n"

			if valid == SYNERROR {
				x += SyntaxErrMsg
			} else if valid == LEXERROR {
				x += LexicalErrMsg
			} else if valid == DECLERROR {
				x += DeclErrMsg
			} else if valid == NEGERROR {
				x += NegErrMsg
			} else if valid == ZEROERROR {
				x += ZeroErrMsg
			}

			return x
		}

		pos = 0
	}

	return intToRoman(variableHolder[lastEle])
}

//Calculate value of 
func calculateExpression() (string, int, int) {

	name := getNextToken()

	if !nextToken() {
		return name, pos, LEXERROR
	}

	nextToken()

	x, valid := calcTerm()

	return name, x, valid
}

//Term => term | term {plus, minus} term
func calcTerm() (int, int) {

	x, valid := calcFactor()

	if valid != 0 {
		return pos, valid
	}

	for {

		if isNextToken("plus") {
			
			var saved = pos 
			if !nextToken() {
				return pos, LEXERROR
			}

			if checkOperator() {
				return pos, SYNERROR
			}

			y, valid2 := calcFactor()

			if valid2 != 0 {
				return pos, valid2
			}

			x += y

			if isZero(x) {
				pos = saved 
				return pos, ZEROERROR
			}

			if isNeg(x) {
				pos = saved
				return pos, NEGERROR
			}

		} else if isNextToken("minus") {
			
			var saved = pos 
			if !nextToken() {
				return pos, LEXERROR
			}

			if checkOperator() {
				return pos, SYNERROR
			}

			y, valid2 := calcFactor()

			if valid2 != NOERROR {
				return pos, valid2
			}

			x -= y

			if isZero(x) {
				pos = saved 
				return pos, ZEROERROR
			}

			if isNeg(x) {
				pos = saved
				return pos, NEGERROR
			}

		} else {
			return x, NOERROR
		}
	}
}

//Factor => term | term {times, divide, modulo, power}
func calcFactor() (int, int) {

	x, valid := calcParen()

	if valid != NOERROR {
		return pos, valid
	}

	for {

		if isNextToken("times") { 
			
			var saved = pos 
			if !nextToken() {
				return pos, LEXERROR
			}

			if checkOperator() {
				return pos, SYNERROR
			}

			y, valid2 := calcParen() 

			if valid2 != NOERROR {
				return pos, valid2
			}

			x *= y

			if isZero(x) {
				pos = saved 
				return pos, ZEROERROR
			}

			if isNeg(x) {
				pos = saved
				return pos, NEGERROR
			}

		} else if isNextToken("divide") {
			
			var saved = pos 
			if !nextToken() {
				return pos, LEXERROR
			}

			if checkOperator() {
				return pos, SYNERROR
			}

			y, valid2 := calcParen()
			
			if valid2 != NOERROR {
				return pos, valid2
			}

			x /= y

			if isZero(x) {
				pos = saved
				return saved, ZEROERROR
			}

			if isNeg(x) {
				pos = saved 
				return pos, NEGERROR
			}

		} else if isNextToken("modulo") {

			var saved = pos 
			if !nextToken() {
				return pos, LEXERROR
			}

			if checkOperator() {
				return pos, SYNERROR
			}

			y, valid2 := calcParen()

			if valid2 != NOERROR {
				return pos, valid2
			}

			x %= y

			if isZero(x) {
				pos = saved 
				return pos, ZEROERROR
			}

			if isNeg(x) {
				pos = saved
				return pos, NEGERROR
			}

		} else if isNextToken("power") {

			var saved = pos 
			if !nextToken() {
				return pos, LEXERROR
			}

			if checkOperator() {
				return pos, SYNERROR
			}

			y, valid2 := calcParen()

			if valid2 != NOERROR {
				return pos, valid2
			}

			x = int(math.Pow(float64(x), float64(y)))

			if isZero(x) {
				pos = saved 
				return pos, ZEROERROR
			}

			if isNeg(x) {
				pos = saved
				return pos, NEGERROR
			}

		} else {
			return x, NOERROR
		}
	}
}

//Paren => term | (expression)
func calcParen() (int, int) {

	for isNextToken("(") {
		
		if !nextToken() {
			return pos, LEXERROR
		}

		x, valid := calcTerm()

		if valid != 0 {
			return pos, valid
		}

		if (!isNextToken(")")) {
			return pos, SYNERROR
		}

		nextToken() 

		return x, NOERROR
	}

	x := getNextToken()
	

	if isRoman(x) {

		if !nextToken() {
			return pos, LEXERROR
		}

		if len(getNextToken()) != 0 && !checkOperator() && getNextToken() != ")" {
			return pos, SYNERROR
		}

		rtn,_ := romanToInt(x)
		return rtn, NOERROR
	}

	y, valid := variableHolder[x]

	if !valid {
		return pos, DECLERROR
	}

	if !nextToken() {
		return pos, LEXERROR
	}

	if len(getNextToken()) != 0 && !checkOperator() && getNextToken() != ")" {
		return pos, SYNERROR
	}



	return y, NOERROR
}


/**
* Converts a roman number to an integer. 
* Returns false if conversion fails.
*/
func romanToInt(roman string) (int, bool) {

	rtn := romanNums[roman[0]]

	for i := 1; i < len(roman); i++ {

		last := romanNums[roman[i - 1]]
		current := romanNums[roman[i]]
		
		if last < current {
			value, exist := subRule[roman[i - 1: i + 1]]

			if !exist {
				return i, false
			}

			rtn += value
			rtn -= last
			
		} else {
			rtn += current
		}
	}

	return rtn, true
}

/**
* Converts a roman integer into a string.
*/
func intToRoman(roman int) string {

	rtn := ""

	romanArray := [13] int {1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
	romanArray2 := [13] string {"M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"}

	i := 0 

	for roman > 0  {

		for roman >= romanArray[i] {
			rtn += romanArray2[i]
			roman -= romanArray[i]
		}
		i++
	}

	return rtn
}

/** readLines reads a whole file into memory
* and returns a slice of its lines.
 */
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Printf("input file missing")
		os.Exit(1)
	}
	assignments, err := readLines(args[0])
	if err != nil {
		fmt.Printf("Error with reading input file: %s", err)
	}
	fmt.Println(calculateAll(assignments))
}
