package util

import (
	"fmt"
	"strings"
)

type ViperKey string

const (
	JiraDomainKey ViperKey = "domain"
	JiraEmailKey  ViperKey = "email"
	JiraTokenKey  ViperKey = "token"

	DefaultProjectIDKey   ViperKey = "default_project_id"
	DefaultIssueTypeIDKey ViperKey = "default_issue_type_id"
)

func SensorString(str string) string {
	// otherwise show the first and last 25% chars (max 4)
	charsToShow := len(str) / 4
	charsToShow = min(charsToShow, 4)

	// limit the sensored part to 20 chars
	sensoredLen := min(len(str)-charsToShow*2, 20)

	return fmt.Sprintf("%s%s%s", str[:charsToShow], strings.Repeat("*", sensoredLen), str[len(str)-charsToShow:])
}
