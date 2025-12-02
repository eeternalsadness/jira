package configure

import (
	"fmt"

	"github.com/eeternalsadness/jira/internal/util"
	"github.com/spf13/viper"
)

func configureCredentials() error {
	// configure jira domain
	err := util.ViperUpsertString("Domain", "Enter the Jira domain", "example.atlassian.net")
	if err != nil {
		return fmt.Errorf("failed to configure Jira domain: %s", err)
	}

	// configure jira email
	err = util.ViperUpsertString("Email", "Enter the email address used for Jira", "example@example.com")
	if err != nil {
		return fmt.Errorf("failed to configure Jira email: %s", err)
	}

	// configure jira api token

	if err = util.ViperUpsertString("Token", "Enter the Jira API token", ""); err != nil {
		return fmt.Errorf("failed to configure Jira API token: %s", err)
	}

	if err := viper.WriteConfig(); err != nil {
		return err
	}

	return nil
}
