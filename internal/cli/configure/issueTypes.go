package configure

import (
	"fmt"

	"github.com/eeternalsadness/jira/internal/util"
	"github.com/spf13/viper"
)

// TODO: make the user search for default issue type using name instead
func configureDefaultIssueType() error {
	// configure default issue type
	var defaultVal *string = nil
	defaultValStr := ""
	defaultIssueTypeIDKey := "default_issue_type_id"

	// if there is an existing value, use it as the default domain
	if viper.IsSet(defaultIssueTypeIDKey) {
		*defaultVal = viper.GetString(defaultIssueTypeIDKey)
		defaultValStr = fmt.Sprintf(" [%s]", *defaultVal)
	}

	defaultIssueTypeID, err := util.UserGetString(
		fmt.Sprintf("Enter the default issue type ID%s: ", defaultValStr),
		defaultVal,
		false)
	if err != nil {
		return err
	}

	viper.Set(defaultIssueTypeIDKey, *defaultIssueTypeID)

	return viper.WriteConfig()
}
