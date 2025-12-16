package configure

import (
	"fmt"

	"github.com/eeternalsadness/jira/internal/util"
	"github.com/spf13/viper"
)

func configureCredentials() error {
	// configure jira domain
	defaultDomain := "example.atlassian.net"
	domainKey := "domain"

	// if there is an existing value, use it as the default domain
	if viper.IsSet(domainKey) {
		defaultDomain = viper.GetString(domainKey)
	}

	domain, err := util.UserGetString(
		fmt.Sprintf("Enter the Jira domain [%s]", defaultDomain),
		&defaultDomain,
		false)
	if err != nil {
		return err
	}

	// configure jira email
	defaultEmail := "example@example.com"
	emailKey := "email"
	if viper.IsSet(emailKey) {
		defaultEmail = viper.GetString(emailKey)
	}
	email, err := util.UserGetString(
		fmt.Sprintf("Enter the email address used for Jira [%s]", defaultEmail),
		&defaultEmail,
		false)
	if err != nil {
		return err
	}

	// configure jira api token
	defaultToken := "example.atlassian.net"
	defaultTokenSensored := defaultToken
	tokenKey := "token"
	if viper.IsSet(tokenKey) {
		defaultToken = viper.GetString(tokenKey)
		defaultTokenSensored = util.SensorString(defaultToken)
	}
	token, err := util.UserGetString(
		fmt.Sprintf("Enter the Jira API token [%s]", defaultTokenSensored),
		&defaultToken,
		false)
	if err != nil {
		return err
	}

	if err := util.ConfigJiraCredentials(domain, email, token); err != nil {
		return err
	}

	return nil
}
