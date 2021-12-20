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

var (
	BASES = [][3]Point{
		{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}},
		{{1, 0, 0}, {0, -1, 0}, {0, 0, -1}},
		{{1, 0, 0}, {0, 0, 1}, {0, -1, 0}},
		{{1, 0, 0}, {0, 0, -1}, {0, 1, 0}},

		{{-1, 0, 0}, {0, 1, 0}, {0, 0, -1}},
		{{-1, 0, 0}, {0, -1, 0}, {0, 0, 1}},
		{{-1, 0, 0}, {0, 0, 1}, {0, 1, 0}},
		{{-1, 0, 0}, {0, 0, -1}, {0, -1, 0}},

		{{0, 1, 0}, {1, 0, 0}, {0, 0, -1}},
		{{0, 1, 0}, {-1, 0, 0}, {0, 0, 1}},
		{{0, 1, 0}, {0, 0, 1}, {1, 0, 0}},
		{{0, 1, 0}, {0, 0, -1}, {-1, 0, 0}},

		{{0, -1, 0}, {1, 0, 0}, {0, 0, 1}},
		{{0, -1, 0}, {-1, 0, 0}, {0, 0, -1}},
		{{0, -1, 0}, {0, 0, 1}, {-1, 0, 0}},
		{{0, -1, 0}, {0, 0, -1}, {1, 0, 0}},

		{{0, 0, 1}, {1, 0, 0}, {0, 1, 0}},
		{{0, 0, 1}, {-1, 0, 0}, {0, -1, 0}},
		{{0, 0, 1}, {0, 1, 0}, {-1, 0, 0}},
		{{0, 0, 1}, {0, -1, 0}, {1, 0, 0}},

		{{0, 0, -1}, {1, 0, 0}, {0, -1, 0}},
		{{0, 0, -1}, {-1, 0, 0}, {0, 1, 0}},
		{{0, 0, -1}, {0, 1, 0}, {1, 0, 0}},
		{{0, 0, -1}, {0, -1, 0}, {-1, 0, 0}},
	}
)

const (
	HIGHLIGHT = "\u001b[31m"
	RESET     = "\u001b[0m"
)

type Point struct {
	x, y, z int
}

func (p Point) String() string {
	return fmt.Sprintf("(%d, %d, %d)", p.x, p.y, p.z)
}

func (p Point) equal(other Point) bool {
	return p.x == other.x && p.y == other.y && p.z == other.z
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func (p Point) manhattanDistance(other Point) int {
	return abs(p.x-other.x) + abs(p.y-other.y) + abs(p.z-other.z)
}

func (p Point) project(rotation int) Point {

	a, b, c := BASES[rotation][0].x, BASES[rotation][0].y, BASES[rotation][0].z
	d, e, f := BASES[rotation][1].x, BASES[rotation][1].y, BASES[rotation][1].z
	g, h, i := BASES[rotation][2].x, BASES[rotation][2].y, BASES[rotation][2].z

	x, y, z := p.x, p.y, p.z

	return Point{
		x: a*x + b*y + c*z,
		y: d*x + e*y + f*z,
		z: g*x + h*y + i*z}
}

func (p Point) offset(other Point) Point {
	return Point{
		x: other.x - p.x,
		y: other.y - p.y,
		z: other.z - p.z,
	}
}

func (p Point) add(other Point) Point {
	return Point{
		x: p.x + other.x,
		y: p.y + other.y,
		z: p.z + other.z,
	}
}

type Scanner struct {
	observations []Point
	position     Point
	rotation     int
	name         string
	located      bool
}

func (s Scanner) hasObservation(p Point) bool {
	for _, o := range s.observations {
		if o.equal(p) {
			return true
		}
	}
	return false
}

func loadMapFromFile(inputPath string) []Scanner {
	output := make([]Scanner, 1)
	index := 0
	new := true

	fileHandle, err := os.Open(inputPath)

	if err != nil {
		panic(err)
	}
	defer fileHandle.Close()

	scanner := bufio.NewScanner(fileHandle)

	for scanner.Scan() {
		currentLine := scanner.Text()
		currentLine = strings.TrimSpace(currentLine)

		if currentLine == "" {
			index++
			output = append(output, Scanner{})
			new = true
			continue
		}

		if new {
			output[index].name = currentLine
			new = false
			continue
		}

		parts := strings.Split(currentLine, ",")
		newPoint := Point{}
		newPoint.x, _ = strconv.Atoi(parts[0])
		newPoint.y, _ = strconv.Atoi(parts[1])
		newPoint.z, _ = strconv.Atoi(parts[2])

		output[index].observations = append(output[index].observations, newPoint)
	}

	return output
}

func (s Scanner) String() string {
	const fmtStr = "Scanner '%s' \n" +
		"\tPosition: %s\n" + "\tRotation: %d\n" + "\tObservations: %d\n"

	return fmt.Sprintf(fmtStr, s.name, s.position, s.rotation, len(s.observations))
}

func (s Scanner) projectAllPoints(basis int) []Point {
	output := make([]Point, 0)

	for _, p := range s.observations {
		output = append(output, p.project(basis))
	}

	return output
}

func (s *Scanner) move(offset Point, rotation int) {
	for i := range s.observations {
		s.observations[i] = s.observations[i].project(rotation).add(offset)
	}
	s.located = true
}

func (scanner *Scanner) calculateOffset(reference Scanner) (Point, int, bool) {
	for newBasis := range BASES {
		// Project all the points into the new orientation
		projectedPoints := scanner.projectAllPoints(newBasis)

		for i := 0; i < len(reference.observations); i++ {
			for j := 0; j < len(scanner.observations); j++ {
				// Find the offset between the two points
				offset := projectedPoints[j].offset(reference.observations[i])
				// Count the number of projectedPoints + offset that collide with the reference
				collisions := 0
				for _, p := range projectedPoints {
					if reference.hasObservation(p.add(offset)) {
						collisions++
					}
				}
				if collisions >= 12 {
					fmt.Printf("Using basis %d and offset %s found %d collisions\n", newBasis, offset, collisions)
					return offset, newBasis, true
				}

			}

		}
	}

	return Point{}, 0, false
}

func main() {
	flag.Parse()

	input := loadMapFromFile(*inputFile)

	input[0].located = true

	for {
		unlocated := 0
		for _, scanner := range input {
			if !scanner.located {
				unlocated++
			}
		}

		if unlocated == 0 {
			break
		}

		for i := range input {
			if !input[i].located {
				continue
			}
			for j := range input {
				if input[j].located {
					continue
				}

				offset, rotation, ok := input[j].calculateOffset(input[i])
				if ok {
					input[j].move(offset, rotation)
					input[j].position = offset
					fmt.Printf("%s%s%s\n", HIGHLIGHT, input[j], RESET)

				}
			}
		}
	}

	pointMap := map[Point]bool{}

	for _, scanner := range input {
		for _, p := range scanner.projectAllPoints(scanner.rotation) {
			pointMap[p] = true
		}
	}
	fmt.Println(len(pointMap))

	maxDistance := 0
	for i := range input {
		for j := range input {
			distance := input[j].position.manhattanDistance(input[i].position)
			if distance > maxDistance {
				maxDistance = distance
			}
		}
	}
	fmt.Println(maxDistance)

}
