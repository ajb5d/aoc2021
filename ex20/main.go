package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	inputFile   = flag.String("input", "data/test.txt", "Input file")
	verboseFlag = flag.Bool("verbose", false, "Verbose output")
)

const (
	HIGHLIGHT = "\u001b[31m"
	RESET     = "\u001b[0m"
)

type Point struct {
	x, y int
}

type TrenchMap struct {
	input        []bool
	data         map[Point]bool
	min, max     Point
	defaultValue bool
}

func (t *TrenchMap) updateExtent() {
	for point := range t.data {
		if point.x < t.min.x {
			t.min.x = point.x
		}
		if point.x > t.max.x {
			t.max.x = point.x
		}
		if point.y < t.min.y {
			t.min.y = point.y
		}
		if point.y > t.max.y {
			t.max.y = point.y
		}
	}
}

func (t TrenchMap) string() string {
	output := ""
	for y := t.min.y; y <= t.max.y; y++ {
		for x := t.min.x; x <= t.max.x; x++ {
			if t.data[Point{x, y}] {
				output += "X"
			} else {
				output += "."
			}
		}
		output += "\n"
	}
	return output
}

func (t TrenchMap) pointValue(x, y int) bool {
	if x < t.min.x || x > t.max.x || y < t.min.y || y > t.max.y {
		return t.defaultValue
	}
	return t.data[Point{x, y}]
}

func loadMapFromFile(inputPath string) TrenchMap {
	output := TrenchMap{data: make(map[Point]bool)}

	fileHandle, err := os.Open(inputPath)

	if err != nil {
		panic(err)
	}
	defer fileHandle.Close()

	scanner := bufio.NewScanner(fileHandle)
	scanner.Scan()
	inputLine := scanner.Text()
	for _, value := range inputLine {
		output.input = append(output.input, value == '#')
	}

	scanner.Scan()

	row := 0
	for scanner.Scan() {
		currentLine := scanner.Text()
		currentLine = strings.TrimSpace(currentLine)

		for column, value := range currentLine {
			output.data[Point{column, row}] = (value == '#')
		}
		row++
	}

	output.updateExtent()
	return output
}

func main() {
	flag.Parse()

	input := loadMapFromFile(*inputFile)
	fmt.Printf("%s", input.string())

	for rounds := 0; rounds < 50; rounds++ {
		newImage := make(map[Point]bool)
		for x := input.min.x - 2; x <= input.max.x+2; x++ {
			for y := input.min.y - 2; y <= input.max.y+2; y++ {
				value := 0
				for yOffset := -1; yOffset <= 1; yOffset++ {
					for xOffset := -1; xOffset <= 1; xOffset++ {
						value <<= 1
						if input.pointValue(x+xOffset, y+yOffset) {
							value++
						}
					}
				}
				if input.input[value] {
					newImage[Point{x, y}] = true
				}
			}
		}

		input.data = newImage
		input.updateExtent()
		input.defaultValue = !input.defaultValue
		fmt.Println()

		fmt.Printf("%s", input.string())
	}

	litCount := 0
	for _, value := range input.data {
		if value {
			litCount++
		}
	}
	fmt.Printf("Lit count: %d\n", litCount)
}
