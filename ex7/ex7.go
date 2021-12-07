package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

var (
	inputFile = flag.String("input", "", "Input file")
)

func loadInputsFromFile(input string) []int {
	output := make([]int, 0)

	bytes, err := ioutil.ReadFile(input)
	if err != nil {
		panic(err)
	}

	parts := strings.Split(strings.TrimSpace(string(bytes)), ",")

	for _, part := range parts {
		val, _ := strconv.Atoi(part)
		output = append(output, val)
	}

	return output
}

const PART_ONE = false

func cost(old, new int) int {
	diff := new - old
	if diff < 0 {
		diff = -diff
	}
	if PART_ONE {
		return diff
	}

	diff = diff * (diff + 1) / 2
	return diff
}

func scorePosition(position int, inputs []int) int {
	sum := 0
	for _, value := range inputs {
		sum += cost(position, value)
	}
	return sum
}

func minMax(input []int) (int, int) {
	min := input[0]
	max := input[0]

	for _, value := range input {
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}

	return min, max
}

func main() {
	flag.Parse()

	input := loadInputsFromFile(*inputFile)
	sort.Ints(input)

	minScore := scorePosition(0, input)
	minPosition := 0

	start, end := minMax(input)

	for i := start; i <= end; i++ {
		score := scorePosition(i, input)
		if score < minScore {
			minScore = score
			minPosition = i
		}
	}

	fmt.Println(minPosition)
	fmt.Println(minScore)
}
