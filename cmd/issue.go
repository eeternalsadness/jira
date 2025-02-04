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
	"text/tabwriter"

	"github.com/eeternalsadness/jira/util"
	"github.com/spf13/cobra"
)

// getIssueCmd represents the issue command when called by the get command
var getIssueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Get your current Jira issues",
	Long: `Get Jira issues that are assigned to the current user (you).
Issues with status 'Done', 'Rejected', or 'Cancelled' are not returned.`,
	Run: func(cmd *cobra.Command, args []string) {
		issues, err := jira.GetAssignedIssues()
		if err != nil {
			fmt.Printf("Failed to get assigned issues: %s\n", err)
			return
		}

		// print out issues
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tIssue\tStatus\t")
		for _, issue := range issues {
			fmt.Fprintf(w, "%s\t[%s] %s\t%s\t\n", issue.Id, issue.Key, issue.Title, issue.Status)
		}
		w.Flush()
	},
}

// createIssueCmd represents the issue command when called by the create command
var createIssueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Create a Jira issue",
	Long:  `Create a Jira issue in the 'KV FnB Web' project (default). The issue is assigned to the current user by default.`,
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

// transitionIssueCmd represents the issue command when called by the transition command
var transitionIssueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Transition a Jira issue",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Must have exactly 1 argument!")
			return
		}

		issueId := args[0]
		transitions, err := jira.GetTransitions(issueId)
		if err != nil {
			fmt.Printf("Failed to get valid transitions for issue: %s\n", err)
			return
		}

		transition, err := selectTransition(transitions)
		if err != nil {
			fmt.Printf("Failed when selecting a transition: %s\n", err)
			return
		}

		err = jira.TransitionIssue(issueId, transition.Id)
		if err != nil {
			fmt.Printf("Failed when transitioning issue %s: %s\n", issueId, err)
			return
		}

		fmt.Printf("Issue %s transitioned to '%s'.\n", issueId, transition.Name)
	},
}

func selectTransition(transitions []util.Transition) (util.Transition, error) {
	// print out available transitions
	fmt.Println("Available transitions:")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "#\tName\tCategory\t")
	for i, transition := range transitions {
		fmt.Fprintf(w, "%d\t%s\t%s\t\n", i+1, transition.Name, transition.Category)
	}
	w.Flush()

	// prompt for transition
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\nEnter the # to transition to [1 - %d, or 'q' to quit]: ", len(transitions))
	inputStr, err := reader.ReadString('\n')
	if err != nil {
		return util.Transition{}, err
	}
	inputStr = inputStr[:len(inputStr)-1]

	// quit if user types 'q'
	if inputStr == "q" {
		os.Exit(0)
	}

	// check number value
	transitionIndex, err := strconv.ParseInt(inputStr, 10, 64)
	if err != nil {
		return util.Transition{}, err
	}

	// check if index value is valid
	if transitionIndex <= 0 || transitionIndex > int64(len(transitions)) {
		return util.Transition{}, fmt.Errorf("transition # be a number between 1 and %d (inclusive)", len(transitions))
	}

	return transitions[transitionIndex-1], nil
}

func init() {
	getCmd.AddCommand(getIssueCmd)
	createCmd.AddCommand(createIssueCmd)
	transitionCmd.AddCommand(transitionIssueCmd)
}
