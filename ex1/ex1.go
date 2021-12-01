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

func IngestIntegerListFile(inputPath string) []int {
	output := []int{}
	fileHandle, err := os.Open(inputPath)

	if err != nil {
		fmt.Println(err)
		return output
	}

	defer fileHandle.Close()

	scanner := bufio.NewScanner(fileHandle)
	for scanner.Scan() {
		val, _ := strconv.Atoi(scanner.Text())
		output = append(output, val)
	}

	return output
}

func main() {
	flag.Parse()
	input := IngestIntegerListFile(*inputFile)

	last := -1
	incCount := 0

	for _, val := range input {
		if last != -1 && val > last {
			incCount++
		}
		last = val
	}

	fmt.Printf("Incremented Count: %d\n", incCount)

	last = -1
	incCount = 0

	for index := range input[2:] {

		currentSum := input[index] + input[index+1] + input[index+2]
		if last != -1 && currentSum > last {
			incCount++
		}
		last = currentSum
	}
	fmt.Printf("Incremented Count: %d\n", incCount)

}
