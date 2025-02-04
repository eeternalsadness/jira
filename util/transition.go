package util

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Transition struct {
	Id       string
	Name     string
	Category string
}

func (jira *Jira) GetTransitions(issueId string) ([]Transition, error) {
	// call api
	path := fmt.Sprintf("rest/api/3/issue/%s/transitions", issueId)
	resp, err := jira.callApi(path, "GET", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call Jira API: %w", err)
	}

	// parse json
	var data map[string]interface{}
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response from Jira API: %w", err)
	}

	// transform json into output
	transitions := data["transitions"].([]interface{})
	outTransitions := make([]Transition, len(transitions))
	for i, transition := range transitions {
		transitionMap := transition.(map[string]interface{})
		toMap := transitionMap["to"].(map[string]interface{})
		statusCategory := toMap["statusCategory"].(map[string]interface{})

		// get the necessary fields for the struct
		id := transitionMap["id"].(string)
		name := transitionMap["name"].(string)
		categoryName := statusCategory["name"].(string)
		outTransitions[i] = Transition{
			Id:       id,
			Name:     name,
			Category: categoryName,
		}
	}

	return outTransitions, nil
}

func (jira *Jira) TransitionIssue(issueId string, transitionId string) error {
	// call api
	body := fmt.Sprintf(`{
    "transition": {
      "id": "%s"
    }
  }`, transitionId)
	path := fmt.Sprintf("rest/api/3/issue/%s/transitions", issueId)
	resp, err := jira.callApi(path, "POST", bytes.NewBuffer([]byte(body)))
	if err != nil {
		return fmt.Errorf("failed to call Jira API: %w", err)
	}

	fmt.Println(string(resp))

	return nil
}
