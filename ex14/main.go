package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	inputFile = flag.String("input", "", "Input file")
)

type PolymerDetails struct {
	start        string
	replacements map[string]string
}

func loadMapFromFile(inputPath string) PolymerDetails {
	output := PolymerDetails{replacements: make(map[string]string)}

	fileHandle, err := os.Open(inputPath)

	if err != nil {
		panic(err)
	}
	defer fileHandle.Close()

	scanner := bufio.NewScanner(fileHandle)

	scanner.Scan()
	output.start = scanner.Text()
	scanner.Scan() // Ignore second line

	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), " -> ")
		output.replacements[parts[0]] = parts[1]
	}

	return output
}

func pairsForString(input string) map[string]int {
	output := make(map[string]int)
	for index := range input[:len(input)-1] {
		currentPair := string(input[index]) + string(input[index+1])
		output[currentPair]++
	}
	return output
}

func main() {
	flag.Parse()

	input := loadMapFromFile(*inputFile)
	fmt.Println("Part 1:")

	counts := pairsForString(input.start)

	for step := 1; step <= 40; step++ {
		newCounts := make(map[string]int)
		for targetPair := range counts {
			replacement, ok := input.replacements[targetPair]
			if ok {
				newCounts[string(targetPair[0])+replacement] += counts[targetPair]
				newCounts[replacement+string(targetPair[1])] += counts[targetPair]
			} else {
				newCounts[targetPair] += counts[targetPair]
			}
		}
		counts = newCounts
		// fmt.Println(counts)
	}

	runeCounts := make(map[string]int)
	for key, value := range counts {
		runeCounts[string(key[1])] += value
	}

	//Add Terminal B
	runeCounts["N"] += 1

	fmt.Println(runeCounts)

	mostCommonCount, leastCommonCount := 0, 0
	for _, count := range runeCounts {
		if count > mostCommonCount {
			mostCommonCount = count
		}
		if count < leastCommonCount || leastCommonCount == 0 {
			leastCommonCount = count
		}
	}
	fmt.Printf("Part 1: %d\n", mostCommonCount-leastCommonCount)

}
