package configure

import (
	"fmt"

	"github.com/eeternalsadness/jira/internal/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure Jira domain and credentials",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		configure()
	},
}

func configure() error {
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
