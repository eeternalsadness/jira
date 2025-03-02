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
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// getIssueCmd represents the issue command when called by the get command
var getIssueCmd = &cobra.Command{
	Use:   "issue (issueId | --all)",
	Short: "Get your current Jira issues",
	Long: `Get Jira issues that are assigned to the current user (you).
Issues with status 'Done', 'Rejected', or 'Cancelled' are not returned.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		getAllIssues, err := cmd.Flags().GetBool("all")
		if err != nil {
			fmt.Printf("Failed to read --all flag: %s\n", err)
			return
		}

		if getAllIssues && len(args) > 0 {
			fmt.Println("Cannot use --all with an issue ID!")
			return
		} else if !getAllIssues && len(args) == 0 {
			fmt.Println("Missing argument or flags!")
			cmd.Help()
			return
		} else if getAllIssues {
			issues, err := jira.GetAssignedIssues()
			if err != nil {
				fmt.Printf("Failed to get assigned issues: %s\n", err)
				return
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
			issue, err := jira.GetIssueById(issueId)
			if err != nil {
				fmt.Printf("Failed to get assigned issue: %s\n", err)
				return
			}

			// print out issue
			fmt.Printf("[%s] %s\n\n", issue.Key, issue.Title)
			fmt.Printf("Description:\n%s\n", issue.Description)
		}
	},
}

// createIssueCmd represents the issue command when called by the create command
var createIssueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Create a Jira issue",
	Long:  `Create a Jira issue in the 'KV FnB Web' project (default). The issue is assigned to the current user by default.`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		// prompt for issue's title
		fmt.Print("Enter the issue's title: ")
		title, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Failed to read user input: %s\n", err)
		}
		title = title[:len(title)-1]

		// title can't be empty
		if len(strings.TrimSpace(title)) == 0 {
			fmt.Println("Issue's title can't be empty!")
			return
		}

		// prompt for issue's description
		fmt.Print("Enter the issue's description (optional): ")
		description, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Failed to read user input: %s\n", err)
			return
		}
		description = description[:len(description)-1]

		// create issue
		// TODO: find a way to create issues for different projects
		// TODO: find a way to create different issue types
		issueKey, err := jira.CreateIssue(jira.DefaultProjectId, jira.DefaultIssueTypeId, title, description)
		if err != nil {
			fmt.Printf("Failed to create Jira issue: %s\n", err)
			return
		}

		fmt.Printf("Issue '%s' created.\nURL: https://%s/browse/%s.\n", title, jira.Domain, issueKey)
	},
}

func init() {
	getCmd.AddCommand(getIssueCmd)
	getIssueCmd.Flags().BoolP("all", "a", false, "get all issues assigned to you")
	createCmd.AddCommand(createIssueCmd)
}
