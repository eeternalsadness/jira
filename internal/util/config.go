package util

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/eeternalsadness/jira/pkg/jira"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func InitConfig(cmd *cobra.Command, cfgFile string) error {
	if cfgFile != "" {
		// use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(fmt.Sprintf("%s/.config/jira/", home))
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")

		// set config key to default config file
		// viper.Set("config", fmt.Sprintf("%s/.config/jira/config.yaml", home))
		cfgFile = fmt.Sprintf("%s/.config/jira/config.yaml", home)
	}

	// make sure config dir is created
	cfgDir := path.Dir(cfgFile)
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		return err
	}

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundErr viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundErr) {
			// if running configure, set the config file
			if cmd.Name() == "configure" {
				viper.SetConfigFile(cfgFile)
				return nil
			}
			fmt.Println("Config file not found! Please run 'jira configure' to configure your Jira credentials.")
		}
		return err
	}

	return viper.BindPFlags(cmd.Flags())
}

func InitJiraConfig() (jira.Jira, error) {
	// get jira config
	var jiraClient jira.Jira
	if err := viper.Unmarshal(&jiraClient); err != nil {
		return jira.Jira{}, fmt.Errorf("failed to read the config file '%s': %s", viper.ConfigFileUsed(), err)
	}

	return jiraClient, nil
}

func ConfigJiraCredentials(domain *string, email *string, token *string) error {
	if domain == nil {
		panic("domain is nil")
	}
	if email == nil {
		panic("email is nil")
	}
	if token == nil {
		panic("token is nil")
	}

	viper.Set("domain", *domain)
	viper.Set("email", *email)
	viper.Set("token", *token)

	if err := viper.WriteConfig(); err != nil {
		return err
	}

	return nil
}
