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

type BingoGame struct {
	NumberSequence []int
	BingoCards     []BingoCard
}

type BingoCard struct {
	Numbers      [5][5]int
	Marked       [5][5]bool
	WinningCard  bool
	WinSequence  int
	WinningValue int
}

func IngestBingoFile(inputPath string) BingoGame {
	output := BingoGame{}
	fileHandle, err := os.Open(inputPath)

	if err != nil {
		fmt.Println(err)
		return output
	}

	defer fileHandle.Close()

	scanner := bufio.NewScanner(fileHandle)

	// Read in the first line with the sequence of numbers
	scanner.Scan()
	firstLine := scanner.Text()

	for _, element := range strings.Split(firstLine, ",") {
		intValue, _ := strconv.Atoi(element)
		output.NumberSequence = append(output.NumberSequence, intValue)
	}

	for scanner.Scan() {
		currentCard := BingoCard{}
		for i := 0; i < 5; i++ {
			scanner.Scan()
			currentLine := scanner.Text()
			for j, element := range strings.Fields(currentLine) {
				intValue, _ := strconv.Atoi(element)
				currentCard.Numbers[i][j] = intValue
			}
		}
		output.BingoCards = append(output.BingoCards, currentCard)
	}

	return output
}

func CheckBingoCard(card BingoCard) bool {
	for i := 0; i < 5; i++ {
		if card.Marked[i][0] && card.Marked[i][1] && card.Marked[i][2] && card.Marked[i][3] && card.Marked[i][4] {
			return true
		}
	}

	for i := 0; i < 5; i++ {
		if card.Marked[0][i] && card.Marked[1][i] && card.Marked[2][i] && card.Marked[3][i] && card.Marked[4][i] {
			return true
		}
	}

	return false
}

func PrintBingoCard(card BingoCard) {
	fmt.Println("+-----+-----+-----+-----+-----+")
	for i := 0; i < 5; i++ {
		fmt.Printf("|")
		for j := 0; j < 5; j++ {
			if card.Marked[i][j] {
				fmt.Printf(" *%2d |", card.Numbers[i][j])
			} else {
				fmt.Printf("  %2d |", card.Numbers[i][j])
			}
		}
		fmt.Printf("\n")
		fmt.Println("+-----+-----+-----+-----+-----+")
	}
}

func ScoreBingoCard(card BingoCard) int {
	sum := 0
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if !card.Marked[i][j] {
				sum += card.Numbers[i][j]
			}
		}
	}
	return sum * card.WinningValue
}

func FindWinningCard(input BingoGame, position int) BingoCard {
	for _, card := range input.BingoCards {
		if card.WinningCard && card.WinSequence == position {
			return card
		}
	}
	return BingoCard{}
}

func main() {
	flag.Parse()

	input := IngestBingoFile(*inputFile)

	currentWinningCardCount := 0

	for index, value := range input.NumberSequence {
		fmt.Printf("%d: %d\n", index, value)

		for cardIndex, card := range input.BingoCards {
			if card.WinningCard {
				continue
			}
			for i := 0; i < 5; i++ {
				for j := 0; j < 5; j++ {
					if card.Numbers[i][j] == value {
						input.BingoCards[cardIndex].Marked[i][j] = true
					}
				}
			}

			if CheckBingoCard(input.BingoCards[cardIndex]) {
				if !input.BingoCards[cardIndex].WinningCard {
					input.BingoCards[cardIndex].WinningValue = value
					input.BingoCards[cardIndex].WinningCard = true
					input.BingoCards[cardIndex].WinSequence = currentWinningCardCount
					currentWinningCardCount++
				}
			}
		}
	}

	firstCard := FindWinningCard(input, 0)
	fmt.Println("First winning card:", ScoreBingoCard(firstCard))
	PrintBingoCard(firstCard)

	lastCard := FindWinningCard(input, currentWinningCardCount-1)
	fmt.Println("Last winning card:", ScoreBingoCard(lastCard))
	PrintBingoCard(lastCard)

}
