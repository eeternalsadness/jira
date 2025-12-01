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
	"errors"
	"fmt"
	"os"

	"github.com/eeternalsadness/jira/internal/cli/issue"
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
	Version: "v0.1.7",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initConfig(cmd)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/jira/config.yaml)")

	rootCmd.AddCommand(issue.NewCommand())
}

func initConfig(cmd *cobra.Command) error {
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

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundErr viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundErr) {
			fmt.Println("Config file not found! Please run 'jira configure' to configure your Jira credentials.")
		}
		return err
	}

	if cfgFile != "" {
		fmt.Fprintln(os.Stdout, "Using config file: ", viper.ConfigFileUsed())
	}
	// get jira config
	err := viper.Unmarshal(&jira)
	if err != nil {
		return fmt.Errorf("failed to read the config file '%s': %s", cfgFile, err)
	}

	err = viper.BindPFlags(cmd.Flags())
	if err != nil {
		return err
	}

	return nil
}
