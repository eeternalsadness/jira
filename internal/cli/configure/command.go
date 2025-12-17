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
package configure

import (
	"fmt"

	"github.com/eeternalsadness/jira/internal/util"
	"github.com/spf13/cobra"
)

var configurationOptions = []string{
	"Credentials",
	"Default issue type",
	"Default project",
}

// NewCommand creates and returns the issue command
func NewCommand() *cobra.Command {
	configureCmd := &cobra.Command{
		Use:   "configure",
		Short: "Configure credentials, issue types, or projects for the CLI tool",
		RunE: func(cmd *cobra.Command, args []string) error {
			configurationOption, err := selectConfigOption()
			if err != nil || configurationOption == "" {
				return err
			}

			switch configurationOption {
			case "Credentials":
				return configureCredentials()
			case "Default issue type":
				return configureDefaultIssueType()
			case "Default project":
				return configureDefaultProject()
			default:
				return fmt.Errorf("invalid configuration option: %s", configurationOption)
			}
		},
	}

	return configureCmd
}

func selectConfigOption() (string, error) {
	fmt.Println("Configuration options:")
	err := util.PrettyPrintStringSlice(configurationOptions)
	if err != nil {
		return "", err
	}

	index, err := util.UserSelectFromRange(len(configurationOptions))
	if err != nil {
		return "", err
	}

	return configurationOptions[index], nil
}
