package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile, userLicense string

	rootCmd = &cobra.Command{
		Use:   "goals",
		Short: "Goals is a light boilerplate generator for graphql in golang",
		Long: `Um monte de balbalbla
pq aqui quebra linha e panz`,
	}
)

// Execute executes the root command
func Execute() {
	rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "author name for copyright attribution")
	rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	viper.SetDefault("author", "DENNYS AZEVEDO <dennys.bd@gmail.com>")
	viper.SetDefault("license", "mit")

	rootCmd.AddCommand(initCmd)
}
