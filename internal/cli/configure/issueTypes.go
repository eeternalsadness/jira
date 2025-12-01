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
	issueTypeIds := viper.GetIntSlice("IssueTypeIds")
	fmt.Println("Current issue type IDs:")
	fmt.Println(issueTypeIds)

	fmt.Print("Enter the new list of issue type IDs (separate by commas): ")
	userInput, _ := reader.ReadString('\n')
	userInput = userInput[:len(userInput)-1]
	userInputSlice := strings.Split(userInput, ",")

	// parse issue type ids
	var issueTypeIdsNew []int
	for i := range userInputSlice {
		issueTypeIdString := strings.TrimSpace(userInputSlice[i])
		issueTypeIdInt, err := strconv.Atoi(issueTypeIdString)
		if err != nil {
			fmt.Printf("Invalid project ID: %s", issueTypeIdString)
			return
		}
		issueTypeIdsNew = append(issueTypeIdsNew, issueTypeIdInt)
	}

	if len(issueTypeIds) > 0 {
		overwrite, err := util.UserYesNo("Overwrite existing issue type IDs?")
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		if overwrite {
			viper.Set("IssueTypeIds", issueTypeIdsNew)
		}
	} else {
		viper.Set("IssueTypeIds", issueTypeIdsNew)
	}

	// set default issue type ID
	issueTypeIds = viper.GetIntSlice("IssueTypeIds")
	if len(issueTypeIds) > 0 {
		// default to first issue type ID
		defaultIssueTypeId := issueTypeIds[0]
		util.ViperUpsertInt("DefaultIssueTypeId", "Enter the default issue type ID", strconv.Itoa(defaultIssueTypeId))
	}

	viper.WriteConfig()
}
