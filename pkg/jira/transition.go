package jira

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Transition struct {
	ID       string
	Name     string
	Category string
}

func (jira *Jira) GetTransitions(issueID string) ([]Transition, error) {
	// call api
	path := fmt.Sprintf("rest/api/3/issue/%s/transitions", issueID)
	resp, err := jira.callApi(path, "GET", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call Jira API: %w", err)
	}

	// parse json
	var data map[string]any
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response from Jira API: %w", err)
	}

	// transform json into output
	transitions := data["transitions"].([]any)
	outTransitions := make([]Transition, len(transitions))
	for i, transition := range transitions {
		transitionMap := transition.(map[string]any)
		toMap := transitionMap["to"].(map[string]any)
		statusCategory := toMap["statusCategory"].(map[string]any)

		// get the necessary fields for the struct
		id := transitionMap["id"].(string)
		name := transitionMap["name"].(string)
		categoryName := statusCategory["name"].(string)
		outTransitions[i] = Transition{
			ID:       id,
			Name:     name,
			Category: categoryName,
		}
	}

	return outTransitions, nil
}

func (jira *Jira) TransitionIssue(issueID string, transitionID string) error {
	// call api
	body := fmt.Sprintf(`{
    "transition": {
      "id": "%s"
    }
  }`, transitionID)
	path := fmt.Sprintf("rest/api/3/issue/%s/transitions", issueID)
	resp, err := jira.callApi(path, "POST", bytes.NewBuffer([]byte(body)))
	if err != nil {
		return fmt.Errorf("failed to call Jira API: %w", err)
	}

	fmt.Println(string(resp))

	return nil
}
