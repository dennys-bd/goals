package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile, userLicense string

	rootCmd = &cobra.Command{
		Use:   "goals",
		Short: "Goals is Golang/GraphQL Framework",
		Long: `Goals is a Golang/GraphQL Framework in development.
This application is a tool to generate files most used in a go/graphql application,
You can create a new project and create new models with it.
Goals is offering many functions to facilitate your life a GraphQL-Go developer.`,
	}
)

// Execute executes the root command
func Execute() {
	rootCmd.Execute()
}

func init() {
	// rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "author name for copyright attribution")
	// rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
	// rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	viper.SetDefault("author", "DENNYS AZEVEDO <dennys.bd@gmail.com>")
	viper.SetDefault("license", "mit")

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(scaffoldCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(runServerCmd)
}
