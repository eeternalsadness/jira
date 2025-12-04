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
package issue

import (
	"github.com/eeternalsadness/jira/internal/util"
	"github.com/eeternalsadness/jira/pkg/jira"
	"github.com/spf13/cobra"
)

var jiraClient jira.Jira

// NewCommand creates and returns the issue command
func NewCommand() *cobra.Command {
	issueCmd := &cobra.Command{
		Use:   "issue",
		Short: "Manage Jira issues",
		Long:  `Create, get, transition, and manage Jira issues.`,
		Example: `# Get all your assigned issues
jira issue get --all

# Get a specific issue by ID
jira issue get PROJ-123

# Create a new issue
jira issue create

# Transition an issue
jira issue transition PROJ-123`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			if cmd.HasParent() {
				if err = cmd.Parent().PersistentPreRunE(cmd.Parent(), args); err != nil {
					return err
				}
			}

			jiraClient, err = util.InitJiraConfig()
			if err != nil {
				return err
			}

			return nil
		},
	}

	// Add subcommands
	issueCmd.AddCommand(newGetCommand())
	issueCmd.AddCommand(newCreateCommand())
	issueCmd.AddCommand(newTransitionCommand())

	return issueCmd
}
