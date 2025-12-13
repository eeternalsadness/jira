package util

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var ErrUserQuit = errors.New("")

// NOTE: this is basically UserGetBool(), but standardized
func UserYesNo(prompt string) (bool, error) {
	userInput, err := UserGetString(
		fmt.Sprintf("%s [y/n]: ", prompt),
		nil,
		false)
	if err != nil {
		return false, err
	}

	switch *userInput {
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

	index, err := UserGetInt(
		fmt.Sprintf("\nSelect an option [1 - %d, or 'q' to quit]: ", max),
		nil,
		true)
	if err != nil {
		return -1, err
	}

	// check if index value is valid
	if *index < 1 || *index > max {
		return -1, fmt.Errorf("you must choose a number between 1 and %d (inclusive)", max)
	}

	return *index - 1, nil
}

func UserGetInt(prompt string, defaultVal *int, hasQuitOption bool) (*int, error) {
	userInput, err := UserGetString(prompt, nil, hasQuitOption)
	if err != nil {
		return nil, err
	}

	// empty user input, use defaultVal
	if *userInput == "" && defaultVal != nil {
		return defaultVal, nil
	}

	if userInputInt, err := strconv.Atoi(*userInput); err != nil {
		return nil, err
	} else {
		return &userInputInt, nil
	}
}

func UserGetString(prompt string, defaultVal *string, hasQuitOption bool) (*string, error) {
	fmt.Print(prompt)

	reader := bufio.NewReader(os.Stdin)
	userInput, err := reader.ReadString('\n')
	if err != nil && errors.Is(err, io.EOF) {
		return nil, err
	}
	userInput = strings.TrimSpace(userInput)

	// user quits
	if hasQuitOption && userInput == "q" {
		return nil, ErrUserQuit
	}

	if userInput == "" && defaultVal != nil {
		return defaultVal, nil
	}

	return &userInput, nil
}
