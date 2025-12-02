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

// configureProjectsCmd represents the configure projects command
var configureProjectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Configure the list of available Jira projects",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		configureProjects()
	},
}

func configureProjects() {
	reader := bufio.NewReader(os.Stdin)

	// get current project ids
	projectIDs := viper.GetIntSlice("ProjectIDs")
	fmt.Println("Current project IDs:")
	fmt.Println(projectIDs)

	fmt.Print("Enter the new list of project IDs (separate by commas): ")
	userInput, _ := reader.ReadString('\n')
	userInput = userInput[:len(userInput)-1]
	userInputSlice := strings.Split(userInput, ",")

	// parse project ids
	var projectIDsNew []int
	for i := range userInputSlice {
		projectIDString := strings.TrimSpace(userInputSlice[i])
		projectIDInt, err := strconv.Atoi(projectIDString)
		if err != nil {
			fmt.Printf("Invalid project ID: %s", projectIDString)
			return
		}
		projectIDsNew = append(projectIDsNew, projectIDInt)
	}

	if len(projectIDs) > 0 {
		overwrite, err := util.UserYesNo("Overwrite existing project IDs?")
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		if overwrite {
			viper.Set("ProjectIDs", projectIDsNew)
		}
	} else {
		viper.Set("ProjectIDs", projectIDsNew)
	}

	// set default project ID
	projectIDs = viper.GetIntSlice("ProjectIDs")
	if len(projectIDs) > 0 {
		// default to first project ID
		defaultProjectID := projectIDs[0]
		util.ViperUpsertInt("DefaultProjectID", "Enter the default project ID", strconv.Itoa(defaultProjectID))
	}

	viper.WriteConfig()
}
