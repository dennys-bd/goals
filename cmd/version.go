package cmd

import "github.com/spf13/cobra"

var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"Version", "v"},
	Short:   "Version says the current version of goals",
	Run: func(cmd *cobra.Command, args []string) {
		println("Goals Version: Beta 0.4")
	},
}
