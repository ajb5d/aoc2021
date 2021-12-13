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

type Map struct {
	values     map[Point]int
	maxX, maxY int
}

func loadMapFromFile(inputPath string) Map {
	output := Map{}
	output.values = make(map[Point]int)

	fileHandle, err := os.Open(inputPath)

	if err != nil {
		panic(err)
	}
	defer fileHandle.Close()

	scanner := bufio.NewScanner(fileHandle)
	row := 0
	for scanner.Scan() {
		col := 0
		for _, element := range strings.Split(scanner.Text(), "") {
			value, _ := strconv.Atoi(element)
			output.values[Point{row, col}] = value
			col++
		}
		if output.maxY == 0 {
			output.maxY = col
		} else {
			if output.maxY != col {
				panic("Inconsistent row length")
			}
		}
		row++
	}

	output.maxX = row
	return output
}

func printMap(x Map) {
	for i := 0; i < x.maxX; i++ {
		for j := 0; j < x.maxY; j++ {
			fmt.Printf(" %2d ", x.values[Point{i, j}])
		}
		fmt.Println()
	}
}

func getNeighbors(p Point, source Map) []Point {
	output := []Point{}
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			if x == 0 && y == 0 {
				continue
			}
			if p.x+x >= 0 && p.x+x < source.maxX && p.y+y >= 0 && p.y+y < source.maxY {
				output = append(output, Point{p.x + x, p.y + y})
			}
		}
	}
	return output
}

func main() {
	flag.Parse()

	input := loadMapFromFile(*inputFile)

	flashCount := 0
	for cycle := 0; cycle < 1000; cycle++ {
		didFlash := map[Point]bool{}
		willFlash := map[Point]bool{}

		for k := range input.values {
			input.values[k] += 1
			if input.values[k] > 9 {
				willFlash[k] = true
			}
		}

		for len(willFlash) > 0 {
			for k := range willFlash {
				didFlash[k] = true
				delete(willFlash, k)

				for _, neighbor := range getNeighbors(k, input) {
					input.values[neighbor] += 1
					if input.values[neighbor] > 9 {
						_, ok := didFlash[neighbor]
						if !ok {
							willFlash[neighbor] = true
						}
					}
				}
			}
		}

		for k := range didFlash {
			input.values[k] = 0
		}

		flashCount += len(didFlash)
		if cycle == 99 {
			fmt.Println("Flash count after 100:", flashCount)
		}

		if len(didFlash) == len(input.values) {
			fmt.Println("All Flash on cycle:", cycle+1)
			break
		}
	}

}
