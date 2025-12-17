package configure

import (
	"fmt"

	"github.com/eeternalsadness/jira/internal/util"
	"github.com/spf13/viper"
)

// TODO: make the user search for projects using keys instead
func configureDefaultProject() error {
	var defaultVal *string = nil
	defaultValStr := ""
	defaultProjectIDKey := "default_project_id"

	// if there is an existing value, use it as the default project ID
	if viper.IsSet(defaultProjectIDKey) {
		*defaultVal = viper.GetString(defaultProjectIDKey)
		defaultValStr = fmt.Sprintf(" [%s]", *defaultVal)
	}

	defaultProjectID, err := util.UserGetString(
		fmt.Sprintf("Enter the default project ID%s: ", defaultValStr),
		defaultVal,
		false)
	if err != nil {
		return err
	}

	viper.Set(defaultProjectIDKey, *defaultProjectID)

	return viper.WriteConfig()
}
