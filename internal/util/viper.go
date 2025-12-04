package util

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

func ViperUpsertString(key string, prompt string, defaultVal *string, sensitive bool) error {
	reader := bufio.NewReader(os.Stdin)

	currentVal := viper.GetString(key)

	// if there is an existing value, use it instead of the default val
	if viper.IsSet(key) {
		defaultVal = &currentVal
	}

	if defaultVal != nil {
		if sensitive {
			fmt.Printf("%s [%s]: ", prompt, sensorString(*defaultVal))
		} else {
			fmt.Printf("%s [%s]: ", prompt, *defaultVal)
		}
	} else {
		fmt.Printf("%s: ", prompt)
	}

	userInput, _ := reader.ReadString('\n')
	userInput = userInput[:len(userInput)-1]

	if len(userInput) > 0 {
		viper.Set(key, userInput)
	} else if defaultVal != nil {
		viper.Set(key, *defaultVal)
	}

	return nil
}

func ViperUpsertInt(key string, prompt string, defaultVal *int) error {
	reader := bufio.NewReader(os.Stdin)

	// if there is an existing value, use it instead of the default val
	if viper.IsSet(key) {
		*defaultVal = viper.GetInt(key)
	}

	if defaultVal != nil {
		fmt.Printf("%s [%d]: ", prompt, *defaultVal)
	} else {
		fmt.Printf("%s: ", prompt)
	}

	userInput, _ := reader.ReadString('\n')
	userInput = userInput[:len(userInput)-1]

	if len(userInput) > 0 {
		// check if input string is valid int
		userInputInt, err := strconv.Atoi(strings.TrimSpace(userInput))
		if err != nil {
			return err
		}
		viper.Set(key, userInputInt)
	} else {
		if defaultVal == nil {
			return fmt.Errorf("invalid number")
		}
		viper.Set(key, *defaultVal)
	}

	return nil
}

func sensorString(str string) string {
	// otherwise show the first and last 25% chars (max 4)
	charsToShow := len(str) / 4
	charsToShow = min(charsToShow, 4)

	// limit the sensored part to 20 chars
	sensoredLen := min(len(str)-charsToShow*2, 20)

	return fmt.Sprintf("%s%s%s", str[:charsToShow], strings.Repeat("*", sensoredLen), str[len(str)-charsToShow:])
}
