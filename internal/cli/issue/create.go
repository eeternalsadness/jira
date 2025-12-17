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
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/eeternalsadness/jira/internal/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	projectID   string
	issueTypeID string
)

func newCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [-p PROJECT_ID] [-t ISSUE_TYPE_ID]",
		Short: "Create a Jira issue",
		Long:  `Create a Jira issue in the specified project. The issue is assigned to the current user by default.`,
		Args:  cobra.MaximumNArgs(2),
		Example: `# Create a Jira issue with the default project and issue type
jira issue create

# Create a Jira issue with a specific project and issue type
jira issue create --project-id 123 --issue-type-id 456`,
		PreRun: func(cmd *cobra.Command, args []string) {
			if projectID == "" {
				projectID = viper.GetString(string(util.DefaultProjectIDKey))
			}

			if issueTypeID == "" {
				issueTypeID = viper.GetString(string(util.DefaultIssueTypeIDKey))
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			return createIssue()
		},
	}

	cmd.Flags().StringVarP(&projectID, "project-id", "p", "", "create an issue in the specified project")
	cmd.Flags().StringVarP(&issueTypeID, "issue-type-id", "t", "", "specify the issue type to create")

	return cmd
}

func createIssue() error {
	reader := bufio.NewReader(os.Stdin)

	// prompt for issue's title
	fmt.Print("Enter the issue's title: ")
	title, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read user input: %s", err)
	}
	title = title[:len(title)-1]

	// title can't be empty
	if len(strings.TrimSpace(title)) == 0 {
		return fmt.Errorf("issue's title can't be empty")
	}

	// prompt for issue's description
	fmt.Print("Enter the issue's description (optional): ")
	description, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read user input: %s", err)
	}
	description = description[:len(description)-1]

	// create issue
	issueKey, err := jiraClient.CreateIssue(projectID, issueTypeID, title, description)
	if err != nil {
		return fmt.Errorf("failed to create Jira issue: %s", err)
	}

	fmt.Printf("Issue '%s' created.\nURL: https://%s/browse/%s.\n", title, jiraClient.Domain, issueKey)
	return nil
}
