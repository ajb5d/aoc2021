package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
)

var (
	inputFile           = flag.String("input", "data/test.txt", "Input file")
	verboseFlag         = flag.Bool("verbose", false, "Verbose output")
	TOTAL_VERSION_COUNT = 0
)

const (
	HIGHLIGHT = "\u001b[31m"
	RESET     = "\u001b[0m"
)

type BitStream struct {
	data    []byte
	bytePos int
	bitPos  int
}

func loadMapFromFile(inputPath string) BitStream {
	output := BitStream{data: make([]byte, 0), bytePos: 0, bitPos: 0}
	fileHandle, err := os.Open(inputPath)

	if err != nil {
		panic(err)
	}
	defer fileHandle.Close()

	scanner := bufio.NewScanner(fileHandle)

	scanner.Scan()
	currentLine := scanner.Text()

	for i := 0; i < len(currentLine); i += 2 {
		value, _ := strconv.ParseUint(string(currentLine[i:i+2]), 16, 8)
		output.data = append(output.data, byte(value))
	}
	if *verboseFlag {
		fmt.Println("Input line:", currentLine)
		fmt.Println("Total Bits In:", len(output.data)*8)
		output.printData()
	}

	return output
}

func (b BitStream) printData() {
	for i := range b.data {
		for j := 7; j >= 0; j-- {
			if b.data[i]&(1<<uint(j)) != 0 {
				fmt.Print(HIGHLIGHT, "1", RESET)
			} else {
				fmt.Print("0")
			}
		}
		fmt.Printf(" ")
	}
	fmt.Println()
}

func (b *BitStream) takeBits(numBits int) uint64 {
	const MostSignificantBit = byte(1 << 7)
	output := uint64(0)
	for i := 0; i < numBits; i++ {
		output <<= 1
		if b.bitPos == 8 {
			b.bitPos = 0
			b.bytePos++
		}
		if b.bytePos >= len(b.data) {
			panic("Ran out of data")
		}
		if b.data[b.bytePos]&(MostSignificantBit>>uint(b.bitPos)) != 0 {
			output |= 1
		}
		b.bitPos++
	}
	return output
}

func readPacket(b *BitStream, depth int) (int, int) {
	prefixString := ""
	for i := 0; i < depth; i++ {
		prefixString += "  "
	}

	// Read Version String
	version := b.takeBits(3)
	typeId := b.takeBits(3)
	totalBits := 6

	TOTAL_VERSION_COUNT += int(version)

	if *verboseFlag {
		fmt.Printf("%sVersion: %d, Type: %d\n", prefixString, version, typeId)
	}

	if typeId == 4 {
		payload, bitsRead := readLiteralValue(b)
		if *verboseFlag {
			fmt.Printf("%s Literal Value: %d\n", prefixString, payload)
		}
		return int(payload), bitsRead + totalBits
	}

	lengthType := b.takeBits(1)
	totalBits++

	subPacketValues := make([]int, 0)

	if lengthType == 0 {
		packetLength := b.takeBits(15)
		totalBits += 15

		if *verboseFlag {
			fmt.Printf("%s Length (DirectBits): %d\n", prefixString, packetLength)
		}

		subtotalBits := 0
		for subtotalBits < int(packetLength) {
			value, bitsRead := readPacket(b, depth+1)
			subPacketValues = append(subPacketValues, value)
			subtotalBits += bitsRead
		}
		totalBits += subtotalBits
	} else {
		packetsToRead := b.takeBits(11)
		totalBits += 11

		if *verboseFlag {
			fmt.Printf("%s Length (Packets): %d\n", prefixString, packetsToRead)
		}

		subtotalBits := 0
		for i := 0; i < int(packetsToRead); i++ {
			value, bitsRead := readPacket(b, depth+1)
			subPacketValues = append(subPacketValues, value)
			subtotalBits += bitsRead
		}
		totalBits += subtotalBits
	}

	switch typeId {
	case 0:
		// Sum
		totalValue := 0
		for _, value := range subPacketValues {
			totalValue += value
		}

		if *verboseFlag {
			fmt.Printf("%s Sum: %d\n", prefixString, totalValue)
		}
		return totalValue, totalBits

	case 1:
		// Product
		totalValue := 1
		for _, value := range subPacketValues {
			totalValue *= value
		}

		if *verboseFlag {
			fmt.Printf("%s Product: %d\n", prefixString, totalValue)
		}
		return totalValue, totalBits

	case 2:
		// Min
		minValue := subPacketValues[0]
		for _, value := range subPacketValues {
			if value < minValue {
				minValue = value
			}
		}

		if *verboseFlag {
			fmt.Printf("%s Min: %d\n", prefixString, minValue)
		}
		return minValue, totalBits

	case 3:
		// Max
		maxValue := subPacketValues[0]
		for _, value := range subPacketValues {
			if value > maxValue {
				maxValue = value
			}
		}

		if *verboseFlag {
			fmt.Printf("%s Max: %d\n", prefixString, maxValue)
		}

		return maxValue, totalBits

	case 5:
		// Greather Than
		if len(subPacketValues) != 2 {
			panic("Greater Than needs 2 values")
		}

		outputValue := 0

		if subPacketValues[0] > subPacketValues[1] {
			outputValue = 1
		}

		if *verboseFlag {
			fmt.Printf("%s Greater Than: %d\n", prefixString, outputValue)
		}

		return outputValue, totalBits

	case 6:
		// Less Than
		if len(subPacketValues) != 2 {
			panic("Less Than needs 2 values")
		}

		outputValue := 0

		if subPacketValues[0] < subPacketValues[1] {
			outputValue = 1
		}

		if *verboseFlag {
			fmt.Printf("%s Less Than: %d\n", prefixString, outputValue)
		}

		return outputValue, totalBits

	case 7:
		// Equal
		if len(subPacketValues) != 2 {
			panic("Equal needs 2 values")
		}

		outputValue := 0

		if subPacketValues[0] == subPacketValues[1] {
			outputValue = 1
		}

		if *verboseFlag {
			fmt.Printf("%s Equal: %d\n", prefixString, outputValue)
		}

		return outputValue, totalBits

	default:
		panic("Unknown type")
	}
}

func readLiteralValue(b *BitStream) (uint64, int) {
	var output uint64
	var bitsRead int
	for {
		currentNibble := b.takeBits(5)
		bitsRead += 5

		output <<= 4
		output |= (currentNibble & 0x0F)

		if currentNibble&0x10 == 0 {
			break
		}
	}

	return output, bitsRead
}

func part1(input BitStream) int {
	TOTAL_VERSION_COUNT = 0
	_, totalBits := readPacket(&input, 0)
	if *verboseFlag {
		fmt.Println("Total bits read:", totalBits)
	}
	return TOTAL_VERSION_COUNT
}

func part2(input BitStream) int {
	value, _ := readPacket(&input, 0)
	return value
}

func main() {
	flag.Parse()
	input := loadMapFromFile(*inputFile)

	fmt.Println("Part 1:", part1(input))
	fmt.Println("Part 2:", part2(input))
}
