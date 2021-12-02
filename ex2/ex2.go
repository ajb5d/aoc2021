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

func IngestCourseFile(inputPath string) [][]string {
	output := [][]string{}
	fileHandle, err := os.Open(inputPath)

	if err != nil {
		fmt.Println(err)
		return output
	}

	defer fileHandle.Close()

	scanner := bufio.NewScanner(fileHandle)
	for scanner.Scan() {
		fields := strings.Fields(strings.ToLower(scanner.Text()))
		output = append(output, fields)
	}

	return output
}

func main() {
	flag.Parse()

	input := IngestCourseFile(*inputFile)

	horizontal, depth := 0, 0

	for _, course := range input {
		direction := course[0]
		magnitude, _ := strconv.Atoi(course[1])

		switch direction {
		case "forward":
			horizontal += magnitude
		case "down":
			depth += magnitude
		case "up":
			depth -= magnitude
		}

	}

	fmt.Printf("%d %d\n", horizontal, depth)
	fmt.Printf("%d\n", horizontal*depth)

	horizontal, depth = 0, 0
	aim := 0

	for _, course := range input {
		direction := course[0]
		magnitude, _ := strconv.Atoi(course[1])

		switch direction {
		case "forward":
			horizontal += magnitude
			depth += magnitude * aim
		case "down":
			aim += magnitude
		case "up":
			aim -= magnitude
		}
	}

	fmt.Printf("%d %d\n", horizontal, depth)
	fmt.Printf("%d\n", horizontal*depth)
}
