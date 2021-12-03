package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
)

var (
	inputFile = flag.String("input", "", "Input file")
)

func IngestInputFile(inputPath string) []string {
	output := []string{}
	fileHandle, err := os.Open(inputPath)

	if err != nil {
		fmt.Println(err)
		return output
	}

	defer fileHandle.Close()

	scanner := bufio.NewScanner(fileHandle)
	for scanner.Scan() {
		value := scanner.Text()
		output = append(output, value)
	}

	return output
}

func part1(input []string) {
	fieldWidth := len(input[0])
	fields := make([]int, fieldWidth)

	for _, value := range input {
		for i, v := range value {
			if v == '1' {
				fields[i]++
			}
		}
	}

	gammaAccumulator := ""
	epsilonAccumulator := ""

	for _, v := range fields {

		if float64(v)/float64(len(input)) >= 0.5 {
			gammaAccumulator += "1"
			epsilonAccumulator += "0"
		} else {
			gammaAccumulator += "0"
			epsilonAccumulator += "1"
		}
	}

	gamma, _ := strconv.ParseUint(gammaAccumulator, 2, 64)
	epsilon, _ := strconv.ParseUint(epsilonAccumulator, 2, 64)
	fmt.Println(gamma * epsilon)
}

func main() {
	flag.Parse()

	input := IngestInputFile(*inputFile)
	part1(input)
	part2(input)
}

func part2(input []string) {
	resultA := part2a(input, false)

	resultB := part2a(input, true)
	fmt.Println(resultA * resultB)
}

func part2a(input []string, invertValues bool) uint64 {
	currentInput := input
	nextInput := []string{}

	defaultValue, alterateValue := '1', '0'
	if invertValues {
		defaultValue, alterateValue = '0', '1'
	}

	for currentPosition := 0; currentPosition < len(input[0]); currentPosition++ {
		digitCount := 0
		inputLength := len(currentInput)

		for _, value := range currentInput {
			if value[currentPosition] == '1' {
				digitCount++
			}
		}

		targetDigit := defaultValue

		if float64(digitCount)/float64(inputLength) < 0.5 {
			targetDigit = alterateValue
		}

		for _, value := range currentInput {
			runeList := []rune(value)
			if runeList[currentPosition] == targetDigit {
				nextInput = append(nextInput, value)
			}
		}

		if len(nextInput) == 1 {
			returnValue, _ := strconv.ParseUint(nextInput[0], 2, 64)
			return returnValue
		}

		currentInput = nextInput
		nextInput = []string{}
	}
	return 0
}
