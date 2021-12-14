package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"reflect"
	"sort"
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
	}

	runeCounts := make(map[string]int)
	for key, value := range counts {
		runeCounts[string(key[1])] += value
	}

	//Add Initial N
	runeCounts["N"] += 1

	keys := reflect.ValueOf(runeCounts).MapKeys()
	sort.Slice(keys, func(i, j int) bool {
		return runeCounts[keys[i].String()] > runeCounts[keys[j].String()]
	})

	mostCommonCount, leastCommonCount := runeCounts[keys[0].String()], runeCounts[keys[len(keys)-1].String()]
	fmt.Printf("%d\n", mostCommonCount-leastCommonCount)

}
