package configure

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/eeternalsadness/jira/internal/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configureIssueTypesCmd represents the configure projects command
var configureIssueTypesCmd = &cobra.Command{
	Use:     "issueTypes",
	Aliases: []string{"types"},
	Short:   "Configure the list of available issue types",
	Args:    cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		configureIssueTypes()
	},
}

func configureIssueTypes() {
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
			fmt.Printf("Invalid project ID: %s", issueTypeIDString)
			return
		}
		issueTypeIDsNew = append(issueTypeIDsNew, issueTypeIDInt)
	}

	if len(issueTypeIDs) > 0 {
		overwrite, err := util.UserYesNo("Overwrite existing issue type IDs?")
		if err != nil {
			fmt.Printf("%s\n", err)
			return
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
		util.ViperUpsertInt("DefaultIssueTypeID", "Enter the default issue type ID", strconv.Itoa(defaultIssueTypeID))
	}

	viper.WriteConfig()
}
