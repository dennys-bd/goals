package cmd

import (
	"github.com/spf13/cobra"
)

var scaffoldCmd = &cobra.Command{
	Use:     "scaffold",
	Aliases: []string{"scaf", "s", "Scaffold"},
	Short:   "Scaffold goals structures",
	Long: `Scaffold (goals scaffold) can create many graphql
structures for goals. Check our commands, to see what we
are able to do.`,
}

func init() {
	scaffoldCmd.AddCommand(authCmd)
	scaffoldCmd.AddCommand(gqlCmd)
}
