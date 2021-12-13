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

type CaveMap struct {
	nodes map[string][]string
}

func loadMapFromFile(inputPath string) CaveMap {
	output := CaveMap{}
	output.nodes = make(map[string][]string)

	fileHandle, err := os.Open(inputPath)

	if err != nil {
		panic(err)
	}
	defer fileHandle.Close()

	scanner := bufio.NewScanner(fileHandle)
	for scanner.Scan() {
		elements := strings.Split(scanner.Text(), "-")

		output.nodes[elements[0]] = append(output.nodes[elements[0]], elements[1])
		output.nodes[elements[1]] = append(output.nodes[elements[1]], elements[0])

	}

	return output
}

func nodeIsLarge(node string) bool {
	return node == strings.ToUpper(node)
}

func pathLegal(newPath []string) bool {
	visitedNodes := map[string]int{}

	for _, element := range newPath {
		visitedNodes[element] += 1
		repeatVisitCount := 0

		for key, value := range visitedNodes {
			if value > 1 && (key == "start" || key == "end") {
				return false
			}

			if !nodeIsLarge(key) {
				if value > 2 {
					return false
				}

				if value > 1 {
					repeatVisitCount += 1
				}
			}

		}

		if repeatVisitCount > 1 {
			return false
		}
	}

	return true
}

func pathName(path []string) string {
	return strings.Join(path, "-")
}

func main() {
	flag.Parse()

	input := loadMapFromFile(*inputFile)

	currentPaths := make([][]string, 0)
	nextPaths := make([][]string, 0)
	completedPaths := map[string]bool{}
	searchedPaths := map[string]bool{}

	for _, value := range input.nodes["start"] {
		newPath := make([]string, 2)
		newPath[0] = "start"
		newPath[1] = value

		currentPaths = append(currentPaths, newPath)
	}

	for len(currentPaths) > 0 {
		// fmt.Println()
		// fmt.Println("Paths to Explore: ", currentPaths)
		// fmt.Println(len(currentPaths))

		for _, path := range currentPaths {
			// fmt.Println("Current path to evaluate: ", path)
			endNode := path[len(path)-1]

			if endNode == "end" {
				completedPaths[pathName(path)] = true
				continue
			}
			// fmt.Println("End node: ", endNode)
			// fmt.Println("Neighbors: ", input.nodes[endNode])

			for _, value := range input.nodes[endNode] {
				newPath := make([]string, len(path)+1)
				copy(newPath, path)
				newPath[len(newPath)-1] = value

				if _, ok := searchedPaths[pathName(newPath)]; !ok {
					searchedPaths[pathName(newPath)] = true
					// fmt.Println("Candidate Path", newPath, "legal?", pathLegal(newPath))
					if pathLegal(newPath) {
						nextPaths = append(nextPaths, newPath)
					}

				}
			}
			// fmt.Println("Next Paths: ", nextPaths)
			// fmt.Println()
		}

		currentPaths = nextPaths
		nextPaths = make([][]string, 0)
	}

	// fmt.Println("Completed paths: ", completedPaths)
	fmt.Println("Number of paths: ", len(completedPaths))

}
