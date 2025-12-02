package util

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"text/tabwriter"
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

func UserSelectFromRange[T any](headerMap map[string]string, options []T) (int, error) {
	if len(options) == 0 {
		return -1, fmt.Errorf("expected non-empty slice, got empty slice")
	}

	// make sure options is a struct slice
	if t := reflect.TypeOf(options[0]).Kind(); t != reflect.Struct {
		return -1, fmt.Errorf("expected a struct slice, got %s", t.String())
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	var builder strings.Builder

	// print out headers
	builder.WriteString("#\t")
	for header := range headerMap {
		_, err := builder.WriteString(fmt.Sprintf("%s\t", header))
		if err != nil {
			return -1, err
		}
	}
	fmt.Fprintln(w, builder.String())
	builder.Reset()

	// print out options
	for i, option := range options {
		v := reflect.ValueOf(option)

		// use the header map to generate tab-separated string of values
		builder.WriteString(fmt.Sprintf("%d\t", i+1))
		for _, fieldName := range headerMap {
			builder.WriteString(fmt.Sprintf("%s\t", v.FieldByName(fieldName)))
		}
		builder.WriteString("\n")
		fmt.Fprintf(w, builder.String())
		builder.Reset()
	}
	w.Flush()

	// prompt for input
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\nSelect an option [1 - %d, or 'q' to quit]: ", len(options))
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
	if index <= 0 || index > len(options) {
		return -1, fmt.Errorf("you must choose a number between 1 and %d (inclusive)", len(options))
	}

	return index - 1, nil
}
