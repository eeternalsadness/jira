package util

import (
	"bufio"
	"fmt"
	"os"
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
