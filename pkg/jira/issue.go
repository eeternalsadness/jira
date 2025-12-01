package jira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type Issue struct {
	Id             string
	Key            string
	Title          string
	Description    string
	Status         string
	StatusCategory string
	Url            string
}

func (jira *Jira) GetAssignedIssues() ([]Issue, error) {
	// call api
	jql := url.QueryEscape("assignee = currentuser() AND statuscategory != \"Done\"")
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
		status := statusMap["name"].(string)
		statusCategory := statusCategoryMap["name"].(string)
		url := issueMap["self"].(string)
		outIssues[i] = Issue{
			Id:             id,
			Key:            key,
			Title:          title,
			Description:    "",
			Status:         status,
			StatusCategory: statusCategory,
			Url:            url,
		}
	}

	return outIssues, nil
}

func (jira *Jira) GetIssueById(issueId string) (Issue, error) {
	fields := url.QueryEscape("summary,description,comment,status")
	path := fmt.Sprintf("rest/api/3/issue/%s?fields=%s", issueId, fields)
	resp, err := jira.callApi(path, "GET", nil)
	if err != nil {
		return Issue{}, fmt.Errorf("failed to call Jira API: %w", err)
	}

	// parse json data
	var data map[string]interface{}
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return Issue{}, fmt.Errorf("failed to unmarshal JSON response from Jira API: %w", err)
	}

	// transform json into output
	fieldsMap := data["fields"].(map[string]interface{})
	statusMap := fieldsMap["status"].(map[string]interface{})
	statusCategoryMap := statusMap["statusCategory"].(map[string]interface{})

	// get the necessary fields for the struct
	id := data["id"].(string)
	key := data["key"].(string)
	title := fieldsMap["summary"].(string)
	description := getIssueDescriptionText(fieldsMap["description"].(map[string]interface{}))
	status := statusMap["name"].(string)
	statusCategory := statusCategoryMap["name"].(string)
	url := data["self"].(string)

	// form return struct
	outIssue := Issue{
		Id:             id,
		Key:            key,
		Title:          title,
		Description:    description,
		Status:         status,
		StatusCategory: statusCategory,
		Url:            url,
	}

	return outIssue, nil
}

func getIssueDescriptionText(descriptionMap map[string]interface{}) string {
	descriptionContent := descriptionMap["content"].([]interface{})
	descriptionSlice := make([]string, len(descriptionContent))
	for _, content := range descriptionContent {
		contentMap := content.(map[string]interface{})
		if contentMap["type"].(string) == "paragraph" {
			for _, contentContent := range contentMap["content"].([]interface{}) {
				contentContentMap := contentContent.(map[string]interface{})
				if contentContentMap["type"].(string) == "text" {
					// NOTE: assume that there's only 1 text field per content object
					descriptionSlice = append(descriptionSlice, contentContentMap["text"].(string))
				}
			}
		}
	}

	return strings.Join(descriptionSlice, "\n")
}

func (jira *Jira) CreateIssue(projectId int, issueTypeId int, title string, description string) (string, error) {
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
        "id": "%d"
      },
      "issuetype": {
        "id": "%d"
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
	var data map[string]any
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON response from Jira API: %w", err)
	}

	// transform json to output
	issueKey := data["key"].(string)

	return issueKey, nil
}
