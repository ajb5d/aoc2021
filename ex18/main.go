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

func loadMapFromFile(inputPath string) []string {
	output := []string{}
	fileHandle, err := os.Open(inputPath)

	if err != nil {
		panic(err)
	}
	defer fileHandle.Close()

	scanner := bufio.NewScanner(fileHandle)

	for scanner.Scan() {
		output = append(output, scanner.Text())
	}

	return output
}

type Number struct {
	leftPtr, rightPtr     *Number
	leftValue, rightValue int

	invalid bool
}

func parseNumber(str string) *Number {

	if !strings.Contains(str, "[") {
		value, _ := strconv.Atoi(str)
		return &Number{leftValue: value, invalid: true}
	}

	bracketDepth := 0
	for index, char := range str {
		switch char {
		case '[':
			bracketDepth++
		case ']':
			bracketDepth--
		case ',':
			if bracketDepth == 1 {
				partA := str[1:index]
				partB := str[index+1 : len(str)-1]
				leftPtr := parseNumber(partA)
				rightPtr := parseNumber(partB)
				return &Number{leftPtr: leftPtr, rightPtr: rightPtr}
			}
		}
	}
	panic("Unable to parse number")
}

func (node *Number) fixTree() {
	if node.leftPtr != nil {
		if node.leftPtr.invalid {
			node.leftValue = node.leftPtr.leftValue
			node.leftPtr = nil
		} else {
			node.leftPtr.fixTree()
		}
	}
	if node.rightPtr != nil {
		if node.rightPtr.invalid {
			node.rightValue = node.rightPtr.leftValue
			node.rightPtr = nil
		} else {
			node.rightPtr.fixTree()
		}
	}
}

func (node *Number) stringValue() string {
	output := "["
	if node.invalid {
		panic("Invalid node")
	}

	if node.leftPtr != nil {
		output += node.leftPtr.stringValue()
	} else {
		output += fmt.Sprintf("%d", node.leftValue)
	}

	output += ","

	if node.rightPtr != nil {
		output += node.rightPtr.stringValue()
	} else {
		output += fmt.Sprintf("%d", node.rightValue)
	}

	return output + "]"
}

// func max(a, b int) int {
// 	if a > b {
// 		return a
// 	}
// 	return b
// }

// func (node *Number) depth() int {
// 	leftDepth, rightDepth := 0, 0
// 	if node.leftPtr != nil {
// 		leftDepth = node.leftPtr.depth()
// 	}

// 	if node.rightPtr != nil {
// 		rightDepth = node.rightPtr.depth()
// 	}

// 	return max(leftDepth, rightDepth) + 1
// }

func (node *Number) depthFirstWalk() []*int {
	output := make([]*int, 0)
	if node.leftPtr != nil {
		output = append(output, node.leftPtr.depthFirstWalk()...)
	} else {
		output = append(output, &node.leftValue)
	}

	if node.rightPtr != nil {
		output = append(output, node.rightPtr.depthFirstWalk()...)
	} else {
		output = append(output, &node.rightValue)
	}
	return output
}

func (node *Number) findPairThatWillExplode(currentDepth int) *Number {
	if currentDepth >= 4 && node.leftPtr == nil && node.rightPtr == nil {
		return node
	}

	if node.leftPtr != nil {
		leftSide := node.leftPtr.findPairThatWillExplode(currentDepth + 1)
		if leftSide != nil {
			return leftSide
		}
	}

	if node.rightPtr != nil {
		rightSide := node.rightPtr.findPairThatWillExplode(currentDepth + 1)
		if rightSide != nil {
			return rightSide
		}
	}

	return nil
}

func (node *Number) findNumberThatWillSplit() *int {

	if node.leftPtr != nil {
		leftSide := node.leftPtr.findNumberThatWillSplit()
		if leftSide != nil {
			return leftSide
		}
	} else {
		if node.leftValue > 9 {
			return &node.leftValue
		}
	}

	if node.rightPtr != nil {
		rightSide := node.rightPtr.findNumberThatWillSplit()
		if rightSide != nil {
			return rightSide
		}
	} else {
		if node.rightValue > 9 {
			return &node.rightValue
		}
	}

	return nil
}

