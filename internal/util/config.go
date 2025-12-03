package util

import (
	"errors"
	"fmt"
	"os"

	"github.com/eeternalsadness/jira/pkg/jira"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func InitConfig(cmd *cobra.Command) (jira.Jira, error) {
	if cfgFile := viper.GetString("config"); cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(fmt.Sprintf("%s/.config/jira/", home))
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundErr viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundErr) {
			// ignore this error for the configure command
			if cmd.Name() == "configure" {
				return jira.Jira{}, nil
			}
			fmt.Println("Config file not found! Please run 'jira configure' to configure your Jira credentials.")
		}
		return jira.Jira{}, err
	}

	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return jira.Jira{}, err
	}

	// get jira config
	var jiraClient jira.Jira
	if err := viper.Unmarshal(&jiraClient); err != nil {
		return jira.Jira{}, fmt.Errorf("failed to read the config file '%s': %s", viper.ConfigFileUsed(), err)
	}

	return jiraClient, nil
}
