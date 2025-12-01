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
	projectIds := viper.GetIntSlice("ProjectIds")
	fmt.Println("Current project IDs:")
	fmt.Println(projectIds)

	fmt.Print("Enter the new list of project IDs (separate by commas): ")
	userInput, _ := reader.ReadString('\n')
	userInput = userInput[:len(userInput)-1]
	userInputSlice := strings.Split(userInput, ",")

	// parse project ids
	var projectIdsNew []int
	for i := range userInputSlice {
		projectIdString := strings.TrimSpace(userInputSlice[i])
		projectIdInt, err := strconv.Atoi(projectIdString)
		if err != nil {
			fmt.Printf("Invalid project ID: %s", projectIdString)
			return
		}
		projectIdsNew = append(projectIdsNew, projectIdInt)
	}

	if len(projectIds) > 0 {
		overwrite, err := util.UserYesNo("Overwrite existing project IDs?")
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		if overwrite {
			viper.Set("ProjectIds", projectIdsNew)
		}
	} else {
		viper.Set("ProjectIds", projectIdsNew)
	}

	// set default project ID
	projectIds = viper.GetIntSlice("ProjectIds")
	if len(projectIds) > 0 {
		// default to first project ID
		defaultProjectId := projectIds[0]
		util.ViperUpsertInt("DefaultProjectId", "Enter the default project ID", strconv.Itoa(defaultProjectId))
	}

	viper.WriteConfig()
}
