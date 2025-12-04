package configure

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/eeternalsadness/jira/internal/util"
	"github.com/spf13/viper"
)

func configureIssueTypes() error {
	reader := bufio.NewReader(os.Stdin)

	// get current project ids
	issueTypeIDs := viper.GetIntSlice("IssueTypeIDs")
	fmt.Println("Current issue type IDs:")
	fmt.Println(issueTypeIDs)

	fmt.Print("Enter the new list of issue type IDs (separate by commas): ")
	userInput, _ := reader.ReadString('\n')
	userInput = userInput[:len(userInput)-1]
	userInputSlice := strings.Split(userInput, ",")

	// parse issue type ids
	var issueTypeIDsNew []int
	for i := range userInputSlice {
		issueTypeIDString := strings.TrimSpace(userInputSlice[i])
		issueTypeIDInt, err := strconv.Atoi(issueTypeIDString)
		if err != nil {
			return fmt.Errorf("invalid project ID: %s", issueTypeIDString)
		}
		issueTypeIDsNew = append(issueTypeIDsNew, issueTypeIDInt)
	}

	if len(issueTypeIDs) > 0 {
		overwrite, err := util.UserYesNo("Overwrite existing issue type IDs?")
		if err != nil {
			return err
		}

		if overwrite {
			viper.Set("IssueTypeIDs", issueTypeIDsNew)
		}
	} else {
		viper.Set("IssueTypeIDs", issueTypeIDsNew)
	}

	// set default issue type ID
	issueTypeIDs = viper.GetIntSlice("IssueTypeIDs")
	if len(issueTypeIDs) > 0 {
		// default to first issue type ID
		defaultIssueTypeID := issueTypeIDs[0]
		err := util.ViperUpsertInt(
			"DefaultIssueTypeID",
			"Enter the default issue type ID",
			&defaultIssueTypeID)
		if err != nil {
			return err
		}
	}

	return viper.WriteConfig()
}
