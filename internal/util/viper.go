package util

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

func ViperUpsertString(key string, prompt string, example string) error {
	reader := bufio.NewReader(os.Stdin)

	if len(example) > 0 {
		fmt.Printf("%s [%s]: ", prompt, example)
	} else {
		fmt.Printf("%s: ", prompt)
	}

	userInput, _ := reader.ReadString('\n')
	userInput = userInput[:len(userInput)-1]

	// check for existing config
	if viper.GetString(key) != "" {
		overwrite, err := UserYesNo(fmt.Sprintf("Configuration for key '%s' already exists. Overwrite?", key))
		if err != nil {
			return err
		}

		if overwrite {
			viper.Set(key, userInput)
		}
	} else {
		viper.Set(key, userInput)
	}

	return nil
}

func ViperUpsertInt(key string, prompt string, example string) error {
	reader := bufio.NewReader(os.Stdin)

	if len(example) > 0 {
		fmt.Printf("%s [%s]: ", prompt, example)
	} else {
		fmt.Printf("%s: ", prompt)
	}

	userInput, _ := reader.ReadString('\n')
	userInput = userInput[:len(userInput)-1]

	// check if input string is valid int
	userInputInt, err := strconv.Atoi(strings.TrimSpace(userInput))
	if err != nil {
		return err
	}

	// check for existing config
	if viper.GetString(key) != "" {
		overwrite, err := UserYesNo(fmt.Sprintf("Configuration for key '%s' already exists. Overwrite?", key))
		if err != nil {
			return err
		}

		if overwrite {
			viper.Set(key, userInputInt)
		}
	} else {
		viper.Set(key, userInputInt)
	}

	return nil
}
