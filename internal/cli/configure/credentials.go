package configure

import (
	"fmt"

	"github.com/eeternalsadness/jira/internal/util"
	"github.com/spf13/viper"
)

func configureCredentials() error {
	// configure jira domain
	example := "example.atlassian.net"
	err := util.ViperUpsertString(
		"domain",
		"Enter the Jira domain",
		&example,
		false)
	if err != nil {
		return fmt.Errorf("failed to configure Jira domain: %s", err)
	}

	// configure jira email
	example = "example@example.com"
	err = util.ViperUpsertString(
		"email",
		"Enter the email address used for Jira",
		&example,
		false)
	if err != nil {
		return fmt.Errorf("failed to configure Jira email: %s", err)
	}

	// configure jira api token
	err = util.ViperUpsertString(
		"token",
		"Enter the Jira API token",
		nil,
		true)
	if err != nil {
		return fmt.Errorf("failed to configure Jira API token: %s", err)
	}

	if err := viper.WriteConfig(); err != nil {
		return err
	}

	return nil
}
