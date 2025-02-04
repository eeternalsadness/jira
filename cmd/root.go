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
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/eeternalsadness/jira/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	cfgPath string
	jira    util.Jira
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "jira",
	Short:   "A CLI tool to do common Jira tasks",
	Long:    `This CLI tool aims to carry out common Jira tasks, helping you to stay in the command line instead of breaking your workflow and going to your web browser for Jira tasks.`,
	Version: "v0.1.3",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		latestVersion, err := getLatestVersion()
		fmt.Println(latestVersion)
		fmt.Println(err)
		// ignore errors
		if err == nil && latestVersion != cmd.Version {
			fmt.Printf("\033[33mVersion '%s' is available. To update to the latest version, run:\ngo install github.com/eeternalsadness/jira@latest\033[0m\n", latestVersion)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/jira/config.yaml)")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		cfgPath = fmt.Sprintf("%s/.config/jira/", home)
		viper.AddConfigPath(cfgPath)
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	// TODO: is it possible to check which command was called?
	if err := viper.ReadInConfig(); err == nil {
		if cfgFile != "" {
			fmt.Fprintln(os.Stdout, "Using config file: ", viper.ConfigFileUsed())
		}
		// get jira config
		err = viper.Unmarshal(&jira)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read the config file '%s': %w", cfgFile, err)
			os.Exit(1)
		}

	} else if _, errStat := os.Stat(cfgFile); cfgFile != "" && os.IsNotExist(errStat) {
		fmt.Fprintln(os.Stderr, "Config file not found: ", viper.ConfigFileUsed())
		os.Exit(1)
	}
}

func getLatestVersion() (string, error) {
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
