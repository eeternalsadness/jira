/*
Copyright © 2025 Bach Nguyen <69bnguyen@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package version

import (
	"fmt"

	"github.com/eeternalsadness/jira/internal/util"
	"github.com/spf13/cobra"
)

// NewCommand creates and returns the version command
func NewCommand() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show the version for the jira CLI tool",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// override the root's PersistentPreRunE since we don't need to init config
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("%s version %s\n", cmd.Root().Name(), util.Version)
			fmt.Printf("Go version: %s\n", util.GoVersion)
			fmt.Printf("Git commit SHA: %s\n", util.GitCommitSHA)
			return util.CheckVersion(cmd)
		},
	}

	return versionCmd
}
