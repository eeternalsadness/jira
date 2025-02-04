package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
)

type Issue struct {
	Id     string
	Key    string
	Title  string
	Status string
	Url    string
}

func (jira *Jira) GetAssignedIssues() ([]Issue, error) {
	// call api
	jql := url.QueryEscape("assignee = currentuser() AND status NOT IN (Done, Rejected, Cancelled)")
	fields := url.QueryEscape("summary,status")
	path := fmt.Sprintf("rest/api/3/search/jql?jql=%s&fields=%s", jql, fields)
	resp, err := jira.callApi(path, "GET", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call Jira API: %w", err)
	}

	// parse json data
	var data map[string]interface{}
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response from Jira API: %w", err)
	}

	// transform json into output
	issues := data["issues"].([]interface{})
	outIssues := make([]Issue, len(issues))
	for i, issue := range issues {
		issueMap := issue.(map[string]interface{})
		fieldsMap := issueMap["fields"].(map[string]interface{})
		statusMap := fieldsMap["status"].(map[string]interface{})
		statusCategoryMap := statusMap["statusCategory"].(map[string]interface{})

		// get the necessary fields for the struct
		id := issueMap["id"].(string)
		key := issueMap["key"].(string)
		title := fieldsMap["summary"].(string)
		status := statusCategoryMap["name"].(string)
		url := issueMap["self"].(string)
		outIssues[i] = Issue{
			Id:     id,
			Key:    key,
			Title:  title,
			Status: status,
			Url:    url,
		}
	}

	return outIssues, nil
}

func (jira *Jira) CreateIssue(projectId string, issueTypeId string, title string, description string) (string, error) {
	// get current user id
	currentUserId, err := jira.getCurrentUserId()
	if err != nil {
		return "", fmt.Errorf("failed to get current user ID: %w", err)
	}

	// form description field if passed in
	descriptionField := ""
	if description != "" {
		descriptionField = fmt.Sprintf(`"description": {
      "content": [
        {
          "content": [
            {
              "text": "%s",
              "type": "text"
            }
          ],
          "type": "paragraph"
        }
      ],
      "type": "doc",
      "version": 1
    },`, description)
	}

	// form request body
	body := fmt.Sprintf(`{
    "fields": {
      "assignee": {
        "id": "%s"
      },
      "project": {
        "id": "%s"
      },
      "issuetype": {
        "id": "%s"
      },
      %s
      "summary": "%s"
    },
    "update": {}
  }`, currentUserId, projectId, issueTypeId, descriptionField, title)

	// call api
	path := "rest/api/3/issue"
	resp, err := jira.callApi(path, "POST", bytes.NewBuffer([]byte(body)))
	if err != nil {
		return "", fmt.Errorf("failed to call Jira API: %w", err)
	}

	// parse json data
	var data map[string]interface{}
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON response from Jira API: %w", err)
	}

	// transform json to output
	issueKey := data["key"].(string)

	return issueKey, nil
}
