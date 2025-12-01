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
package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

// TODO: determine what goes into the config file and what doesn't
const (
	checkVersionInterval = 10 * time.Minute
	tmpDir               = "/tmp/jira-cobra"
)

// TODO: maybe refactor this to another package
func checkVersion(cmd *cobra.Command) {
	// only check every once in a while
	lastCheckVersionTime, err := getLastCheckVersionTime()
	if err != nil {
		return
	}

	if time.Since(lastCheckVersionTime) < checkVersionInterval {
		return
	}

	latestVersion, err := getLatestVersion()
	cobra.CheckErr(updateLastCheckVersionTime())

	// ignore errors
	if err == nil && latestVersion != cmd.Root().Version {
		// TODO: potentially add an update command instead of telling the user to update manually
		fmt.Printf("\n\033[33mVersion '%s' is available. To update to the latest version, run:\n  go install github.com/eeternalsadness/jira@latest\033[0m\n", latestVersion)
	}
}

func getLastCheckVersionTime() (time.Time, error) {
	tmpFilePath := fmt.Sprintf("%s/lastCheckVersionTime", tmpDir)
	data, err := os.ReadFile(tmpFilePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return time.Time{}, fmt.Errorf("failed to read check version time file at '%s': %w", tmpFilePath, err)
		} else {
			return time.Unix(0, 0), nil
		}
	}

	checkVersionTime, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse time from file: %w", err)
	}

	return time.Unix(checkVersionTime, 0), nil
}

func updateLastCheckVersionTime() error {
	// make sure tmp dir is set up
	err := os.MkdirAll(tmpDir, 0o755)
	cobra.CheckErr(err)

	// write to file
	// TODO: maybe put this in the config file
	tmpFilePath := fmt.Sprintf("%s/lastCheckVersionTime", tmpDir)
	err = os.WriteFile(tmpFilePath, []byte(strconv.FormatInt(time.Now().Unix(), 10)), 0o755)
	return err
}

// TODO: put this in a different package
func getLatestVersion() (string, error) {
	// FIXME: move hard coded value elsewhere
	githubEndpoint := "https://api.github.com/repos/eeternalsadness/jira/releases/latest"

	// call github releases api endpoint
	resp, err := http.Get(githubEndpoint)
	if err != nil {
		return "", fmt.Errorf("failed to reach the Github endpoint to check version: %w", err)
	}
	defer resp.Body.Close()

	// non-200 status code
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("%s", resp.Status)
	}

	// read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body from Github: %w", err)
	}

	// parse as json
	var data map[string]interface{}
	err = json.Unmarshal(respBody, &data)
	if err != nil {
		return "", fmt.Errorf("failed to parse response body as JSON: %w", err)
	}

	// extract tag
	tag, ok := data["tag_name"].(string)
	if !ok {
		return "", fmt.Errorf("expected 'tag_name' to be a string, got %T", data["tag_name"])
	}

	return tag, nil
}
