package main

import (
	"fmt"
	"strconv"
	"strings"
)

func parseInputs(input string) []int {
	output := make([]int, 32)
	parts := strings.Split(input, ",")

	for _, part := range parts {
		val, _ := strconv.Atoi(part)
		output[val] += 1
	}

	return output

}

const (
	testInput   = "3,4,1,5"
	puzzleInput = "1,1,3,5,1,1,1,4,1,5,1,1,1,1,1,1,1,3,1,1,1,1,2,5,1,1,1,1,1,2,1,4,1,4,1,1,1,1,1,3,1,1,5,1,1,1,4,1,1,1,4,1,1,3,5,1,1,1,1,4,1,5,4,1,1,2,3,2,1,1,1,1,1,1,1,1,1,1,1,1,1,5,1,1,1,1,1,1,1,1,2,2,1,1,1,1,1,5,1,1,1,3,4,1,1,1,1,3,1,1,1,1,1,4,1,1,3,1,1,3,1,1,1,1,1,3,1,5,2,3,1,2,3,1,1,2,1,2,4,5,1,5,1,4,1,1,1,1,2,1,5,1,1,1,1,1,5,1,1,3,1,1,1,1,1,1,4,1,2,1,1,1,1,1,1,1,1,1,1,1,1,1,3,2,1,1,1,1,2,2,1,2,1,1,1,5,5,1,1,1,1,1,1,1,1,1,1,1,1,2,2,1,1,4,2,1,4,1,1,1,1,1,1,1,2,1,2,1,1,1,1,1,1,1,1,1,1,1,1,1,2,2,1,5,1,1,1,1,1,1,1,1,3,1,1,3,3,1,1,1,3,5,1,1,4,1,1,1,1,1,4,1,1,3,1,1,1,1,1,1,1,1,2,1,5,1,1,1,1,1,1,1,1,1,1,4,1,1,1,1"
)

func main() {
	state := parseInputs(puzzleInput)

	for day := 0; day < 256; day++ {
		newState := make([]int, 32)
		newState[6] += state[0]
		newState[8] += state[0]
		for i := 1; i < len(state); i++ {
			newState[i-1] += state[i]
		}
		state = newState

		sum := 0
		for i := 0; i < len(state); i++ {
			sum += state[i]
		}
		fmt.Printf("Day %d: %d\n", day, sum)
	}
}
