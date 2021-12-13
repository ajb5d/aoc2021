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

type Point struct {
	x, y int
}

type Fold struct {
	axis     string
	position int
}

type CaveMap struct {
	nodes     map[Point]bool
	maxExtent Point
	folds     []Fold
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func loadMapFromFile(inputPath string) CaveMap {
	output := CaveMap{nodes: make(map[Point]bool)}
	// output.nodes = make(map[Point]bool)

	fileHandle, err := os.Open(inputPath)

	if err != nil {
		panic(err)
	}
	defer fileHandle.Close()

	scanner := bufio.NewScanner(fileHandle)
	xMax, yMax := 0, 0
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			break
		}

		elements := strings.Split(scanner.Text(), ",")

		x, _ := strconv.Atoi(elements[0])
		y, _ := strconv.Atoi(elements[1])

		xMax = max(xMax, x)
		yMax = max(yMax, y)

		output.nodes[Point{x, y}] = true
	}

	output.maxExtent = Point{xMax, yMax}
	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.Split(strings.Fields(line)[2], "=")
		x, _ := strconv.Atoi(parts[1])
		output.folds = append(output.folds, Fold{parts[0], x})
	}

	return output
}

func printMap(input CaveMap) {
	fmt.Println("Extent:", input.maxExtent)
	for y := 0; y <= input.maxExtent.y; y++ {
		for x := 0; x <= input.maxExtent.x; x++ {
			if input.nodes[Point{x, y}] {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func printMapWithYLine(input CaveMap, yLine int) {
	for y := 0; y <= input.maxExtent.y; y++ {
		for x := 0; x <= input.maxExtent.x; x++ {
			if y == yLine {
				fmt.Print("-")
			} else {
				if input.nodes[Point{x, y}] {
					fmt.Print("#")
				} else {
					fmt.Print(".")
				}
			}
		}
		fmt.Println()
	}
}

func main() {
	flag.Parse()

	input := loadMapFromFile(*inputFile)

	fmt.Println("Part 1:")

	for index, fold := range input.folds {

		newMap := CaveMap{nodes: make(map[Point]bool), maxExtent: input.maxExtent}

		if fold.axis == "x" {
			newMap.maxExtent.x = fold.position - 1
			for node := range input.nodes {
				if node.x < fold.position {
					newMap.nodes[node] = true
				} else {
					newPoint := Point{input.maxExtent.x - node.x, node.y}
					newMap.nodes[newPoint] = true
				}
			}
		} else {
			newMap.maxExtent.y = fold.position - 1
			for node := range input.nodes {
				if node.y < fold.position {
					newMap.nodes[node] = true
				} else {
					newPoint := Point{node.x, input.maxExtent.y - node.y}
					newMap.nodes[newPoint] = true
				}
			}
		}
		input = newMap
		fmt.Printf("Fold %d: %d\n", index+1, len(input.nodes))
	}

	printMap(input)

}
