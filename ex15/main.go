package main

import (
	"bufio"
	"container/heap"
	"flag"
	"fmt"
	"os"
)

var (
	inputFile = flag.String("input", "", "Input file")
)

const (
	HIGHLIGHT  = "\u001b[31m"
	RESET      = "\u001b[0m"
	REALLY_BIG = 999999999
)

type Point struct {
	row, col int
}

func loadMapFromFile(inputPath string) [][]int {
	output := make([][]int, 0)
	fileHandle, err := os.Open(inputPath)

	if err != nil {
		panic(err)
	}
	defer fileHandle.Close()

	scanner := bufio.NewScanner(fileHandle)

	for scanner.Scan() {
		currentLine := scanner.Text()
		output = append(output, make([]int, len(currentLine)))
		for i, c := range currentLine {
			output[len(output)-1][i] = int(c) - 48
		}
	}

	return output
}

func (p Point) up() Point {
	return Point{row: p.row - 1, col: p.col}
}

func (p Point) down() Point {
	return Point{row: p.row + 1, col: p.col}
}

func (p Point) left() Point {
	return Point{row: p.row, col: p.col - 1}
}

func (p Point) right() Point {
	return Point{row: p.row, col: p.col + 1}
}

func valueForPoint(point Point, input [][]int) int {
	mappedPoint := Point{row: point.row, col: point.col}
	costAdjustment := 0
	extent := len(input)

	for mappedPoint.row >= extent {
		mappedPoint.row -= extent
		costAdjustment += 1
	}

	for mappedPoint.col >= extent {
		mappedPoint.col -= extent
		costAdjustment += 1
	}

	cost := input[mappedPoint.row][mappedPoint.col] + costAdjustment
	if cost > 9 {
		cost = cost%10 + 1
	}
	return cost
}

type HeapElement struct {
	value    Point
	priority int
	index    int
}

type PriorityPointQueue []*HeapElement

func (p PriorityPointQueue) Len() int {
	return len(p)
}

func (pq PriorityPointQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityPointQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityPointQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*HeapElement)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityPointQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func (pq PriorityPointQueue) findInHeap(value Point) (int, bool) {
	for index := range pq {
		if pq[index].value == value {
			return index, true
		}
	}
	return -1, false
}

func pathSearch(input [][]int, extent int, progress bool) int {
	/**
	* Implementation of Dijkstra's algorithm to find the shortest path from 0,0 to extent-1,extent-1
	**/

	distance := map[Point]int{}
	previous := map[Point]Point{}
	toVisit := make(PriorityPointQueue, extent*extent)

	for i := 0; i < extent; i++ {
		for j := 0; j < extent; j++ {
			distance[Point{row: i, col: j}] = REALLY_BIG
			toVisit[i*extent+j] = &HeapElement{value: Point{row: i, col: j}, priority: REALLY_BIG}
		}
	}

	distance[Point{row: 0, col: 0}] = 0
	toVisit[0] = &HeapElement{value: Point{row: 0, col: 0}, priority: 0}
	heap.Init(&toVisit)

	for len(toVisit) > 0 {
		if progress && len(toVisit)%1000 == 0 {
			fmt.Println("To Visit:", len(toVisit))
		}
		currentNode := heap.Pop(&toVisit).(*HeapElement).value

		if currentNode.row == extent-1 && currentNode.col == extent-1 {
			break
		}

		for _, neighbor := range []Point{currentNode.up(), currentNode.down(), currentNode.left(), currentNode.right()} {
			index, inHeap := toVisit.findInHeap(neighbor)
			if inHeap {
				alternateDistance := distance[currentNode] + valueForPoint(neighbor, input)
				if alternateDistance < distance[neighbor] {
					distance[neighbor] = alternateDistance
					previous[neighbor] = currentNode
					toVisit[index].priority = alternateDistance
					heap.Fix(&toVisit, index)
				}
			}
		}
		// heap.Init(&toVisit)
	}

	return distance[Point{row: extent - 1, col: extent - 1}]
}

func part1(input [][]int) {
	extent := len(input)
	fmt.Println("Part 1: ", pathSearch(input, extent, false))
}

func part2(input [][]int) {
	extent := len(input) * 5
	fmt.Println("Part 2: ", pathSearch(input, extent, true))
}

func main() {
	flag.Parse()
	input := loadMapFromFile(*inputFile)
	part1(input)
	part2(input)
}
