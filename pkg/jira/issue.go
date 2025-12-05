package jira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// NOTE: follow Jira API reference
type Issue struct {
	ID             string
	Key            string
	Title          string
	Description    string
	Status         string
	StatusCategory string
	ProjectID      string
	IssueTypeID    string
	URL            string
}

func (jira *Jira) GetAssignedIssues() ([]Issue, error) {
	// call api
	jql := url.QueryEscape("assignee = currentuser() AND statuscategory != \"Done\"")
	fields := url.QueryEscape("summary,status,project,issuetype")
	path := fmt.Sprintf("rest/api/3/search/jql?jql=%s&fields=%s", jql, fields)
	resp, err := jira.callAPI(path, "GET", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call Jira API: %w", err)
	}

	// parse json data
	var data map[string]any
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response from Jira API: %w", err)
	}

	// transform json into output
	issues := data["issues"].([]any)
	outIssues := make([]Issue, len(issues))
	for i, issue := range issues {
		issueMap := issue.(map[string]any)
		fieldsMap := issueMap["fields"].(map[string]any)
		statusMap := fieldsMap["status"].(map[string]any)
		statusCategoryMap := statusMap["statusCategory"].(map[string]any)
		projectMap := fieldsMap["project"].(map[string]any)
		issueTypeMap := fieldsMap["issuetype"].(map[string]any)

		// get the necessary fields for the struct
		id := issueMap["id"].(string)
		key := issueMap["key"].(string)
		title := fieldsMap["summary"].(string)
		status := statusMap["name"].(string)
		statusCategory := statusCategoryMap["name"].(string)
		projectID := projectMap["id"].(string)
		issueTypeID := issueTypeMap["id"].(string)
		url := issueMap["self"].(string)
		outIssues[i] = Issue{
			ID:             id,
			Key:            key,
			Title:          title,
			Description:    "",
			Status:         status,
			StatusCategory: statusCategory,
			ProjectID:      projectID,
			IssueTypeID:    issueTypeID,
			URL:            url,
		}
	}

	return outIssues, nil
}

func (jira *Jira) GetIssueByID(issueID string) (Issue, error) {
	fields := url.QueryEscape("summary,description,comment,status")
	path := fmt.Sprintf("rest/api/3/issue/%s?fields=%s", issueID, fields)
	resp, err := jira.callAPI(path, "GET", nil)
	if err != nil {
		return Issue{}, fmt.Errorf("failed to call Jira API: %w", err)
	}

	// parse json data
	var data map[string]any
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return Issue{}, fmt.Errorf("failed to unmarshal JSON response from Jira API: %w", err)
	}

	// var out bytes.Buffer
	// json.Indent(&out, resp, "", "  ")
	// fmt.Println(out.String())

	// transform json into output
	fieldsMap := data["fields"].(map[string]any)
	statusMap := fieldsMap["status"].(map[string]any)
	statusCategoryMap := statusMap["statusCategory"].(map[string]any)

	// get the necessary fields for the struct
	id := data["id"].(string)
	key := data["key"].(string)
	title := fieldsMap["summary"].(string)
	descriptionMap, ok := fieldsMap["description"].(map[string]any)
	description := ""
	if ok {
		description = getIssueDescriptionText(descriptionMap)
	}
	status := statusMap["name"].(string)
	statusCategory := statusCategoryMap["name"].(string)
	url := data["self"].(string)

	// form return struct
	outIssue := Issue{
		ID:             id,
		Key:            key,
		Title:          title,
		Description:    description,
		Status:         status,
		StatusCategory: statusCategory,
		URL:            url,
	}

	return outIssue, nil
}

func getIssueDescriptionText(descriptionMap map[string]any) string {
	descriptionContent := descriptionMap["content"].([]any)
	descriptionSlice := make([]string, len(descriptionContent))
	for _, content := range descriptionContent {
		contentMap := content.(map[string]any)
		if contentMap["type"].(string) == "paragraph" {
			for _, contentContent := range contentMap["content"].([]any) {
				contentContentMap := contentContent.(map[string]any)
				if contentContentMap["type"].(string) == "text" {
					// NOTE: assume that there's only 1 text field per content object
					descriptionSlice = append(descriptionSlice, contentContentMap["text"].(string))
				}
			}
		}
	}

	return strings.Join(descriptionSlice, "\n")
}

func (jira *Jira) CreateIssue(projectID string, issueTypeID string, title string, description string) (string, error) {
	// get current user id
	currentUserID, err := jira.getCurrentUserID()
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
  }`, currentUserID, projectID, issueTypeID, descriptionField, title)

	// call api
	path := "rest/api/3/issue"
	resp, err := jira.callAPI(path, "POST", bytes.NewBuffer([]byte(body)))
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
