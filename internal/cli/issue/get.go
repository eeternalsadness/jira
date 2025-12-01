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
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var isAll bool

func newGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get your current Jira issues",
		Long: `Get Jira issues that are assigned to the current user (you).
	Issues with status category 'Done' are not returned.`,
		Args: cobra.MaximumNArgs(1),
		Example: `# Get all your assigned issues
jira issue get --all

# Get a specific issue by ID
jira issue get PROJ-123`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return getIssue(cmd, args)
		},
	}

	cmd.SilenceUsage = true
	cmd.Flags().BoolVarP(&isAll, "all", "a", false, "get all issues assigned to you")

	return cmd
}

func getIssue(cmd *cobra.Command, args []string) error {
	if isAll && len(args) > 0 {
		return fmt.Errorf("cannot use --all with an issue ID")
	} else if !isAll && len(args) == 0 {
		cmd.Usage()
		return fmt.Errorf("missing argument or flags")
	} else if isAll {
		issues, err := jiraClient.GetAssignedIssues()
		if err != nil {
			return fmt.Errorf("failed to get assigned issues: %s", err)
		}

		// print out issues
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tIssue\tStatus\tStatus Category\t")
		for _, issue := range issues {
			fmt.Fprintf(w, "%s\t[%s] %s\t%s\t%s\t\n", issue.Id, issue.Key, issue.Title, issue.Status, issue.StatusCategory)
		}
		w.Flush()
	} else {
		issueId := args[0]
		issue, err := jiraClient.GetIssueById(issueId)
		if err != nil {
			return fmt.Errorf("failed to get assigned issue: %s", err)
		}

		// print out issue
		fmt.Printf("[%s] %s\n\n", issue.Key, issue.Title)
		fmt.Printf("Description:\n%s\n", issue.Description)
	}

	return nil
}
