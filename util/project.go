package util

import (
	"encoding/json"
	"fmt"
)

type Project struct {
	ID   string
	Key  string
	Name string
	URL  string
}

func (jira *Jira) GetProjectByID(projectID int) (Project, error) {
	path := fmt.Sprintf("rest/api/3/project/%d", projectID)
	resp, err := jira.callApi(path, "GET", nil)
	if err != nil {
		return Project{}, fmt.Errorf("failed to call Jira API: %w", err)
	}

	// parse json data
	var data map[string]any
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return Project{}, fmt.Errorf("failed to unmarshal JSON response from Jira API: %w", err)
	}

	// get the necessary fields for the struct
	id := data["id"].(string)
	key := data["key"].(string)
	name := data["name"].(string)
	url := data["self"].(string)

	// form return struct
	outProject := Project{
		ID:   id,
		Key:  key,
		Name: name,
		URL:  url,
	}

	return outProject, nil
}
