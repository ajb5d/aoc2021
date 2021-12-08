package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

/*
  0:      1:      2:      3:      4:
 aaaa    ....    aaaa    aaaa    ....
b    c  .    c  .    c  .    c  b    c
b    c  .    c  .    c  .    c  b    c
 ....    ....    dddd    dddd    dddd
e    f  .    f  e    .  .    f  .    f
e    f  .    f  e    .  .    f  .    f
 gggg    ....    gggg    gggg    ....

  5:      6:      7:      8:      9:
 aaaa    aaaa    aaaa    aaaa    aaaa
b    .  b    .  .    c  b    c  b    c
b    .  b    .  .    c  b    c  b    c
 dddd    dddd    ....    dddd    dddd
.    f  e    f  .    f  e    f  .    f
.    f  e    f  .    f  e    f  .    f
 gggg    gggg    ....    gggg    gggg

 a: 8
 b: 6
 c: 8
 d: 7
 e: 4
 f: 9
 g: 7
*/
var (
	inputFile = flag.String("input", "", "Input file")
)

var digit_lookup = map[string]string{
	"abcefg":  "0",
	"cf":      "1",
	"acdeg":   "2",
	"acdfg":   "3",
	"bcdf":    "4",
	"abdfg":   "5",
	"abdefg":  "6",
	"acf":     "7",
	"abcdefg": "8",
	"abcdfg":  "9",
}

type Line struct {
	inputs, outputs []string
}

func loadInputsFromFile(inputPath string) []Line {
	output := make([]Line, 0)

	fileHandle, err := os.Open(inputPath)

	if err != nil {
		panic(err)
	}
	defer fileHandle.Close()

	scanner := bufio.NewScanner(fileHandle)

	for scanner.Scan() {
		line := scanner.Text()
		currentLine := Line{}

		parts := strings.Split(line, " | ")

		currentLine.inputs = strings.Fields(parts[0])
		currentLine.outputs = strings.Fields(parts[1])
		output = append(output, currentLine)
	}
	return output
}

func getMapping(line Line) [7]int {
	partCounts := [7]int{}
	assignments := [7]int{-1, -1, -1, -1, -1, -1, -1}

	for _, element := range line.inputs {
		for _, part := range element {
			partCounts[int(part-'a')]++
		}
	}

	// Assign based on frequency
	for i, count := range partCounts {
		switch count {
		case 4:
			//e
			assignments[4] = i
		case 6:
			//b
			assignments[1] = i
		case 9:
			//f
			assignments[5] = i
		}
	}
	// Known: bef
	// Assign based on Pattern: Pattern for '1' has segments 'cf' -> Find 'c'
	for _, element := range line.inputs {
		if len(element) == 2 {
			fSegment := byte(assignments[5] + 'a')
			if element[0] == fSegment {
				assignments[2] = int(element[1] - 'a')
			} else {
				assignments[2] = int(element[0] - 'a')
			}
			break
		}
	}
	// Known bcef
	// Pattern for '7' has segments 'acf' -> find 'a'
	for _, element := range line.inputs {
		if len(element) == 3 {
			fSegment := byte(assignments[5] + 'a')
			cSegment := byte(assignments[2] + 'a')

			for j := 0; j < 3; j++ {
				if element[j] != fSegment && element[j] != cSegment {
					assignments[0] = int(element[j] - 'a')
					break
				}
			}
		}
	}
	// Known abcef
	// Pattern for '4' has 'bcdf' -> find 'd'
	for _, element := range line.inputs {
		if len(element) == 4 {
			fSegment := byte(assignments[5] + 'a')
			cSegment := byte(assignments[2] + 'a')
			bSegment := byte(assignments[1] + 'a')

			for j := 0; j < 4; j++ {
				if element[j] != fSegment && element[j] != cSegment && element[j] != bSegment {
					assignments[3] = int(element[j] - 'a')
					break
				}
			}
		}
	}
	// Known acbdef
	// The last one is 'g'
	for i := 0; i < 8; i++ {
		found := false
		for _, element := range assignments {
			if element == i {
				found = true
				break
			}
		}
		if !found {
			assignments[6] = i
			break
		}
	}

	// fmt.Println(line.inputs)
	// for i := 0; i < 7; i++ {
	// 	fmt.Printf("%c -> %c ", byte(i+'a'), byte(assignments[i]+'a'))
	// }
	// fmt.Println()
	return assignments
}

func translateString(input string, mapping [7]int) []string {
	output := []string{}
	reverseMapping := reverseMapping(mapping)
	for _, element := range input {
		output = append(output, string(byte(reverseMapping[int(element-'a')]+'a')))
	}
	return output
}

func reverseMapping(mapping [7]int) [7]int {
	output := [7]int{}
	for i, element := range mapping {
		output[element] = i
	}
	return output
}

func main() {
	flag.Parse()

	input := loadInputsFromFile(*inputFile)

	count := 0
	total := 0
	for _, line := range input {
		//Part 1
		for _, element := range line.outputs {
			if len(element) == 2 || len(element) == 3 || len(element) == 4 || len(element) == 7 {
				count++
			}
		}

		//Part 2
		assignments := getMapping(line)
		outputSignal := ""
		for _, element := range line.outputs {
			outputStrings := translateString(element, assignments)
			sort.Strings(outputStrings)
			outputString := strings.Join(outputStrings, "")
			outputSignal += digit_lookup[outputString]
		}
		value, _ := strconv.Atoi(outputSignal)
		total += value
	}
	fmt.Println(count)
	fmt.Println(total)

}
