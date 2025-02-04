package util

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Jira struct {
	Domain             string
	Email              string
	Token              string
	DefaultProjectId   string
	DefaultIssueTypeId string
}

func (jira *Jira) getAuthToken() string {
	auth := fmt.Sprintf("%s:%s", jira.Email, jira.Token)
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func (jira *Jira) callApi(path string, method string, body io.Reader) ([]byte, error) {
	// form http request
	client := &http.Client{}
	req, err := http.NewRequest(method, fmt.Sprintf("https://%s/%s", jira.Domain, path), body)
	if err != nil {
		return nil, fmt.Errorf("failed to form a HTTP request: %w", err)
	}

	// set headers
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", jira.getAuthToken()))
	req.Header.Add("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// send http request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call the Jira API: %w", err)
	}
	defer resp.Body.Close()

	// read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the response from the Jira API: %w", err)
	}

	// non-200 status code
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		if len(respBody) > 0 {
			var data map[string]interface{}
			err = json.Unmarshal(respBody, &data)
			if err != nil {
				return nil, fmt.Errorf("failed to parse JSON response from Jira API: %w", err)
			}

			jsonOutput, err := json.MarshalIndent(data, "", "  ")
			if err != nil {
				return nil, fmt.Errorf("failed to format JSON response from Jira API: %w", err)
			}

			return nil, fmt.Errorf("%s:\n%s", resp.Status, jsonOutput)
		} else {
			return nil, fmt.Errorf("%s", resp.Status)
		}
	}

	return respBody, nil
}

func (jira *Jira) getCurrentUserId() (string, error) {
	// call api
	path := "rest/api/3/myself"
	resp, err := jira.callApi(path, "GET", nil)
	if err != nil {
		return "", fmt.Errorf("failed to call Jira API: %w", err)
	}

	// parse json
	var data map[string]interface{}
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON response from Jira API: %w", err)
	}

	// transform json to output
	accountId := data["accountId"].(string)

	return accountId, nil
}
