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

func configureProjects() error {
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
			return fmt.Errorf("invalid project ID: %s", projectIDString)
		}
		projectIDsNew = append(projectIDsNew, projectIDInt)
	}

	if len(projectIDs) > 0 {
		overwrite, err := util.UserYesNo("Overwrite existing project IDs?")
		if err != nil {
			return err
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
		err := util.ViperUpsertInt("DefaultProjectID", "Enter the default project ID", strconv.Itoa(defaultProjectID))
		if err != nil {
			return err
		}
	}

	return viper.WriteConfig()
}
