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

	"github.com/eeternalsadness/jira/internal/util"
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

			issueID := args[0]
			return transitionIssue(issueID)
		},
	}

	return cmd
}

func transitionIssue(issueID string) error {
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
		return fmt.Errorf("failed when transitioning issue %s: %s", issueID, err)
	}

	fmt.Printf("Issue %s transitioned to '%s'.\n", issueID, transition.Name)
	return nil
}

func selectTransition(transitions []jira.Transition) (jira.Transition, error) {
	// form header map
	headerMap := map[string]string{
		"Name":     "Name",
		"Category": "Category",
	}

	// prompt user for transition
	fmt.Println("Available transitions:")
	err := util.PrettyPrintStructSlice(headerMap, transitions)
	if err != nil {
		return jira.Transition{}, err
	}

	transitionIndex, err := util.UserSelectFromRange(len(transitions))
	if err != nil {
		if err == util.ErrUserQuit {
			return jira.Transition{}, nil
		} else {
			return jira.Transition{}, err
		}
	}

	return transitions[transitionIndex], nil
}
