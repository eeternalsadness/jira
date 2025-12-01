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
	"strconv"
	"text/tabwriter"

	"github.com/eeternalsadness/jira/pkg/jira"
	"github.com/spf13/cobra"
)

func newTransitionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transition ISSUE_ID",
		Short: "Transition a Jira issue",
		Long:  `Transition a Jira issue using the issue ID. The user is prompted to select a transition from a list of valid transitions for the issue.`,
		Args:  cobra.ExactArgs(1),
		Example: `# Create a Jira issue with the default project and issue type
jira issue create

# Create a Jira issue with a specific project and issue type
jira issue create --project-id 123 --issue-type-id 456`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			issueID, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			return transitionIssue(issueID)
		},
	}

	return cmd
}

func transitionIssue(issueID int) error {
	transitions, err := jiraClient.GetTransitions(issueID)
	if err != nil {
		return fmt.Errorf("failed to get valid transitions for issue: %s", err)
	}

	transition, err := selectTransition(transitions)
	if err != nil {
		return fmt.Errorf("failed when selecting a transition: %s", err)
	}

	// user quits
	if transition == (jira.Transition{}) {
		return nil
	}

	err = jiraClient.TransitionIssue(issueID, transition.ID)
	if err != nil {
		return fmt.Errorf("failed when transitioning issue %d: %s", issueID, err)
	}

	fmt.Printf("Issue %d transitioned to '%s'.\n", issueID, transition.Name)
	return nil
}

func selectTransition(transitions []jira.Transition) (jira.Transition, error) {
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
		return jira.Transition{}, err
	}
	inputStr = inputStr[:len(inputStr)-1]

	// quit if user types 'q'
	if inputStr == "q" {
		return jira.Transition{}, nil
	}

	// check number value
	transitionIndex, err := strconv.ParseInt(inputStr, 10, 64)
	if err != nil {
		return jira.Transition{}, err
	}

	// check if index value is valid
	if transitionIndex <= 0 || transitionIndex > int64(len(transitions)) {
		return jira.Transition{}, fmt.Errorf("transition # be a number between 1 and %d (inclusive)", len(transitions))
	}

	return transitions[transitionIndex-1], nil
}
