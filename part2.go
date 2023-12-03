package main

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
)

const (
	TOKEN_TYPE_NUMBER = iota
	TOKEN_TYPE_PERIOD = iota
	TOKEN_TYPE_SYMBOL = iota
	TOKEN_TYPE_NULL   = iota
)

type Token struct {
	value      string
	token_type int
	span       int
}

type Transform struct {
	x      int
	y      int
	width  int
	height int
}

type Cell struct {
	token     Token
	transform Transform
}

func readFile(path string) []string {
	file, _ := os.Open(path)

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}

func readToken(str string, offset int) Token {
	regexNumber, _ := regexp.Compile(`^(\d+)`)
	regexPeriod, _ := regexp.Compile(`^\.`)
	regexSymbol, _ := regexp.Compile(`^\S`)

	str = str[offset:len(str)]

	if match := regexNumber.FindString(str); match != "" {
		//println("number " + match + "  len: " + strconv.Itoa(len(match)))
		return Token{
			match,
			TOKEN_TYPE_NUMBER,
			len(match),
		}
	}
	if match := regexPeriod.FindString(str); match != "" {
		//println("period " + match + "  len: " + strconv.Itoa(len(match)))
		return Token{
			match,
			TOKEN_TYPE_PERIOD,
			len(match),
		}
	}
	if match := regexSymbol.FindString(str); match != "" {
		//println("symbol " + match + "  len: " + strconv.Itoa(len(match)))
		return Token{
			match,
			TOKEN_TYPE_SYMBOL,
			len(match),
		}
	}

	//println("null")
	return Token{
		"",
		TOKEN_TYPE_NULL,
		0,
	}
}

func getGrid(lines []string) [][]Cell {
	tokens := make([][]Token, len(lines))
	cells := make([][]Cell, len(lines))
	for i := range lines {
		tokens[i] = make([]Token, 0)
	}

	for i, line := range lines {
		//fmt.Println(i, line)

		j := 0
		offset := 0
		token := Token{}
		for token.token_type != TOKEN_TYPE_NULL {
			token = readToken(line, offset)

			if token.token_type == TOKEN_TYPE_NULL {
				continue
			}

			var x = offset
			var y = i
			var width = token.span
			var height = 1

			cells[i] = append(cells[i], Cell{
				token,
				Transform{x, y, width, height},
			})

			tokens[i] = append(tokens[i], token)
			offset += token.span
			j++
		}
	}

	return cells
}

func intersects(a Transform, b Transform) bool {
	return a.x <= b.x+b.width &&
		a.x+a.width >= b.x &&
		a.y <= b.y+b.height &&
		a.y+a.height >= b.y
}

func intersectsSymbol(a Transform, grid [][]Cell) bool {
	cellsWithSymbols := make([]Cell, 0)

	for _, row := range grid {
		for _, cell := range row {
			if cell.token.token_type == TOKEN_TYPE_SYMBOL {
				cellsWithSymbols = append(cellsWithSymbols, cell)
			}
		}
	}

	for _, cell := range cellsWithSymbols {

		if intersects(a, cell.transform) {
			//println(cell.transform.x, cell.transform.y, cell.transform.width, cell.transform.height, cell.token.value, "intersects")
			return true
		}
	}

	return false
}

func intersectsNumbers(a Transform, grid [][]Cell) []int {
	numbers := make([]int, 0)
	cellsWithNumbers := make([]Cell, 0)

	for _, row := range grid {
		for _, cell := range row {
			if cell.token.token_type == TOKEN_TYPE_NUMBER {
				cellsWithNumbers = append(cellsWithNumbers, cell)
			}
		}
	}

	intersectionCells := make([]Cell, 0)
	for _, cell := range cellsWithNumbers {
		if intersects(a, cell.transform) {
			num, _ := strconv.Atoi(cell.token.value)
			intersectionCells = append(intersectionCells, cell)
			numbers = append(numbers, num)
			//println(cell.transform.x, cell.transform.y, cell.transform.width, cell.transform.height, cell.token.value, "intersects")
		}
	}

	return numbers
}

func main() {
	//read file into array of lines
	fileLines := readFile("input.txt")

	//2d slice of tokens
	grid := getGrid(fileLines)

	//print grid
	/*
		for _, line := range grid {
			for _, cell := range line {
				println(cell.transform.x, cell.transform.y, cell.transform.width, cell.transform.height, cell.token.value)
			}
		}
	*/

	numberCells := make([]Cell, 0)
	for _, row := range grid {
		for _, cell := range row {
			if cell.token.token_type == TOKEN_TYPE_NUMBER {
				numberCells = append(numberCells, cell)
			}
		}
	}

	cellsWithAsterisk := make([]Cell, 0)
	for _, cell := range grid {
		for _, cell := range cell {
			if cell.token.token_type == TOKEN_TYPE_SYMBOL && cell.token.value == "*" {
				cellsWithAsterisk = append(cellsWithAsterisk, cell)
			}
		}
	}

	prod_sum := 0
	for _, cell := range cellsWithAsterisk {
		intersections := intersectsNumbers(cell.transform, grid)
		if len(intersections) >= 2 {
			println("found intersections ", cell.transform.x, cell.transform.y, cell.transform.width, cell.transform.height, cell.token.value)
			prod := 1
			for _, num := range intersections {
				prod *= num
				println(num)
			}
			prod_sum += prod
			println("prod: ", prod)
			println()
		}
	}

	println("prod_sum: ", prod_sum)

	//for i, line := range fileLines {
	//	fmt.Println(i, line)
	//}
}
