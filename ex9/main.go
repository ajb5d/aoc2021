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

var (
	inputFile = flag.String("input", "", "Input file")
)

type HeightMap struct {
	depth [][]int
}

type LowPoint struct {
	y, x int
}

func loadMapFromFile(inputPath string) HeightMap {
	output := HeightMap{}

	fileHandle, err := os.Open(inputPath)

	if err != nil {
		panic(err)
	}
	defer fileHandle.Close()

	scanner := bufio.NewScanner(fileHandle)

	for scanner.Scan() {
		line := scanner.Text()
		currentLine := make([]int, 0)
		for _, element := range strings.Split(line, "") {
			value, _ := strconv.Atoi(element)
			currentLine = append(currentLine, value)
		}
		output.depth = append(output.depth, currentLine)
	}
	return output
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func InBasin(depthA, depthB int) bool {
	if depthA == 9 {
		return false
	}

	if depthA >= depthB {
		return true
	}

	return false
}

func basinSize(depthMap HeightMap, basin LowPoint) int {
	basinMap := map[LowPoint]bool{}
	newPoints := []LowPoint{basin}
	for len(newPoints) > 0 {
		currentPoint := newPoints[0]
		currentDepth := depthMap.depth[currentPoint.y][currentPoint.x]
		newPoints = newPoints[1:]

		basinMap[currentPoint] = true

		// Check Left
		if currentPoint.x > 0 && InBasin(depthMap.depth[currentPoint.y][currentPoint.x-1], currentDepth) {
			_, ok := basinMap[LowPoint{currentPoint.y, currentPoint.x - 1}]
			if !ok {
				newPoints = append(newPoints, LowPoint{currentPoint.y, currentPoint.x - 1})
			}
		}

		// Check Right
		if currentPoint.x < len(depthMap.depth[currentPoint.y])-1 && InBasin(depthMap.depth[currentPoint.y][currentPoint.x+1], currentDepth) {
			_, ok := basinMap[LowPoint{currentPoint.y, currentPoint.x + 1}]
			if !ok {
				newPoints = append(newPoints, LowPoint{currentPoint.y, currentPoint.x + 1})
			}
		}

		// Check Up
		if currentPoint.y > 0 && InBasin(depthMap.depth[currentPoint.y-1][currentPoint.x], currentDepth) {
			_, ok := basinMap[LowPoint{currentPoint.y - 1, currentPoint.x}]
			if !ok {
				newPoints = append(newPoints, LowPoint{currentPoint.y - 1, currentPoint.x})
			}
		}

		// Check Down
		if currentPoint.y < len(depthMap.depth)-1 && InBasin(depthMap.depth[currentPoint.y+1][currentPoint.x], currentDepth) {
			_, ok := basinMap[LowPoint{currentPoint.y + 1, currentPoint.x}]
			if !ok {
				newPoints = append(newPoints, LowPoint{currentPoint.y + 1, currentPoint.x})
			}
		}
	}
	return len(basinMap)
}

func main() {
	flag.Parse()

	input := loadMapFromFile(*inputFile)
	lowPoints := make([]LowPoint, 0)

	for y := 0; y < len(input.depth); y++ {
		for x := 0; x < len(input.depth[y]); x++ {
			closestMaxHeight := 10

			// Check Up
			if y > 0 {
				closestMaxHeight = min(closestMaxHeight, input.depth[y-1][x])
			}

			// Check Down
			if y < len(input.depth)-1 {
				closestMaxHeight = min(closestMaxHeight, input.depth[y+1][x])
			}

			// Check Left
			if x > 0 {
				closestMaxHeight = min(closestMaxHeight, input.depth[y][x-1])
			}

			// Check Right
			if x < len(input.depth[y])-1 {
				closestMaxHeight = min(closestMaxHeight, input.depth[y][x+1])
			}

			if input.depth[y][x] < closestMaxHeight {
				lowPoints = append(lowPoints, LowPoint{y, x})
			}
		}
	}

	totalRisk := 0
	basinSizes := make([]int, 0)
	for _, lowPoint := range lowPoints {
		totalRisk += input.depth[lowPoint.y][lowPoint.x] + 1
		basinSizes = append(basinSizes, basinSize(input, lowPoint))

		fmt.Println("Low Point", lowPoint, "Risk", input.depth[lowPoint.y][lowPoint.x]+1, "Basin Size", basinSizes[len(basinSizes)-1])
	}
	fmt.Println("Total risk: ", totalRisk)
	sort.Sort(sort.Reverse(sort.IntSlice(basinSizes)))
	fmt.Println("Basin Sizes: ", basinSizes)
	fmt.Println("Top Three: ", basinSizes[0]*basinSizes[1]*basinSizes[2])
}
