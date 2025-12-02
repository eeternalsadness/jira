package util

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func UserYesNo(prompt string) (bool, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("%s [y/n]: ", prompt)
	userInput, _ := reader.ReadString('\n')
	userInput = userInput[:len(userInput)-1]

	switch userInput {
	case "y":
		return true, nil
	case "n":
		return false, nil
	default:
		return false, fmt.Errorf("input must be 'y' or 'n'")
	}
}

func UserSelectFromRange(max int) (int, error) {
	if max < 1 {
		panic("max must be >= 1!")
	}

	// prompt for input
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\nSelect an option [1 - %d, or 'q' to quit]: ", max)
	inputStr, err := reader.ReadString('\n')
	if err != nil {
		return -1, err
	}
	inputStr = inputStr[:len(inputStr)-1]

	// quit if user types 'q'
	if inputStr == "q" {
		return -1, nil
	}

	// check number value
	index, err := strconv.Atoi(inputStr)
	if err != nil {
		return -1, err
	}

	// check if index value is valid
	if index <= 0 || index > max {
		return -1, fmt.Errorf("you must choose a number between 1 and %d (inclusive)", max)
	}

	return index - 1, nil
}
