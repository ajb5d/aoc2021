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
	inputFile   = flag.String("input", "data/test.txt", "Input file")
	verboseFlag = flag.Bool("verbose", false, "Verbose output")
)

const (
	HIGHLIGHT = "\u001b[31m"
	RESET     = "\u001b[0m"
)

type Point struct {
	x, y int
}

type Target struct {
	min, max Point
}

func loadMapFromFile(inputPath string) Target {
	output := Target{}
	fileHandle, err := os.Open(inputPath)

	if err != nil {
		panic(err)
	}
	defer fileHandle.Close()

	scanner := bufio.NewScanner(fileHandle)

	scanner.Scan()
	currentLine := scanner.Text()

	if !strings.HasPrefix(currentLine, "target area: ") {
		panic("Invalid input file")
	}

	currentLine = strings.TrimPrefix(currentLine, "target area: ")
	parts := strings.Split(currentLine, ", ")
	xParts := strings.Split(strings.TrimPrefix(parts[0], "x="), "..")
	yParts := strings.Split(strings.TrimPrefix(parts[1], "y="), "..")

	output.min.x, _ = strconv.Atoi(xParts[0])
	output.min.y, _ = strconv.Atoi(yParts[0])

	output.max.x, _ = strconv.Atoi(xParts[1])
	output.max.y, _ = strconv.Atoi(yParts[1])

	if output.max.x < output.min.x {
		output.max.x, output.min.x = output.min.x, output.max.x
	}

	if output.max.y < output.min.y {
		output.max.y, output.min.y = output.min.y, output.max.y
	}

	if *verboseFlag {
		fmt.Println(output)
	}

	return output
}

func (p Point) inTarget(target Target) bool {
	return target.min.x <= p.x && p.x <= target.max.x && target.min.y <= p.y && p.y <= target.max.y
}

func (p Point) inPath(path []Point) bool {
	for _, p2 := range path {
		if p.x == p2.x && p.y == p2.y {
			return true
		}
	}
	return false
}

func checkPath(velocity Point, target Target) (bool, int) {
	currentPoint := Point{0, 0}
	maxY := 0
	steps := 0
	path := make([]Point, 1)

	for {
		currentPoint.x += velocity.x
		currentPoint.y += velocity.y

		path = append(path, currentPoint)

		steps++

		if currentPoint.y > maxY {
			maxY = currentPoint.y
		}

		if velocity.x > 0 {
			velocity.x--
		} else if velocity.x < 0 {
			velocity.x++
		}

		velocity.y--

		if currentPoint.inTarget(target) {
			if *verboseFlag {
				printPath(path, target)
			}
			return true, maxY
		}

		if (velocity.y < 0 && currentPoint.y < target.min.y) || // If falling and below target
			(velocity.x == 0 && currentPoint.x < target.min.x) || // No horizontal movement and coming up short
			(currentPoint.x > target.max.x) { //Over the target
			if *verboseFlag {
				printPath(path, target)
			}
			return false, maxY
		}

		if steps > 10000 {
			panic("Infinite loop")
		}

	}
}

func printPath(path []Point, target Target) {
	min, max := target.min, target.max

	for _, p := range path {
		if p.x < min.x {
			min.x = p.x
		}

		if p.x > max.x {
			max.x = p.x
		}

		if p.y < min.y {
			min.y = p.y
		}

		if p.y > max.y {
			max.y = p.y
		}
	}

	for y := max.y; y >= min.y; y-- {
		for x := min.x; x <= max.x; x++ {
			p := Point{x, y}
			if p.inPath(path) {
				fmt.Print(HIGHLIGHT, "#", RESET)
			} else if p.inTarget(target) {
				fmt.Print("T")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()

	}
}

func main() {
	flag.Parse()
	input := loadMapFromFile(*inputFile)

	biggestY := 0
	hits := map[Point]bool{}

	for x := 1; x <= 1000; x++ {
		for y := -1000; y <= 1000; y++ {
			velocity := Point{x, y}
			hit, maxY := checkPath(velocity, input)

			if hit {
				if maxY > biggestY {
					biggestY = maxY
				}
				hits[velocity] = true
			}
		}
	}

	fmt.Println("Biggest Y:", biggestY)
	fmt.Println("Hits:", len(hits))
}
