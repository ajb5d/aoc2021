package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
)

var (
	inputFile = flag.String("input", "", "Input file")
)

func loadMapFromFile(inputPath string) []string {
	output := []string{}

	fileHandle, err := os.Open(inputPath)

	if err != nil {
		panic(err)
	}
	defer fileHandle.Close()

	scanner := bufio.NewScanner(fileHandle)

	for scanner.Scan() {
		line := scanner.Text()
		output = append(output, strings.TrimSpace(line))
	}
	return output
}

var openClose = map[rune]rune{
	'{': '}',
	'[': ']',
	'(': ')',
	'<': '>',
}

var scoresForMismatched = map[rune]int{
	')': 3,
	']': 57,
	'}': 1197,
	'>': 25137,
}

var scoresForUnmatched = map[rune]int{
	'(': 1,
	'[': 2,
	'{': 3,
	'<': 4,
}

const (
	OPEN_SYMBOLS  = "({[<"
	CLOSE_SYMBOLS = ")}]>"
)

func runesMatch(new, stack rune) bool {
	value, ok := openClose[stack]
	if !ok || value != new {
		return false
	}
	return true
}

func lineIsValid(line string) (int, int) {

	stack := []rune{}
	for _, char := range line {
		if strings.ContainsRune(OPEN_SYMBOLS, char) {
			stack = append(stack, char)
		}
		if strings.ContainsRune(CLOSE_SYMBOLS, char) {
			if runesMatch(char, stack[len(stack)-1]) {
				stack = stack[:len(stack)-1]
			} else {
				//fmt.Printf("Expected %s, got %s\n", string(openClose[stack[len(stack)-1]]), string(char))
				return 0, scoresForMismatched[char]
			}
		}
	}
	fmt.Println("Stack:", string(stack))
	totalScore := 0
	for i := len(stack) - 1; i >= 0; i-- {
		totalScore = totalScore*5 + scoresForUnmatched[stack[i]]
	}
	return totalScore, 0
}

func main() {
	flag.Parse()

	input := loadMapFromFile(*inputFile)
	totalMismatched := 0
	unmatchedScores := []int{}
	for _, value := range input {
		unmatchedScore, mismatchedScore := lineIsValid(value)
		totalMismatched += mismatchedScore
		if unmatchedScore > 0 {
			unmatchedScores = append(unmatchedScores, unmatchedScore)
		}
	}

	fmt.Println("Total mismatched score:", totalMismatched)

	sort.Ints(unmatchedScores)

	fmt.Println("Unmatched scores:", unmatchedScores[len(unmatchedScores)/2])
}
