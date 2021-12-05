package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	inputFile = flag.String("input", "", "Input file")
)

type FloorMap struct {
	Extent Point
	Lines  []Line
}

type Line struct {
	start Point
	end   Point
}

type Point struct {
	X int
	Y int
}

func StringToPoint(input string) Point {
	output := Point{}
	for i, element := range strings.Split(input, ",") {
		intValue, _ := strconv.Atoi(element)
		if i == 0 {
			output.X = intValue
		} else {
			output.Y = intValue
		}
	}
	return output
}

func IngestInputFile(inputPath string) FloorMap {
	output := FloorMap{}
	fileHandle, err := os.Open(inputPath)

	if err != nil {
		fmt.Println(err)
		return output
	}

	defer fileHandle.Close()

	scanner := bufio.NewScanner(fileHandle)
	for scanner.Scan() {
		currentLine := Line{}
		inputLine := scanner.Text()
		fields := strings.Fields(inputLine)
		currentLine.start = StringToPoint(fields[0])
		currentLine.end = StringToPoint(fields[2])
		output.Lines = append(output.Lines, currentLine)

		if currentLine.start.X > output.Extent.X {
			output.Extent.X = currentLine.start.X
		}
		if currentLine.start.Y > output.Extent.Y {
			output.Extent.Y = currentLine.start.Y
		}

		if currentLine.end.X > output.Extent.X {
			output.Extent.X = currentLine.end.X
		}

		if currentLine.end.Y > output.Extent.Y {
			output.Extent.Y = currentLine.end.Y
		}
	}

	return output
}

func main() {
	flag.Parse()

	input := IngestInputFile(*inputFile)

	scoreMap(input, false)
	scoreMap(input, true)

}

func scoreMap(input FloorMap, includeDiagnonal bool) int {
	xStride := input.Extent.X + 1
	yStride := input.Extent.Y + 1

	field := make([][]byte, xStride)
	cells := make([]byte, xStride*yStride)

	for i := range field {
		field[i], cells = cells[:yStride], cells[yStride:]
	}

	for _, line := range input.Lines {
		targets := expandLine(line)

		if !includeDiagnonal && lineIsDiagnoal(line) {
			continue
		}

		for _, target := range targets {
			field[target.X][target.Y] += 1
		}
	}

	fieldCount := 0
	for i := range field {
		for j := range field[0] {
			if field[i][j] > 1 {
				fieldCount += 1
			}
		}
	}

	fmt.Println(fieldCount)
	return fieldCount
}

func lineIsDiagnoal(line Line) bool {
	return line.start.X != line.end.X && line.start.Y != line.end.Y
}

func expandSequence(start int, end int) []int {
	output := []int{}
	if start < end {
		for i := start; i <= end; i++ {
			output = append(output, i)
		}
	} else {
		for i := start; i >= end; i-- {
			output = append(output, i)
		}
	}
	return output
}

func expandLine(line Line) []Point {
	output := []Point{}

	if line.start.X == line.end.X {
		for _, value := range expandSequence(line.start.Y, line.end.Y) {
			output = append(output, Point{line.start.X, value})
		}
		return output
	}

	if line.start.Y == line.end.Y {
		for _, value := range expandSequence(line.start.X, line.end.X) {
			output = append(output, Point{value, line.start.Y})
		}
		return output
	}

	xSequence := expandSequence(line.start.X, line.end.X)
	ySequence := expandSequence(line.start.Y, line.end.Y)

	for index := range xSequence {
		output = append(output, Point{xSequence[index], ySequence[index]})
	}

	return output
}