func (node *Number) findAndSplitTarget(target *int) bool {
	if node.leftPtr == nil {
		if &node.leftValue == target {
			node.leftPtr = &Number{}
			node.leftPtr.leftValue = *target / 2
			node.leftPtr.rightValue = *target/2 + (*target % 2)
			return true
		}
	} else {
		if node.leftPtr.findAndSplitTarget(target) {
			return true
		}
	}

	if node.rightPtr == nil {
		if &node.rightValue == target {
			node.rightPtr = &Number{}
			node.rightPtr.leftValue = *target / 2
			node.rightPtr.rightValue = *target/2 + (*target % 2)
			return true
		}
	} else {
		if node.rightPtr.findAndSplitTarget(target) {
			return true
		}
	}
	return false
}

func (node *Number) replaceWithZero(target *Number) {
	if node.leftPtr == target {
		node.leftValue = 0
		node.leftPtr = nil
		return
	}

	if node.rightPtr == target {
		node.rightValue = 0
		node.rightPtr = nil
		return
	}

	if node.leftPtr != nil {
		node.leftPtr.replaceWithZero(target)
	}

	if node.rightPtr != nil {
		node.rightPtr.replaceWithZero(target)
	}
}

func (node *Number) reduce() {
	for {
		explodeTarget := node.findPairThatWillExplode(0)
		splitTarget := node.findNumberThatWillSplit()

		if explodeTarget != nil {
			if *verboseFlag {
				fmt.Printf("Pre explode   %s\n", node.stringValue())
			}
			values := node.depthFirstWalk()
			for index, valuePtr := range values {
				if valuePtr == &explodeTarget.leftValue {
					if index > 0 {
						*(values[index-1]) += *valuePtr
					}
				}
				if valuePtr == &explodeTarget.rightValue {
					if index < len(values)-1 {
						*(values[index+1]) += *valuePtr
					}
				}
			}
			node.replaceWithZero(explodeTarget)
			if *verboseFlag {
				fmt.Printf("After explode %s\n\n", node.stringValue())
			}
			continue
		}

		if splitTarget != nil {
			if *verboseFlag {
				fmt.Printf("Pre split     %s\n", node.stringValue())
			}
			result := node.findAndSplitTarget(splitTarget)
			if !result {
				panic("Failed?")
			}
			if *verboseFlag {
				fmt.Printf("After split   %s\n\n", node.stringValue())
			}
			continue
		}

		if explodeTarget == nil && splitTarget == nil {
			return
		}
	}
}

func (node *Number) magnitude() int {
	leftValue, rightValue := node.leftValue, node.rightValue

	if node.leftPtr != nil {
		leftValue = node.leftPtr.magnitude()
	}

	if node.rightPtr != nil {
		rightValue = node.rightPtr.magnitude()
	}

	return 3*leftValue + 2*rightValue
}

func part1(input []string) {
	value := parseNumber(input[0])
	value.fixTree()

	for _, line := range input[1:] {
		new := parseNumber(line)
		new.fixTree()

		newValue := &Number{leftPtr: value, rightPtr: new}
		value = newValue
		value.reduce()
	}

	finalValue := value.stringValue()

	fmt.Println()
	fmt.Println(finalValue)
	fmt.Println(value.magnitude())
}

func part2(input []string) {
	maxValue := 0

	for _, value1 := range input {
		for _, value2 := range input {
			if value1 == value2 {
				continue
			}

			tree1 := parseNumber(value1)
			tree2 := parseNumber(value2)

			tree1.fixTree()
			tree2.fixTree()

			product := Number{leftPtr: tree1, rightPtr: tree2}
			product.reduce()

			productValue := product.magnitude()

			if productValue > maxValue {
				maxValue = productValue
			}

		}
	}

	fmt.Println(maxValue)
}

func main() {
	flag.Parse()

	input := loadMapFromFile(*inputFile)
	// part1(input)
	part2(input)

}
