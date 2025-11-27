/*
Copyright © 2025 Bach Nguyen <69bnguyen@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure Jira domain and credentials",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		configure()
	},
}

// configureProjectsCmd represents the configure projects command
var configureProjectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Configure the list of available Jira projects",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		configureProjects()
	},
}

func configure() {
	reader := bufio.NewReader(os.Stdin)

	// create config folder
	os.MkdirAll(cfgPath, 0o755)

	// configure jira domain
	err := viperUpsertString("Domain", "Enter the Jira domain", "example.atlassian.net")
	if err != nil {
		fmt.Printf("Failed to configure Jira domain: %s\n", err)
		return
	}

	// configure jira email
	viperUpsertString("Email", "Enter the email address used for Jira", "example@example.com")
	if err != nil {
		fmt.Printf("Failed to configure Jira email: %s\n", err)
		return
	}

	// configure jira api token
	viperUpsertString("Token", "Enter the Jira API token", "")
	if err != nil {
		fmt.Printf("Failed to configure Jira API token: %s\n", err)
		return
	}

	viper.WriteConfigAs(fmt.Sprintf("%s/config.yaml", cfgPath))
}

func configureProjects() {
	reader := bufio.NewReader(os.Stdin)

	projectIds := viper.GetIntSlice("ProjectIds")
	fmt.Println("Current project IDs:")
	fmt.Println(projectIds)

	viper.Set("DefaultProjectId", defaultProjectId)
	viper.SetDefault("DefaultProjectId", "10140")

	// configure jira email
	// fmt.Print("Enter the default issue type ID for the project [12345]: ")
	// defaultIssueTypeId, _ := reader.ReadString('\n')
	// defaultIssueTypeId = defaultIssueTypeId[:len(defaultIssueTypeId)-1]
	// viper.Set("DefaultIssueTypeId", defaultIssueTypeId)
	// FIXME: set default for now, to think of a better way to use issue type id
	viper.SetDefault("DefaultIssueTypeId", "10101")

	viper.WriteConfigAs(fmt.Sprintf("%s/config.yaml", cfgPath))
}

func addProjects(currentProjectIds []int) ([]int, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter the list of project IDs to add (separate by commas): ")
	userInput, _ := reader.ReadString('\n')
	userInput = userInput[:len(userInput)-1]
	userInputSlice := strings.Split(userInput, ",")
	// parse project ids
	projectIdsToAdd := make([]int, 10)
	for i := 0; i < len(userInputSlice); i++ {
		projectIdString := strings.TrimSpace(userInputSlice[i])
		projectIdInt, err := strconv.Atoi(projectIdString)
		if err != nil {
			return nil, fmt.Errorf("Invalid project ID: %s", projectIdString)
		}
		projectIdsToAdd = append(projectIdsToAdd, projectIdInt)
	}

	return append(currentProjectIds, projectIdsToAdd...), nil
}

func userYesNo(prompt string) (bool, error) {
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
		return false, fmt.Errorf("Input must be 'y' or 'n'.")
	}
}

func viperUpsertString(key string, prompt string, example string) error {
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
		overwrite, err := userYesNo(fmt.Sprintf("Configuration for key '%s' already exists. Overwrite?", key))
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

func init() {
	rootCmd.AddCommand(configureCmd)
	configureCmd.AddCommand(configureProjectsCmd)
}
